package srv

import (
	"context"
	"fmt"
	"golang.conradwood.net/protorenderer/cmdline"
	"golang.conradwood.net/protorenderer/v2/common"
	"golang.conradwood.net/protorenderer/v2/compilers/golang"
	"golang.conradwood.net/protorenderer/v2/compilers/java"
	"golang.conradwood.net/protorenderer/v2/helpers"
	"golang.conradwood.net/protorenderer/v2/interfaces"
	"golang.conradwood.net/protorenderer/v2/meta_compiler"
	//	"golang.conradwood.net/protorenderer/v2/store"
	"golang.conradwood.net/protorenderer/v2/versioninfo"
	"strings"
	//	pb1 "golang.conradwood.net/apis/protorenderer"
	pb "golang.conradwood.net/apis/protorenderer2"
	//	"golang.conradwood.net/go-easyops/errors"
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
	return err
}
func compile(srv compile_serve_req, save_on_success bool) error {
	compile_lock.Lock()
	defer compile_lock.Unlock()
	ce := CompileEnv
	fmt.Printf("Compile directories:\n")
	fmt.Printf("  NewProtosDir   : %s\n", ce.NewProtosDir())
	fmt.Printf("  CompilerOutDir : %s\n", ce.CompilerOutDir())

	err := utils.RecreateSafely(ce.NewProtosDir())
	if err != nil {
		return err
	}
	od := ce.CompilerOutDir()
	err = utils.RecreateSafely(od)
	if err != nil {
		return err
	}

	opt, err := receive(ce, srv)
	if err != nil {
		return err
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
		return err
	}

	if opt.Save || opt.IncludeMeta {
		fmt.Printf("[compile] starting meta compiler with %d files\n", len(pfs))
		meta_compiler := meta_compiler.New()
		err = meta_compiler.Compile(ctx, ce, pfs, od, scr)
		if err != nil {
			return err
		}
	}
	// after compiling with meta, we remove any proto files that failed compilation by meta
	// it is not worth compiling them with any other compiles
	pfs = remove_broken(pfs, scr)
	if len(pfs) == 0 {
		return fmt.Errorf("meta compiler failed to compile any files")
	}
	versioninfo.New()

	// compile protos
	err = compile_all_compilers(ctx, ce, scr, pfs)
	if err != nil {
		return err
	}

	// now send return
	err = send(ce, srv, od)
	if err != nil {
		return err
	}

	err = send_failures(srv, scr, pfs)
	if err != nil {
		return err
	}
	for _, pf := range pfs {
		if helpers.ContainsFailure(scr.GetResults(pf)) {
			return fmt.Errorf("at least one error occured")
		}
	}

	if opt.Save {
		err = recompile_dependencies_with_err(ctx, ce, pfs)
		if err != nil {
			fmt.Printf("Recompiling dependencies: %s\n", err)
			return fmt.Errorf("failed to recompile dependencies: %s\n", err)
		}
		fmt.Printf("[compile] saving new protos to store...\n")
		//helpers.MergeCompilerEnvironment(ce, true)
	}
	fmt.Printf("[compile] completed\n")

	return nil
}

// compile all files with all enabled compilers and place results in ce.CompilerOutDir()
func compile_all_compilers(ctx context.Context, ce interfaces.CompilerEnvironment, scr interfaces.CompileResult, pfs []interfaces.ProtoFile) error {
	od := ce.CompilerOutDir()
	compilers := []interfaces.Compiler{golang.New()}
	if cmdline.GetCompilerEnabledJava() {
		compilers = append(compilers, java.New())
	}

	for _, comp := range compilers {
		fmt.Printf("[compile] starting \"%s\" compiler with %d files\n", comp.ShortName(), len(pfs))
		err := comp.Compile(ctx, ce, pfs, od, scr)
		if err != nil {
			return err
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
			return err
		}
	}

	return nil
}

func send(ce interfaces.CompilerEnvironment, srv compile_serve_req, dir string) error {
	bs := utils.NewByteStreamSender(func(key, filename string) error {
		// start new file
		err := srv.Send(&pb.FileTransfer{Filename: filename})
		return err
	},
		// send contents
		func(b []byte) error {
			err := srv.Send(&pb.FileTransfer{Data: b})
			return err
		},
	)
	err := utils.DirWalk(dir, func(root, relfil string) error {
		for _, ign := range IGNORE_FILES {
			if relfil == ign {
				return nil
			}
		}
		ct, err := utils.ReadFile(root + "/" + relfil)
		if err != nil {
			return err
		}
		fmt.Printf("Sending %s (%d bytes)\n", relfil, len(ct))
		err = bs.SendBytes("foo", relfil, ct)
		return err
	},
	)
	return err
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
					return nil, fmt.Errorf("Invalid filename starting with two proto dirs (%s)", rcv.Filename)
				}
				if strings.HasPrefix(rcv.Filename, "protos/") {
					return nil, fmt.Errorf("Invalid filename starting with 'protos' (%s)", rcv.Filename)
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

// recompile all the dependencies on the given file(s)...
func recompile_dependencies_with_err(ctx context.Context, ce interfaces.CompilerEnvironment, pfs []interfaces.ProtoFile) error {
	err := MetaCache.readAllIfNecessary()
	if err != nil {
		return err
	}
	scr := &common.StandardCompileResult{}
	for _, pf := range pfs {
		err := recompile_dependencies(ctx, ce, scr, pf)
		if err != nil {
			return err
		}
	}
	return nil
}

// recompile any files that directly or indirectly import "pf"
func recompile_dependencies(ctx context.Context, ce interfaces.CompilerEnvironment, scr interfaces.CompileResult, pf interfaces.ProtoFile) error {
	pfs, err := MetaCache.AllWithDependencyOn(pf.GetFilename(), 0)
	if err != nil {
		return err
	}
	var cpfs []interfaces.ProtoFile
	for _, npf := range pfs {
		spf := &helpers.StandardProtoFile{Filename: npf.ProtoFile.Name}
		cpfs = append(cpfs, spf)
	}

	err = compile_all_compilers(ctx, ce, scr, cpfs)
	if err != nil {
		return err
	}
	for _, cpf := range cpfs {
		if helpers.ContainsFailure(scr.GetResults(cpf)) {
			return fmt.Errorf("failed to compile dependency \"%s\"\n", pf.GetFilename())
		}
	}
	return nil
}
