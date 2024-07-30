package main

import (
	"flag"
	"fmt"
	cma "golang.conradwood.net/apis/common"
	pb "golang.conradwood.net/apis/protorenderer2"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/linux"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/protorenderer/v2/common"
	"golang.conradwood.net/protorenderer/v2/compilers/golang"
	//	"golang.conradwood.net/protorenderer/v2/helpers"
	"io"
	"os"
	"time"
	//	"golang.conradwood.net/protorenderer/v2/compilers/meta" // NOTE: different type of compiler! (parses comments, assigns IDs etc)
	"golang.conradwood.net/protorenderer/v2/interfaces"
	"golang.yacloud.eu/yatools/gitrepo"
	"path/filepath"
	"strings"
)

var (
	display_versioninfo = flag.String("display_versioninfo", "", "if set to a filename, will display versioninfo.pbbin as yaml")
	deps                = flag.String("dependencies", "", "if set, finds files that depend on that filename")
	get_vinfo           = flag.Bool("get_versioninfo", false, "if true get versioninfo, print it out and exit")
	save                = flag.Bool("save", false, "if true, and compilation is successful, add the proto and its artefacts to protorenderers store")
	version             = flag.Uint64("version", 0, "work on this version, 0 (default) is latest")
	edit_store          = flag.Bool("edit_store", false, "if true, checkout store, wait for key and save it again (needs -token and -ge_disable_user_token)")
)

func main() {
	flag.Parse()
	if *edit_store {
		utils.Bail("failed to edit store", EditStore())
	} else if *display_versioninfo != "" {
		utils.Bail("failed to display version info", DisplayVersionInfo())
	} else if *deps != "" {
		utils.Bail("failed to do reverse dependencies", ReverseDeps())
	} else if *get_vinfo {
		utils.Bail("failed to get versioninfo", GetVersionInfo())
	} else {
		if len(flag.Args()) != 0 {
			for _, a := range flag.Args() {
				submit_protos_with_dir(a)
			}
			os.Exit(0)
		}
		submit_protos()
	}
	fmt.Printf("Done\n")
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
	sce := &StandardCompilerEnvironment{workdir: "/tmp/pr/v2"}

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
	utils.Bail("not a yagitrepo", err)
	proto_dir := gr.Path() + "/protos/"
	submit_protos_with_dir(proto_dir)
}
func submit_protos_with_dir(proto_dir string) {
	ctx := authremote.ContextWithTimeout(time.Duration(30) * time.Second)

	// repoid
	repoid := uint32(0)
	gr, err := gitrepo.NewYAGitRepo(proto_dir)
	if err != nil {
		fmt.Printf("Not a YAGitRepo: \"%s\"\n", proto_dir)
	} else {
		repoid = uint32(gr.RepositoryID())
	}
	srv, err := pb.GetProtoRenderer2Client().Submit(ctx)
	utils.Bail("failed to stream files to server", err)
	so := &pb.SubmitOption{Save: *save}
	err = srv.Send(&pb.FileTransfer{SubmitOption: so})
	utils.Bail("failed to stream options to server", err)
	bs := utils.NewByteStreamSender(func(key, filename string) error {
		// start new file
		err := srv.Send(&pb.FileTransfer{Filename: filename, RepositoryID: repoid})
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
	utils.Bail("failed to send files", err)
	// receiving the results now...
	for {
		recv, err := srv.Recv()
		if recv != nil {
			if recv.Filename != "" {
				fmt.Printf("Receiving: filename=%s, bytes=%d\n", recv.Filename, len(recv.Data))
			}
			if recv.Result != nil {
				fmt.Printf("Failure for %s: %v\n", recv.Result.Filename, recv.Result.Filename)
				for _, result := range recv.Result.CompileResults {
					fmt.Printf("    compiler: \"%s\" (success=%v)\n", result.CompilerName, result.Success)
					fmt.Printf("    error: %s\n", result.ErrorMessage)
					fmt.Printf("    output: %s\n", result.Output)
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
func GetVersionInfo() error {
	ctx := authremote.Context()
	vi, err := pb.GetProtoRenderer2Client().GetVersionInfo(ctx, &cma.Void{})
	if err != nil {
		return err
	}
	t := utils.Table{}
	t.AddHeaders("filename", "failed", "msg")
	for _, vf := range vi.Files {
		t.AddString(vf.Filename)
		failed := false
		s := "OK"
		for _, cr := range vf.FileResult.CompileResults {
			if !cr.Success {
				failed = true
				s = cr.ErrorMessage + "\n" + cr.Output
			}
		}
		t.AddBool(failed)
		t.AddString(s)
		t.NewRow()
	}
	fmt.Println(t.ToPrettyString())
	return nil
}

func ReverseDeps() error {
	fname := *deps
	fmt.Printf("Getting dependencies for \"%s\"\n", fname)
	ctx := authremote.Context()
	req := &pb.ReverseDependenciesRequest{Filename: fname, MaxDepth: 0}
	res, err := pb.GetProtoRenderer2Client().GetReverseDependencies(ctx, req)
	if err != nil {
		return err
	}
	fmt.Printf("%d files depend on \"%s\":\n", len(res.Filenames), req.Filename)
	for _, filename := range res.Filenames {
		fmt.Printf("%s\n", filename)
	}

	return nil
}

func DisplayVersionInfo() error {
	b, err := utils.ReadFile(*display_versioninfo)
	if err != nil {
		return err
	}
	vi := &pb.VersionInfo{}
	err = utils.UnmarshalBytes(b, vi)
	if err != nil {
		return err
	}
	fname := "/tmp/versioninfo.yaml"
	err = utils.WriteYaml(fname, vi)
	if err != nil {
		return err
	}
	fmt.Printf("Written to %s\n", fname)
	return nil
}
