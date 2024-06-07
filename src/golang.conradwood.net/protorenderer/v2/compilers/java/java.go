package java

import (
	"context"
	//	"flag"
	"fmt"
	gocmdline "golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/linux"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/protorenderer/cmdline"
	"golang.conradwood.net/protorenderer/v2/helpers"
	"golang.conradwood.net/protorenderer/v2/interfaces"
	"path/filepath"
)

var (
// java_compiler_bin   = flag.String("java_compiler", "/usr/bin/javac", "path to javac binary")
// java_release        = flag.String("java_release", "11", "flag --target [java_release] for javac: build for specific target version")
// java_use_std_protoc = flag.Bool("java_std_protoc", true, "if set use standard protoc compiler (the one with the OS rather than shipped in this repo")
// java_plugin_name = flag.String("java_plugin_name", "protoc-gen-grpc-java-1.13.1-linux-x86_64.exe", "the name of the java gprc plugin in extra/compilers")
)

/*
java compiles in two stages:
1) .proto to .java
2) .java to .class
*/
type JavaCompiler struct {
	stage string
	/*
		WorkDir     string
		javaSrc     string
		javaClasses string
		protofiles  []string
		err         error
		fl          *filelayouter.FileLayouter
		stage       string
		lastVersion int
			cc          CompilerCallback
	*/
}

func New() interfaces.Compiler {
	return &JavaCompiler{}
}
func (gc JavaCompiler) ShortName() string { return "java" }
func (gc *JavaCompiler) Compile(ctx context.Context, ce interfaces.CompilerEnvironment, files []interfaces.ProtoFile, outdir string, cr interfaces.CompileResult) error {
	dir := ce.WorkDir() + "/" + ce.NewProtosDir()
	targetdir := outdir + "/" + gc.ShortName()
	err := helpers.Mkdir(targetdir)
	if err != nil {
		return err
	}
	import_dirs := []string{
		dir,
		ce.WorkDir() + "/" + ce.AllKnownProtosDir(),
	}

	var proto_file_names []string
	for _, pf := range files {
		proto_file_names = append(proto_file_names, pf.GetFilename())
	}
	fmt.Printf("Start java (.proto->.java) compilation...\n")
	l := linux.New()
	//	j.javaSrc = j.WorkDir + "/java/src"
	//	j.javaClasses = j.WorkDir + "/java/classes"
	javaSrc := targetdir + "/src"
	javaclasses := targetdir + "/classes"
	err = helpers.Mkdir(javaSrc)
	if err != nil {
		return err
	}
	err = helpers.Mkdir(javaclasses)
	if err != nil {
		return err
	}
	compiler_exe := gocmdline.GetYACloudDir() + "/ctools/dev/go/current/protoc/protoc"

	if !cmdline.GetJavaStdProtoC() {
		// find the non-std compiler from this repo:
		compiler_exe, err = utils.FindFile(fmt.Sprintf("extra/compilers/%s/weird.exe", cmdline.GetCompilerVersion()))
		if err != nil {
			return err
		}
		compiler_exe, err = filepath.Abs(compiler_exe)
		if err != nil {
			return err
		}
	}

	// find the compiler plugin:
	plugin_exe, err := utils.FindFile(fmt.Sprintf("extra/compilers/%s/%s", cmdline.GetCompilerVersion(), cmdline.GetJavaPluginName()))
	if err != nil {
		return err
	}
	plugin_exe, err = filepath.Abs(plugin_exe)
	if err != nil {
		return err
	}
	//	j.SetStage("gen-protobuf")

	cmd := []string{
		compiler_exe,
		"--proto_path=.",
		fmt.Sprintf("--java_out=%s", javaSrc),
	}
	for _, id := range import_dirs {
		cmd = append(cmd, fmt.Sprintf("-I%s", id))
	}
	cmdandfile := append(cmd, proto_file_names...)
	out, err := l.SafelyExecuteWithDir(cmdandfile, dir, nil)
	if err != nil {
		fmt.Printf("Output:\n%s\n", out)
		return err
	}
	/*
	   	cmd = []string{
	   		compiler_exe,
	   		"--proto_path=.",
	   		fmt.Sprintf("--plugin=protoc-gen-grpc-java=%s", plugin_exe),
	   		fmt.Sprintf("--grpc-java_out=%s", javaSrc),
	   	}

	   cmdandfile = append(cmd, j.protofiles...)
	   out, err = l.SafelyExecuteWithDir(cmdandfile, dir, nil)

	   	if err != nil {
	   		j.Printf("Java proto compilation failed: %s\n%s\n", out, err)
	   		j.err = j.Errorf("gen_java_stubs", err)
	   		return j.err
	   	}
	*/
	return nil
}
