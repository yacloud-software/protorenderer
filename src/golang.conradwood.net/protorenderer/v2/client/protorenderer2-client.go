package main

import (
	"flag"
	"fmt"
	pb "golang.conradwood.net/apis/protorenderer2"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/linux"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/protorenderer/v2/common"
	"golang.conradwood.net/protorenderer/v2/compilers/golang"
	"io"
	"os"
	//	"golang.conradwood.net/protorenderer/v2/compilers/meta" // NOTE: different type of compiler! (parses comments, assigns IDs etc)
	"golang.conradwood.net/protorenderer/v2/interfaces"
	"golang.yacloud.eu/yatools/gitrepo"
	"path/filepath"
	"strings"
)

func main() {
	flag.Parse()
	if len(flag.Args()) != 0 {
		for _, a := range flag.Args() {
			submit_protos_with_dir(a)
		}
		os.Exit(0)
	}
	submit_protos()
}
func local_compile() {
	outdir := "compile_result/protos"
	// copy the test files across
	fd, err := utils.FindFile("extra/test_protos/previous_protos")
	utils.Bail("previous protos dir not found", err)

	protodir, pfs, err := find_protofiles()
	utils.Bail("failed to find protofiles", err)
	golang_compiler := golang.New()
	ctx := authremote.Context()
	for _, pf := range pfs {
		fmt.Printf("Compiling file: %s\n", pf.GetFilename())
	}
	scr := &common.StandardCompileResult{}
	sce := &StandardCompilerEnvironment{workdir: "/tmp/pr/v2", knownprotosdir: "proto_files/protos", newprotosdir: "new_protos/protos"}

	mkdir(sce.WorkDir())
	mkdir(sce.WorkDir() + "/" + outdir)
	mkdir(sce.AllKnownProtosDir())
	mkdir(sce.NewProtosDir())

	// copy the files from extra into the right place
	err = linux.CopyDir(protodir, sce.WorkDir()+"/"+sce.NewProtosDir())
	utils.Bail("failed to copy new protos", err)
	err = linux.CopyDir(fd, filepath.Dir(sce.WorkDir()+"/"+sce.AllKnownProtosDir())) // to avoid doubling "protos", strip it
	utils.Bail("failed to copy prev protos", err)

	//	meta.New().Compile(ctx, "randomstring", 4012, sce, pfs, sce.WorkDir()+"/"+outdir, scr) // need server for this to work
	pb.GetProtoRenderer2Client()
	golang_compiler.Compile(ctx, sce, pfs, sce.WorkDir()+"/"+outdir, scr)

}

func find_protofiles() (string, []interfaces.ProtoFile, error) {
	dir, err := utils.FindFile("extra/test_protos/new_protos/protos")
	if err != nil {
		return dir, nil, err
	}
	var filenames []string
	utils.DirWalk(dir, func(root, relfil string) error {
		if !strings.HasSuffix(relfil, ".proto") {
			return nil
		}
		filenames = append(filenames, relfil)
		return nil
	},
	)
	var res []interfaces.ProtoFile
	for _, fn := range filenames {
		res = append(res, &StandardProtoFile{filename: fn})
	}
	return dir, res, nil
}

func mkdir(dir string) {
	err := linux.CreateIfNotExists(dir, 0777)
	utils.Bail("failed to create dir", err)
}

func submit_protos() {
	path, err := os.Getwd()
	utils.Bail("no current dir", err)
	gr, err := gitrepo.NewGitRepo(path)
	proto_dir := gr.Path() + "/protos/"
	submit_protos_with_dir(proto_dir)
}
func submit_protos_with_dir(proto_dir string) {
	ctx := authremote.Context()
	srv, err := pb.GetProtoRenderer2Client().Compile(ctx)
	utils.Bail("failed to stream files to server", err)
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
	err = utils.DirWalk(proto_dir, func(root, relfil string) error {
		if !strings.HasSuffix(relfil, ".proto") {
			return nil
		}
		ct, err := utils.ReadFile(proto_dir + "/" + relfil)
		if err != nil {
			return err
		}
		fmt.Printf("Submitting %s (%d bytes)\n", "protos/"+relfil, len(ct))
		err = bs.SendBytes("foo", relfil, ct)
		return err
	},
	)
	utils.Bail("failed to read proto files", err)
	err = srv.Send(&pb.FileTransfer{TransferComplete: true}) // switching to recv mode now
	// receiving the results now...
	for {
		recv, err := srv.Recv()
		if recv != nil {
			if recv.Filename != "" {
				fmt.Printf("Receiving: filename=%s, bytes=%d\n", recv.Filename, len(recv.Data))
			}
			if recv.Result != nil {
				fmt.Printf("Failure for %s: %v\n", recv.Result.Filename, recv.Result.Filename)
				if recv.Result.Failed {
					for _, result := range recv.Result.Failures {
						fmt.Printf("    compiler: \"%s\"\n", result.CompilerName)
						fmt.Printf("    error: %s\n", result.ErrorMessage)
						fmt.Printf("    output: %s\n", result.Output)
					}
				}
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			utils.Bail("failed to receive files", err)
		}
	}
	fmt.Printf("Done\n")

}
