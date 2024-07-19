package java

import (
	"context"
	"fmt"
	gocmdline "golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/linux"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/protorenderer/cmdline"
	"golang.conradwood.net/protorenderer/v2/helpers"
	"golang.conradwood.net/protorenderer/v2/interfaces"
	"path/filepath"
	"strings"
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
	cp, _ := classpath()
	fmt.Printf("CLASSPATH=%s\n", strings.Join(cp, ":"))
	return &JavaCompiler{}
}

func (gc JavaCompiler) ShortName() string { return "java" }

// compiles new .proto files to .java and to .class
// assumes that _all_ existing .proto files are already compiled to .java and .class files
// (otherwise the imports won't work)
func (gc *JavaCompiler) Compile(ctx context.Context, ce interfaces.CompilerEnvironment, files []interfaces.ProtoFile, outdir string, cr interfaces.CompileResult) error {
	dir := ce.WorkDir() + "/" + ce.NewProtosDir()
	targetdir := outdir + "/" + gc.ShortName()
	err := helpers.Mkdir(targetdir)
	if err != nil {
		return err
	}
	import_dirs := []string{
		dir,
		ce.AllKnownProtosDir(),
	}

	var proto_file_names []string
	for _, pf := range files {
		proto_file_names = append(proto_file_names, pf.GetFilename())
	}

	fmt.Printf("Start java (.proto->.java) compilation...\n")
	l := linux.New()
	//	j.javaSrc = j.WorkDir + "/java/src"
	javaSrc := targetdir + "/src"
	err = helpers.Mkdir(javaSrc)
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
	nf := helpers.NewFileFinder(javaSrc)
	nf.FindNew() // build internal list
	fmt.Printf("Compiling files: \"%s\"\n", strings.Join(proto_file_names, " "))
	cmdandfile := append(cmd, proto_file_names...)
	out, err := l.SafelyExecuteWithDir(cmdandfile, dir, nil)
	if err != nil {
		fmt.Printf("Output:\n%s\n", out)
		return err
	}

	// inject the custom headers to the .java files
	new_files, err := nf.FindNew()
	if err != nil {
		return fmt.Errorf("(1) unable to find new files: %w", err)
	}

	for _, new_file := range new_files {
		fm := helpers.NewFileModifierFromFilename(javaSrc + "/" + new_file)
		fm.AddHeader(fmt.Sprintf("// created by protorenderer, run (sources: %s)\n", strings.Join(proto_file_names, " ")))
		err = fm.Save()
		if err != nil {
			return fmt.Errorf("unable to save modified file: %w", err)
		}
	}

	err = gc.compileJava2Class(ctx, ce, files, outdir, cr)
	if err != nil {
		return err
	}
	return nil
}

// this compiles .java files to .class
func (gc *JavaCompiler) compileJava2Class(ctx context.Context, ce interfaces.CompilerEnvironment, files []interfaces.ProtoFile, outdir string, cr interfaces.CompileResult) error {
	dir := outdir + "/" + gc.ShortName() + "/src"           // this is where the .java files are
	targetdir := outdir + "/" + gc.ShortName() + "/classes" // this where the .class files go
	err := helpers.Mkdir(targetdir)
	if err != nil {
		return err
	}

	java_files, err := helpers.FindFiles(dir, ".java")
	if err != nil {
		return err
	}

	fmt.Printf("Start java (.java->.class) compilation (src=%s,target=%s)...\n", dir, targetdir)
	javac := "/etc/java-home/bin/javac"
	cmd := []string{
		javac,
		"-Xlint:none",
		"-Xdoclint:none",
		/*
			"-J-Xmx200000000", // "heap space"
			"-J-Xms5000000",
			"-J-Xss1000000",
		*/
		"--source-path", "/tmp/pr/v2/compile_outdir/java/src",
		"-d",
		targetdir,
		//		"-cp",
		//		classpath,
		"-encoding",
		"UTF-8",
		"-target",
		cmdline.GetJavaRelease(),
		"-source",
		cmdline.GetJavaRelease(),
	}
	cp, err := classpath()
	if err != nil {
		return err
	}
	cmdandfile := append(cmd, java_files...)
	l := linux.New()
	env := []string{
		"CLASSPATH=" + strings.Join(cp, ":"),
	}
	l.SetEnvironment(env)
	out, err := l.SafelyExecuteWithDir(cmdandfile, dir, nil)
	if err != nil {
		fmt.Printf(".java->.class failed: %s\n", string(out))
		return err
	}
	return nil
}

func classpath() ([]string, error) {
	jars, err := utils.FindFile(fmt.Sprintf("extra/jars"))
	//	jars, err := utils.FindFile(fmt.Sprintf("extra/%s/jars", cmdline.GetCompilerVersion()))
	if err != nil {
		return nil, err
	}
	jars, err = filepath.Abs(jars)
	if err != nil {
		return nil, err
	}
	jarfiles, err := helpers.FindFiles(jars, ".jar")
	if err != nil {
		return nil, err
	}
	var res []string
	for _, j := range jarfiles {
		jf := jars + "/" + j
		res = append(res, jf)
	}
	return res, nil
}
