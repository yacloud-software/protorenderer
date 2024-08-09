package java

import (
	"context"
	"fmt"
	gocmdline "golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/linux"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/protorenderer/cmdline"
	"golang.conradwood.net/protorenderer/v2/helpers"
	"golang.conradwood.net/protorenderer/v2/interfaces"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	MARKER = "// this marker placed here for parsing. SOURCE:"
)

var (
	numeric_regex = regexp.MustCompile(`\.\d`)
)

/*
java compiles in two stages:
1) .proto to .java
2) .java to .class
*/
type JavaCompiler struct {
	compiled_to_class map[string]bool
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
	//	cp, _ := classpath()
	//	fmt.Printf("CLASSPATH=%s\n", strings.Join(cp, ":"))
	return &JavaCompiler{compiled_to_class: make(map[string]bool)}
}

func (gc JavaCompiler) ShortName() string { return "java" }

// compiles new .proto files to .java and to .class
// assumes that _all_ existing .proto files are already compiled to .java and .class files
// (otherwise the imports won't work)
func (gc *JavaCompiler) Compile(ctx context.Context, ce interfaces.CompilerEnvironment, files []interfaces.ProtoFile, outdir string, cr interfaces.CompileResult) error {

	x := len(files)
	files = gc.filter_known_non_working(ce, files, cr)
	if len(files) == 0 {
		fmt.Printf("WARNING: javacompiler invoked with %d files, but after filtering broken ones, none left\n", x)
		//		return errors.Errorf("javacompiler invoked with %d files, but after filtering broken ones, none left", x)
		return nil
	}
	dir := ce.NewProtosDir()
	targetdir := outdir + "/" + gc.ShortName()
	err := helpers.Mkdir(targetdir)
	if err != nil {
		return errors.Wrap(err)
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
		return errors.Wrap(err)
	}

	compiler_exe := gocmdline.GetYACloudDir() + "/ctools/dev/go/current/protoc/protoc"

	if !cmdline.GetJavaStdProtoC() {
		// find the non-std compiler from this repo:
		compiler_exe, err = utils.FindFile(fmt.Sprintf("extra/compilers/%s/weird.exe", cmdline.GetCompilerVersion()))
		if err != nil {
			return errors.Wrap(err)
		}
		compiler_exe, err = filepath.Abs(compiler_exe)
		if err != nil {
			return errors.Wrap(err)
		}
	}

	// find the compiler plugin:
	plugin_exe, err := utils.FindFile(fmt.Sprintf("extra/compilers/%s/%s", cmdline.GetCompilerVersion(), cmdline.GetJavaPluginName()))
	if err != nil {
		return errors.Wrap(err)
	}
	plugin_exe, err = filepath.Abs(plugin_exe)
	if err != nil {
		return errors.Wrap(err)
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

	var successful_files []interfaces.ProtoFile
	for _, pfn := range files {
		fname := pfn.GetFilename()
		fmt.Printf("[javacompiler] Compiling file: \"%s\"\n", fname)
		cmdandfile := append(cmd, fname)
		out, err := l.SafelyExecuteWithDir(cmdandfile, dir, nil)
		if err != nil {
			cr.AddFailed(gc, pfn, err, []byte(out))
		} else {
			successful_files = append(successful_files, pfn)
		}
		new_files, err := nf.FindNew() // to inject the custom headers to the .java files
		if err != nil {
			return errors.Errorf("(1) unable to find new files: %w", err)
		}

		for _, new_file := range new_files {
			fm := helpers.NewFileModifierFromFilename(javaSrc + "/" + new_file)
			fm.AddHeader(fmt.Sprintf("%s%s\n", MARKER, fname))
			err = fm.Save()
			if err != nil {
				return errors.Errorf("unable to save modified file: %w", err)
			}
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
	err = gc.compileJava2Class(ctx, ce, successful_files, outdir, cr)
	if err != nil {
		return errors.Wrap(err)
	}
	return nil
}

// this compiles .java files to .class
func (gc *JavaCompiler) compileJava2Class(ctx context.Context, ce interfaces.CompilerEnvironment, files []interfaces.ProtoFile, outdir string, cr interfaces.CompileResult) error {
	dir := outdir + "/" + gc.ShortName() + "/src"           // this is where the .java files are
	targetdir := outdir + "/" + gc.ShortName() + "/classes" // this where the .class files go
	err := helpers.Mkdir(targetdir)
	if err != nil {
		return errors.Wrap(err)
	}

	java_files, err := helpers.FindFiles(dir, ".java")
	if err != nil {
		return errors.Wrap(err)
	}

	fmt.Printf("Start java (.java->.class) compilation (src=%s,target=%s) (%d files, submitted %d files)...\n", dir, targetdir, len(java_files), len(files))
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
		"--source-path", fmt.Sprintf("%s/java/src", ce.CompilerOutDir()),
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
	cp, err := classpath(ce)
	if err != nil {
		return errors.Wrap(err)
	}
	sort.Slice(java_files, func(i, j int) bool {
		return java_files[i] < java_files[j]
	})

	// filter those that were done already
	var njava_files []string
	for _, nf := range java_files {
		if gc.compiled_to_class[nf] {
			continue
		}
		njava_files = append(njava_files, nf)
	}
	java_files = njava_files
	if len(java_files) == 0 {
		return nil
	}
	/*
		at first, we are optimistic, try compiling all at once. this is the "good path". if none fails, it is the fastest means to compile
		if the compiler fails, we go through each file to find out which one failed
	*/
	fmt.Printf("[javacompiler] Compiling %d java files to class..\n", len(java_files))
	cmdandfile := append(cmd, java_files...)
	l := linux.New()
	env := []string{
		"CLASSPATH=" + strings.Join(cp, ":"),
	}
	l.SetEnvironment(env)
	l.SetMaxRuntime(time.Duration(60) * time.Second)
	_, err = l.SafelyExecuteWithDir(cmdandfile, dir, nil)
	if err == nil {
		// none failed
		for _, jf := range java_files {
			gc.compiled_to_class[jf] = true
		}
		return nil
	}
	// try each file
	for i, java_file := range java_files {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		source_proto_name, err := gc.get_source_from_java_file(dir + "/" + java_file)
		if err != nil {
			return errors.Wrap(err)
		}

		fmt.Printf("[javacompiler] Compiling %s->%s->class file (%d of %d)\n", source_proto_name, java_file, i+1, len(java_files))
		cmdandfile := append(cmd, java_file)
		l := linux.New()
		env := []string{
			"CLASSPATH=" + strings.Join(cp, ":"),
		}
		l.SetEnvironment(env)
		l.SetMaxRuntime(time.Duration(60) * time.Second)
		out, err := l.SafelyExecuteWithDir(cmdandfile, dir, nil)
		if err != nil {
			fmt.Printf("%s->.class failed: %s\n", java_file, string(out))
			got_pf := false
			for _, pf := range files {
				if pf.GetFilename() == source_proto_name {
					got_pf = true
					cr.AddFailed(gc, pf, fmt.Errorf("failed to compile .class: %s", err), []byte(out))
				}
			}
			if !got_pf {
				return errors.Errorf("file %s: compiling %s to .class, was unable to find protofile", source_proto_name, java_file)
			}
			continue
		}
		gc.compiled_to_class[java_file] = true
	}
	return nil
}

func classpath(ce interfaces.CompilerEnvironment) ([]string, error) {
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
	// order really matters!
	res = append(res, fmt.Sprintf("%s", ce.CompilerOutDir()+"/java/classes"))
	res = append(res, fmt.Sprintf("%s", ce.StoreDir()+"/java/classes"))
	return res, nil
}

func (gc *JavaCompiler) filter_known_non_working(ce interfaces.CompilerEnvironment, files []interfaces.ProtoFile, cr interfaces.CompileResult) []interfaces.ProtoFile {
	var res []interfaces.ProtoFile
	for _, pf := range files {
		pfi := ce.MetaCache().ByProtoFile(pf)
		if pfi == nil {
			// what should we do if we have no metadata here?
			fmt.Printf("[javacompiler] filter: no metadata for %s\n", pf.GetFilename())
			continue
		}
		// check java package name
		java_package := pfi.PackageJava
		if java_package == "" {
			java_package = pfi.Package
		}
		if java_package == "" {
			cr.AddFailed(gc, pf, fmt.Errorf("filter: unable to determine java package name"), nil)
			continue
		}
		if strings.Contains(java_package, " ") {
			cr.AddFailed(gc, pf, fmt.Errorf("filter: Invalid space in Packagename: \"%s\"", java_package), nil)
			continue
		}
		if numeric_regex.MatchString(java_package) {
			cr.AddFailed(gc, pf, fmt.Errorf("filter: Invalid digit in Packagename: \"%s\"", java_package), nil)
			continue
		}
		res = append(res, pf)
	}
	return res
}

// given an absolute filename finds the header //source: and returns the .proto name
func (gc *JavaCompiler) get_source_from_java_file(abs_filename string) (string, error) {
	b, err := utils.ReadFile(abs_filename)
	if err != nil {
		return "", err
	}
	i := 0
	for _, line := range strings.Split(string(b), "\n") {
		i++
		if i > 10 {
			return "", errors.Errorf("no header found in first %d lines of file %s\n", i, abs_filename)
		}
		if strings.Contains(line, MARKER) {
			source := strings.TrimPrefix(line, MARKER)
			source = strings.Trim(source, "\n")
			return source, nil
		}
	}
	return "", errors.Errorf("no header found in file %s\n", abs_filename)
}
