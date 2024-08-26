package srv

import (
	"context"
	"fmt"
	"golang.conradwood.net/protorenderer/cmdline"
	"golang.conradwood.net/protorenderer/v2/common"
	"golang.conradwood.net/protorenderer/v2/compilers/golang"
	"golang.conradwood.net/protorenderer/v2/compilers/java"
	"golang.conradwood.net/protorenderer/v2/compilers/nanopb"
	"golang.conradwood.net/protorenderer/v2/helpers"
	"golang.conradwood.net/protorenderer/v2/interfaces"
	"golang.conradwood.net/protorenderer/v2/meta_compiler"
	"golang.conradwood.net/protorenderer/v2/store"
	"path/filepath"
	//	"golang.conradwood.net/protorenderer/v2/store"
	"strings"
	//	pb1 "golang.conradwood.net/apis/protorenderer"
	pb "golang.conradwood.net/apis/protorenderer2"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/utils"
	"io"
	"sync"
)

var (
	compile_lock sync.Mutex
	IGNORE_FILES = []string{".goeasyops-dir"}
)

type compile_serve_req interface {
	Recv() (*pb.FileTransfer, error)
	Context() context.Context
	Send(*pb.FileTransfer) error
}

func (pr *protoRenderer) Submit(srv pb.ProtoRenderer2_SubmitServer) error {
	err := compile(srv, false)
	if err != nil {
		fmt.Printf("Error compiling: %s\n", err)
	}
	return errors.Wrap(err)
}
func compile(srv compile_serve_req, save_on_success bool) error {
	compile_lock.Lock()
	defer compile_lock.Unlock()
	fmt.Printf("[compile] -------------------- started ------------------\n")
	ce := CompileEnv.Fork()
	//	sce, itis := ce.(*StandardCompilerEnvironment)
	//	if itis {
	ce.server = srv
	//	}
	utils.RecreateSafely(ce.CompilerOutDir())
	fmt.Printf("Compile directories:\n")
	fmt.Printf("  NewProtosDir   : %s\n", ce.NewProtosDir())
	fmt.Printf("  CompilerOutDir : %s\n", ce.CompilerOutDir())

	err := utils.RecreateSafely(ce.NewProtosDir())
	if err != nil {
		return errors.Wrap(err)
	}
	od := ce.CompilerOutDir()
	err = utils.RecreateSafely(od)
	if err != nil {
		return errors.Wrap(err)
	}

	opt, err := receive(ce, srv)
	if err != nil {
		return errors.Wrap(err)
	}
	if opt != nil {
		fmt.Printf("[compile] Options: %#v\n", opt)
	} else {
		opt = &pb.SubmitOption{}
	}
	scr := &common.StandardCompileResult{}
	ctx := srv.Context()

	pfs, err := helpers.FindProtoFiles(ce.NewProtosDir())
	if err != nil {
		return errors.Wrap(err)
	}
	submitted_proto_list := pfs // need that, because we are removing from pfs all those which the meta compiler failed

	//	if opt.Save || opt.IncludeMeta {
	ce.Printf("[compile] starting meta compiler with %d files\n", len(pfs))
	meta_compiler := meta_compiler.New()
	err = meta_compiler.Compile(ctx, ce, pfs, od, scr)
	if err != nil {
		return errors.Wrap(err)
	}

	// after compiling with meta, we remove any proto files that failed compilation by meta
	// it is not worth compiling them with any other compiles
	pfs = remove_broken(pfs, scr)
	if len(pfs) == 0 {
		for _, failed := range scr.GetFailed() {
			fmt.Printf("Meta compiler failed: %s\n", failed.ErrorMessage)
		}
		send_all_broken_ones(ce, submitted_proto_list, scr)
		return errors.Errorf("meta compiler failed to compile any files")
	}

	err = check_if_files_match_packages(ctx, ce, scr, pfs)
	if err != nil {
		send_all_broken_ones(ce, submitted_proto_list, scr)
		return errors.Errorf("failed to check if filenames match packages: %s", err)
	}

	pfs = remove_broken(pfs, scr)
	if len(pfs) == 0 {
		send_all_broken_ones(ce, submitted_proto_list, scr)
		return errors.Errorf("after package check, no files left to compile")
	}

	// compile protos
	compilers := getCompilerArray()

	err = compile_all_compilersWithCompilerArray(ctx, ce, scr, pfs, compilers)
	if err != nil {
		return errors.Wrap(err)
	}

	// now send return
	err = send(ce, srv, od)
	if err != nil {
		return errors.Wrap(err)
	}

	err = send_failures(srv, scr, pfs)
	if err != nil {
		return errors.Wrap(err)
	}
	send_all_broken_ones(ce, submitted_proto_list, scr)

	if opt.Save {
		cr := NewCompareResult(scr, &common.StandardCompileResult{}) // TODO fix compile results
		if cr.IsWorse() {
			return errors.Errorf("the new protos are worse than what is in store. not submitting")
		}
		dc := NewDependencyCompiler(ce, pfs)
		err = dc.Recompile(ctx)
		//		err = recompile_dependencies_with_err(ctx, ce, pfs, compilers)
		if err != nil {
			fmt.Printf("Recompiling dependencies: %s\n", err)
			return errors.Errorf("failed to recompile dependencies: %s\n", err)
		}
		fmt.Printf("[compile] saving new protos to store...\n")
		err := helpers.MergeCompilerEnvironment(ce, CompileEnv, true)
		if err != nil {
			return errors.Wrap(err)
		}
		//		currentVersionInfo.SetDirty()
		err = saveVersionInfo() // trigger storage merge
		if err != nil {
			return errors.Wrap(err)
		}
		store.TriggerUpload(ce.StoreDir())

	}
	fmt.Printf("[compile] -------------------- completed ------------------\n")
	//	utils.RecreateSafely(ce.CompilerOutDir())
	return nil
}

