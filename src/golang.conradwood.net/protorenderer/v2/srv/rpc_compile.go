package srv

import (
	"fmt"
	"golang.conradwood.net/protorenderer/v2/compilers/golang"
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

func (pr *protoRenderer) Compile(srv pb.ProtoRenderer2_CompileServer) error {
	compile_lock.Lock()
	defer compile_lock.Unlock()
	ce := CompileEnv

	od := ce.WorkDir() + "/compile_outdir/protos"
	err := utils.RecreateSafely(od)
	if err != nil {
		return err
	}

	err = receive(ce, srv)
	if err != nil {
		return err
	}
	scr := &StandardCompileResult{}
	ctx := srv.Context()

	pfs, err := helpers.FindFiles(ce.WorkDir()+"/"+ce.NewProtosDir(), ".proto")
	if err != nil {
		return err
	}

	golang_compiler := golang.New()
	err = golang_compiler.Compile(ctx, ce, pfs, od, scr)
	if err != nil {
		return err
	}

	// TODO: check compile result!!

	err = send(ce, srv, od)
	if err != nil {
		return err
	}

	return nil
}
func send(ce interfaces.CompilerEnvironment, srv pb.ProtoRenderer2_CompileServer, dir string) error {
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
func receive(ce interfaces.CompilerEnvironment, srv pb.ProtoRenderer2_CompileServer) error {
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
