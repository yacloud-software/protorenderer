package meta

import (
	"context"
	"fmt"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/linux"
	cc "golang.conradwood.net/protorenderer/v2/compilers/common"
	"golang.conradwood.net/protorenderer/v2/interfaces"
	"time"
)

type MetaCompiler struct {
}

func New() *MetaCompiler {
	return &MetaCompiler{}
}

/*
the meta compilers works slightly different than the others. the protoc plugin is a small RPC stub, which then calls protorenderer
*/
func (gc *MetaCompiler) Compile(ctx context.Context, token string, port int, ce interfaces.CompilerEnvironment, files []interfaces.ProtoFile, outdir string, cr interfaces.CompileResult) error {
	pcfname := cc.FindCompiler("protoc-gen-meta")
	dir := ce.WorkDir() + "/" + ce.NewProtosDir()

	import_dirs := []string{
		dir,
		ce.WorkDir() + "/" + ce.AllKnownProtosDir(),
	}

	l := linux.New()
	l.SetMaxRuntime(time.Duration(600) * time.Second)
	sctx, err := auth.SerialiseContextToString(ctx)
	if err != nil {
		fmt.Printf("Meta-Compiler: Unable to serialise context: %s\n", err)
		return err
	}

	cmd := []string{
		cmdline.GetYACloudDir() + "/ctools/dev/go/current/protoc/protoc",
		fmt.Sprintf("--plugin=protoc-gen-meta=%s", pcfname),
		"--meta_out=/tmp", // has no output
		fmt.Sprintf("--meta_opt=%s,%s,%d,%s", token, sctx, port, cmdline.GetClientRegistryAddress()),
	}
	for _, id := range import_dirs {
		cmd = append(cmd, fmt.Sprintf("-I%s", id))
	}
	for _, pf := range files {
		filename := pf.GetFilename()
		cmdfl := append(cmd, filename)

		out, err := l.SafelyExecuteWithDir(cmdfl, dir, nil)
		if err != nil {
			fmt.Printf("protoc output: %s\n", out)
			fmt.Printf("Failed to compile: %s\n", err)
			continue
		}
	}
	return nil
}