// compile all files with all enabled compilers and place results in ce.CompilerOutDir()
func compile_all_compilers(ctx context.Context, ce interfaces.CompilerEnvironment, scr interfaces.CompileResult, pfs []interfaces.ProtoFile) error {
	compilers := getCompilerArray()
	return compile_all_compilersWithCompilerArray(ctx, ce, scr, pfs, compilers)
}
func compile_all_compilersWithCompilerArray(ctx context.Context, ce interfaces.CompilerEnvironment, scr interfaces.CompileResult, pfs []interfaces.ProtoFile, compilers []interfaces.Compiler) error {
	od := ce.CompilerOutDir()

	for _, comp := range compilers {
		ce.Printf("[compile] starting \"%s\" compiler with %d files\n", comp.ShortName(), len(pfs))
		if len(pfs) == 0 {
			continue
		}
		err := comp.Compile(ctx, ce, pfs, od, scr)
		if err != nil {
			fmt.Printf("[compile] Compiler \"%s\" failed: %s\n", comp.ShortName(), err)
			return errors.Wrap(err)
		}
		for _, pf := range pfs {
			if !helpers.ContainsFailure(scr.GetResults(pf)) {
				scr.AddSuccess(comp, pf)
			}
		}
		fmt.Printf("[compile] compiler \"%s\" completed\n", comp.ShortName())
	}

	return nil
}

func send_failures(srv compile_serve_req, scr *common.StandardCompileResult, pfs []interfaces.ProtoFile) error {
	// build the compile result
	result := make(map[string]*pb.FileResult) // filename->result (only failed files)
	for _, pf := range pfs {                  // send result for every file that failed
		if !helpers.ContainsFailure(scr.GetResults(pf)) {
			continue
		}

		fr := result[pf.GetFilename()]
		if fr == nil {
			fr = &pb.FileResult{Filename: pf.GetFilename()}
			result[pf.GetFilename()] = fr
		}
		fr.CompileResults = append(fr.CompileResults, scr.GetResults(pf)...)
	}
	fmt.Printf("%d failures\n", len(result))
	for _, failure := range result {
		err := srv.Send(&pb.FileTransfer{Result: failure})
		if err != nil {
			return errors.Wrap(err)
		}
	}

	return nil
}

