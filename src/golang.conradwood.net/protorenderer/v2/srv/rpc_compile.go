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
	//	pb1 "golang.conradwood.net/apis/protorenderer"
	pb "golang.conradwood.net/apis/protorenderer2"
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

	od := ce.WorkDir() + "/compile_outdir"
	err := utils.RecreateSafely(od)
	if err != nil {
		return err
	}

	err = receive(ce, srv)
	if err != nil {
		return err
	}
	scr := &common.StandardCompileResult{}
	ctx := srv.Context()

	pfs, err := helpers.FindProtoFiles(ce.WorkDir() + "/" + ce.NewProtosDir())
	if err != nil {
		return err
	}

	golang_compiler := golang.New()
	scr.SetCompiler(golang_compiler) // to mark errors as such
	err = golang_compiler.Compile(ctx, ce, pfs, od, scr)
	if err != nil {
		return err
	}
	if cmdline.GetCompilerEnabledJava() {
		java_compiler := java.New()
		scr.SetCompiler(java_compiler) // to mark errors as such
		err = java_compiler.Compile(ctx, ce, pfs, od, scr)
		if err != nil {
			return err
		}
	}

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
		fmt.Printf("Submitting %s (%d bytes)\n", "protos/"+relfil, len(ct))
		err = bs.SendBytes("foo", relfil, ct)
		return err
	},
	)
	return err
}
func receive(ce interfaces.CompilerEnvironment, srv compile_serve_req) error {
	bsr := utils.NewByteStreamReceiver(ce.WorkDir() + "/" + ce.NewProtosDir())
	for {
		rcv, err := srv.Recv()
		if rcv != nil {
			if rcv.TransferComplete {
				break
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
