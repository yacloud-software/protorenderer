package golang

import (
	"context"
	"fmt"
	"golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/linux"
	cc "golang.conradwood.net/protorenderer/v2/compilers/common"
	"golang.conradwood.net/protorenderer/v2/interfaces"
)

type golangCompiler struct{}

func New() interfaces.Compiler {
	return &golangCompiler{}
}
func (gc *golangCompiler) ShortName() string { return "golang" }
func (gc *golangCompiler) Compile(ctx context.Context, ce interfaces.CompilerEnvironment, files []interfaces.ProtoFile, outdir string, cr interfaces.CompileResult) error {
	targetdir := outdir
	dir := ce.WorkDir() + "/" + ce.NewProtosDir()
	import_dirs := []string{
		dir,
		ce.WorkDir() + "/" + ce.AllKnownProtosDir(),
	}
	/***************************** compile .proto -> .pb.go ******************************/
	pcfname := cc.FindCompiler("protoc-gen-go")
	fmt.Printf("Using: %s\n", pcfname)
	cmd := []string{
		cmdline.GetYACloudDir() + "/ctools/dev/go/current/protoc/protoc",
		fmt.Sprintf("--plugin=protoc-gen-go=%s", pcfname),
		fmt.Sprintf("--go_out=plugins=grpc:%s", targetdir),
	}
	for _, id := range import_dirs {
		cmd = append(cmd, fmt.Sprintf("-I%s", id))
	}
	for _, pf := range files {
		f := pf.GetFilename()
		Debugf("Compiler working dir: %s, compiling %s\n", dir, f)
		cmdandfile := append(cmd, f)
		l := linux.New()
		out, err := l.SafelyExecuteWithDir(cmdandfile, dir, nil)
		if err != nil {
			cr.AddFailed(pf, err)
			fmt.Printf("Compiler working dir: %s\n", dir)
			fmt.Printf("%s Failed to compile: %s: %s\n", pcfname, f, err)
			fmt.Printf("Compiler output: %s\n", out)
		}
	}
	/***************************** compile create.go ******************************/
	pcfname = cc.FindCompiler("protoc-gen-cnw")
	fmt.Printf("Using: %s\n", pcfname)
	cmd = []string{
		cmdline.GetYACloudDir() + "/ctools/dev/go/current/protoc/protoc",
		fmt.Sprintf("--plugin=protoc-gen-cnw=%s", pcfname),
		fmt.Sprintf("--cnw_out=%s", targetdir),
	}
	for _, id := range import_dirs {
		cmd = append(cmd, fmt.Sprintf("-I%s", id))
	}

	//	fmt.Printf("Compiler working dir: %s\n", dir)
	for _, pf := range files {
		f := pf.GetFilename()
		//		fmt.Printf("Compiler working dir: %s, compiling %s\n", dir, f)
		cmdandfile := append(cmd, f)
		l := linux.New()
		out, err := l.SafelyExecuteWithDir(cmdandfile, dir, nil)
		if err != nil {
			cr.AddFailed(pf, err)
			fmt.Printf("%s Failed to compile: %s: %s\n", pcfname, f, err)
			fmt.Printf("Compiler output: %s\n", out)
		}
	}

	fmt.Printf("Compiling go completed\n")
	return nil
}