func send(ce interfaces.CompilerEnvironment, srv compile_serve_req, dir string) error {
	bs := utils.NewByteStreamSender(func(key, filename string) error {
		// start new file
		err := srv.Send(&pb.FileTransfer{Filename: filename})
		return errors.Wrap(err)
	},
		// send contents
		func(b []byte) error {
			err := srv.Send(&pb.FileTransfer{Data: b})
			return errors.Wrap(err)
		},
	)
	err := utils.DirWalk(dir, func(root, relfil string) error {
		for _, ign := range IGNORE_FILES {
			if relfil == ign {
				return nil
			}
			b := filepath.Base(relfil)
			if b == ign {
				return nil
			}
		}
		ct, err := utils.ReadFile(root + "/" + relfil)
		if err != nil {
			return errors.Wrap(err)
		}
		//		fmt.Printf("Sending %s (%d bytes)\n", relfil, len(ct))
		err = bs.SendBytes("foo", relfil, ct)
		return errors.Wrap(err)
	},
	)
	return errors.Wrap(err)
}
func receive(ce interfaces.CompilerEnvironment, srv compile_serve_req) (*pb.SubmitOption, error) {
	//	ctx := srv.Context()
	bsr := utils.NewByteStreamReceiver(ce.NewProtosDir())
	var opt *pb.SubmitOption
	for {
		rcv, err := srv.Recv()
		if rcv != nil {
			if rcv.TransferComplete {
				break
			}
			if rcv.SubmitOption != nil {
				opt = rcv.SubmitOption
				continue
			}
			if rcv.Filename != "" {
				if strings.HasPrefix(rcv.Filename, "protos/protos/") {
					return nil, errors.Errorf("Invalid filename starting with two proto dirs (%s)", rcv.Filename)
				}
				if strings.HasPrefix(rcv.Filename, "protos/") {
					return nil, errors.Errorf("Invalid filename starting with 'protos' (%s)", rcv.Filename)
				}
				/*
					if do_persist_filename {
						_, err := store.GetOrCreateFile(ctx, rcv.Filename, uint64(rcv.RepositoryID))
						if err != nil {
							return nil, err
						}
					}
				*/
			}
			err = bsr.NewData(rcv)
			if err != nil {
				return nil, err
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}
	err := bsr.Close()
	if err != nil {
		return nil, err
	}
	return opt, nil
}

func remove_broken(pfs []interfaces.ProtoFile, scr interfaces.CompileResult) []interfaces.ProtoFile {
	var res []interfaces.ProtoFile
	for _, pf := range pfs {
		if helpers.ContainsFailure(scr.GetResults(pf)) {
			continue
		}
		res = append(res, pf)
	}
	return res
}

func send_all_broken_ones(ce interfaces.CompilerEnvironment, pfs []interfaces.ProtoFile, cr interfaces.CompileResult) {
	for _, pf := range pfs {
		results := cr.GetResults(pf)
		if !helpers.ContainsFailure(results) {
			continue
		}
		fmt.Printf("file \"%s\" failed:\n", pf.GetFilename())
		for _, comp_err := range results {
			if comp_err.Success {
				continue
			}
			ce.Printf("    %s, file: \"%s\": %s (%s)", comp_err.CompilerName, pf.GetFilename(), comp_err.ErrorMessage, comp_err.Output)
		}
	}
}

func getCompilerArray() []interfaces.Compiler {
	compilers := []interfaces.Compiler{golang.New()}
	if cmdline.GetCompilerEnabledJava() {
		compilers = append(compilers, java.New())
	}
	if cmdline.GetCompilerEnabledNanoPB() {
		compilers = append(compilers, nanopb.New())
	}
	return compilers
}
