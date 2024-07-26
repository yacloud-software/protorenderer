package meta_compiler

import (
	"context"
	"fmt"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/linux"
	"golang.conradwood.net/go-easyops/utils"
	cc "golang.conradwood.net/protorenderer/v2/compilers/common"
	"golang.conradwood.net/protorenderer/v2/helpers"
	"golang.conradwood.net/protorenderer/v2/interfaces"
	"time"
)

// TODO - instead of using gRPC use go-easyops IPC

type MetaCompiler struct {
	id    string
	ce    interfaces.CompilerEnvironment
	files []interfaces.ProtoFile
	cr    interfaces.CompileResult
}

func New() *MetaCompiler {
	mc := &MetaCompiler{id: utils.RandomString(64)}
	meta_compilers.Put(mc.id, mc)
	return mc
}

/*
the meta compilers works slightly different than the others. the protoc plugin is a small RPC stub, which then calls protorenderer. this function invokes protoc and protoc-gen-meta2 plugin, which then will call this process via gRPC
*/
func (gc *MetaCompiler) Compile(ctx context.Context, port int, ce interfaces.CompilerEnvironment, files []interfaces.ProtoFile, outdir string, cr interfaces.CompileResult) error {
	// keep the variables, we need them for the callback (protoc plugin calls this process via gRPC)
	gc.ce = ce
	gc.files = files
	gc.cr = cr
	helpers.Mkdir(ce.CompilerOutDir() + "/info")
	helpers.Mkdir(ce.StoreDir() + "/info")
	pcfname := cc.FindCompiler("protoc-gen-meta2")
	fmt.Printf("Using compiler: \"%s\"\n", pcfname)
	dir := ce.NewProtosDir()

	import_dirs := []string{
		dir,
		ce.AllKnownProtosDir(),
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
		fmt.Sprintf("--meta2_opt=%s,%s,%d,%s", gc.id, sctx, port, cmdline.GetClientRegistryAddress()),
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

func (mc *MetaCompiler) FileByName(name string) (interfaces.ProtoFile, error) {
	for _, pf := range mc.files {
		if pf.GetFilename() == name {
			return pf, nil
		}
	}
	return nil, fmt.Errorf("File \"%s\" not part of the meta compiler", name)
}
func (mc *MetaCompiler) CompilerEnvironment() interfaces.CompilerEnvironment {
	return mc.ce
}
