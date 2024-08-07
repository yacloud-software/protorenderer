package main

import (
	"flag"
	"fmt"
	cma "golang.conradwood.net/apis/common"
	pb "golang.conradwood.net/apis/protorenderer2"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/linux"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/protorenderer/v2/client/protosubmitter"
	"golang.conradwood.net/protorenderer/v2/common"
	"golang.conradwood.net/protorenderer/v2/compilers/golang"
	"golang.conradwood.net/protorenderer/v2/interfaces"
	"os"
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
	var err error
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
				if *save {
					err = protosubmitter.SubmitProtos(a)
				} else {
					err = protosubmitter.CompileProtos(a)
				}
				utils.Bail("failed to compile", err)
			}
			os.Exit(0)
		}
		utils.Bail("failed to submit", protosubmitter.SubmitProtosGit())
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
