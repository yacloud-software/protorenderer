package meta_compiler

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

// TODO - instead of using gRPC use go-easyops IPC

type MetaCompiler struct {
}

func New() *MetaCompiler {
	return &MetaCompiler{}
}

/*
the meta compilers works slightly different than the others. the protoc plugin is a small RPC stub, which then calls protorenderer
*/
func (gc *MetaCompiler) Compile(ctx context.Context, token string, port int, ce interfaces.CompilerEnvironment, files []interfaces.ProtoFile, outdir string, cr interfaces.CompileResult) error {
	pcfname := cc.FindCompiler("protoc-gen-meta2")
	fmt.Printf("Using compiler: \"%s\"\n", pcfname)
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
		fmt.Sprintf("--plugin=protoc-gen-meta2=%s", pcfname),
		"--meta2_out=/tmp", // has no output
		fmt.Sprintf("--meta2_opt=%s,%s,%d,%s", token, sctx, port, cmdline.GetClientRegistryAddress()),
	}
	for _, id := range import_dirs {
		cmd = append(cmd, fmt.Sprintf("-I%s", id))
	}
	for _, pf := range files {
		filename := pf.GetFilename()
		cmdfl := append(cmd, filename)

		out, err := l.SafelyExecuteWithDir(cmdfl, dir, nil)
		if err != nil {
			fmt.Printf("[metacompiler] protoc output: %s\n", out)
			fmt.Printf("[metacompiler] Failed to compile: %s\n", err)
			cr.AddFailed(pf, err, []byte(out))
			continue
		}
	}
	return nil
}
