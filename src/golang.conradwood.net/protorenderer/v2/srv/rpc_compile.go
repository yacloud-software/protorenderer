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
	"golang.conradwood.net/protorenderer/v2/store"
	"golang.conradwood.net/protorenderer/v2/versioninfo"
	"strings"
	//	pb1 "golang.conradwood.net/apis/protorenderer"
	pb "golang.conradwood.net/apis/protorenderer2"
	//	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/utils"
	"io"
)

var (
	IGNORE_FILES = []string{".goeasyops-dir"}
)

type compile_serve_req interface {
	Recv() (*pb.FileTransfer, error)
	Context() context.Context
	Send(*pb.FileTransfer) error
}

func (pr *protoRenderer) Compile(srv pb.ProtoRenderer2_CompileServer) error {
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

	err = receive(ce, srv, save_on_success)
	if err != nil {
		return err
	}
	scr := &common.StandardCompileResult{}
	ctx := srv.Context()

	pfs, err := helpers.FindProtoFiles(ce.NewProtosDir())
	if err != nil {
		return err
	}

	fmt.Printf("[compile] starting meta compiler with %d files\n", len(pfs))
	meta_compiler := meta_compiler.New()
	err = meta_compiler.Compile(ctx, ce, pfs, od, scr)
	if err != nil {
		return err
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
			if len(scr.GetFailures(pf)) == 0 {
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
		if len(scr.GetFailures(pf)) == 0 {
			continue
		}

		fr := result[pf.GetFilename()]
		if fr == nil {
			fr = &pb.FileResult{Filename: pf.GetFilename(), Failed: true}
			result[pf.GetFilename()] = fr
		}
		fr.Failures = append(fr.Failures, scr.GetFailures(pf)...)
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
		fmt.Printf("Submitting %s (%d bytes)\n", relfil, len(ct))
		err = bs.SendBytes("foo", relfil, ct)
		return err
	},
	)
	return err
}
func receive(ce interfaces.CompilerEnvironment, srv compile_serve_req, do_persist_filename bool) error {
	ctx := srv.Context()
	bsr := utils.NewByteStreamReceiver(ce.NewProtosDir())
	for {
		rcv, err := srv.Recv()
		if rcv != nil {
			if rcv.TransferComplete {
				break
			}
			if rcv.Filename != "" {
				if strings.HasPrefix(rcv.Filename, "protos/protos/") {
					return fmt.Errorf("Invalid filename starting with two proto dirs (%s)", rcv.Filename)
				}
				if strings.HasPrefix(rcv.Filename, "protos/") {
					return fmt.Errorf("Invalid filename starting with 'protos' (%s)", rcv.Filename)
				}
				if do_persist_filename {
					_, err := store.GetOrCreateFile(ctx, rcv.Filename, uint64(rcv.RepositoryID))
					if err != nil {
						return err
					}
				}
			}
			err = bsr.NewData(rcv)
			if err != nil {
				return err
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}
	err := bsr.Close()
	if err != nil {
		return err
	}
	return nil
}

func remove_broken(pfs []interfaces.ProtoFile, scr interfaces.CompileResult) []interfaces.ProtoFile {
	var res []interfaces.ProtoFile
	for _, pf := range pfs {
		if len(scr.GetFailures(pf)) != 0 {
			continue
		}
		res = append(res, pf)
	}
	return res
}
