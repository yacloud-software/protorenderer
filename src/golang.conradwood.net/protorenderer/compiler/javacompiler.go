package compiler

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	pr "golang.conradwood.net/apis/protorenderer"
	"golang.conradwood.net/go-easyops/linux"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/protorenderer/common"
	"golang.conradwood.net/protorenderer/compiler/java"
	"golang.conradwood.net/protorenderer/filelayouter"
	//	"golang.conradwood.net/protorenderer/meta"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

var (
	java_compiler_bin   = flag.String("java_compiler", "/usr/bin/javac", "path to javac binary")
	java_release        = flag.String("java_release", "11", "flag --target [java_release] for javac: build for specific target version")
	java_use_std_protoc = flag.Bool("java_std_protoc", true, "if set use standard protoc compiler (the one with the OS rather than shipped in this repo")
	java_plugin_name    = flag.String("java_plugin_name", "protoc-gen-grpc-java-1.13.1-linux-x86_64.exe", "the name of the java gprc plugin in extra/compilers")
)

/*
java compiles in two stages:
1) .proto to .java
2) .java to .class
*/
type JavaCompiler struct {
	WorkDir     string
	javaSrc     string
	javaClasses string
	protofiles  []string
	err         error
	fl          *filelayouter.FileLayouter
	stage       string
	lastVersion int
	cc          CompilerCallback
}

func NewJavaCompiler(cc CompilerCallback) Compiler {
	res := &JavaCompiler{
		cc:      cc,
		fl:      cc.GetFileLayouter(),
		WorkDir: cc.GetFileLayouter().TopDir() + "build",
	}
	return res
}
func (j *JavaCompiler) SetStage(s string) {
	j.stage = s
}
func (j *JavaCompiler) Error() error {
	return j.err
}
func (g *JavaCompiler) Name() string { return "java" }
func (j *JavaCompiler) Compile(rt ResultTracker) error {
	dir := j.fl.SrcDir()
	j.err = nil
	j.SetStage("prepare")
	pfiles, err := AllProtos(dir)
	if err != nil {
		j.err = j.Errorf("allprotos", err)
		return j.err
	}
	j.protofiles = pfiles
	j.Printf("Start java (.proto->.java) compilation...\n")
	l := linux.New()
	j.javaSrc = j.WorkDir + "/java/src"
	j.javaClasses = j.WorkDir + "/java/classes"
	err = common.RecreateSafely(j.javaSrc)
	if err != nil {
		j.err = j.Errorf("recreate-javasrc", err)
		return j.err
	}
	err = common.RecreateSafely(j.javaClasses)
	if err != nil {
		j.err = j.Errorf("recreate-javaclass", err)
		return j.err
	}

	compiler_exe := "/opt/cnw/ctools/dev/go/current/protoc/protoc"

	if !*java_use_std_protoc {
		// find the non-std compiler from this repo:
		compiler_exe, err = utils.FindFile(fmt.Sprintf("extra/compilers/%s/weird.exe", common.GetCompilerVersion()))
		if err != nil {
			j.err = j.Errorf("findweird.exe", err)
			return j.err
		}
		compiler_exe, err = filepath.Abs(compiler_exe)
		if err != nil {
			j.err = j.Errorf("absweirdexe", err)
			return j.err
		}
	}

	// find the compiler plugin:
	plugin_exe, err := utils.FindFile(fmt.Sprintf("extra/compilers/%s/%s", common.GetCompilerVersion(), *java_plugin_name))
	if err != nil {
		j.err = j.Errorf("findgengrpcplugin", err)
		return j.err
	}
	plugin_exe, err = filepath.Abs(plugin_exe)
	if err != nil {
		j.err = j.Errorf("absgengrpcplugin", err)
		return j.err
	}
	j.SetStage("gen-protobuf")

	cmd := []string{
		compiler_exe,
		"--proto_path=.",
		fmt.Sprintf("--java_out=%s", j.javaSrc),
	}
	cmdandfile := append(cmd, j.protofiles...)
	out, err := l.SafelyExecuteWithDir(cmdandfile, dir, nil)
	if err != nil {
		j.Printf("Java compilation failed: %s\n%s\n", out, err)
		j.err = j.Errorf("gen_java_protobuf", err)
		return j.err
	}
	j.SetStage("gen-grpc")

	cmd = []string{
		compiler_exe,
		"--proto_path=.",
		fmt.Sprintf("--plugin=protoc-gen-grpc-java=%s", plugin_exe),
		fmt.Sprintf("--grpc-java_out=%s", j.javaSrc),
	}
	cmdandfile = append(cmd, j.protofiles...)
	out, err = l.SafelyExecuteWithDir(cmdandfile, dir, nil)
	if err != nil {
		j.Printf("Java proto compilation failed: %s\n%s\n", out, err)
		j.err = j.Errorf("gen_java_stubs", err)
		return j.err
	}
	j.SetStage("addextra")
	err = j.AddExtra()
	if err != nil {
		return err
	}
	j.SetStage("compileclass")
	j.Printf("Start java (.java->.class) compilation...\n")
	err = j.CompileToClass()
	if err != nil {
		return err
	}
	j.Printf("Compiling java completed\n")
	return nil
}

/*
* this compiles .java to .class
* this turns out to be difficult to do without creating a
* memory requirement dependency on amount of total classes.
* thus it does incremental compilation, one file at a time
* not efficient, but safe and reliable.
 */
func (j *JavaCompiler) CompileToClass() error {
	compiling_version := j.fl.CurrentVersionNumber()
	jars, err := utils.FindFile(fmt.Sprintf("extra/%s/jars", common.GetCompilerVersion()))
	if err != nil {
		return j.Errorf("findjars", err)
	}
	jars, err = filepath.Abs(jars)
	if err != nil {
		return j.Errorf("absjars", err)
	}
	jarfiles, err := AllFiles(jars, ".jar")
	if err != nil {
		return j.Errorf("alljars", err)
	}
	deli := ""
	classpath := ""
	for _, j := range jarfiles {
		classpath = classpath + deli + jars + "/" + j
		deli = ":"
	}
	classpath = classpath + deli + j.javaClasses
	// well kept secret apparently: classpath needs to include current dir
	// it is unclear why (link to documentation anyone?)
	// either way, without current directory it doesn't work.
	classpath = classpath + deli + "."
	//	j.Printf("Classpath: %s\n", classpath)
	l := linux.New()
	l.SetMaxRuntime(time.Duration(1200) * time.Second)
	cmd := []string{
		*java_compiler_bin,
		"-Xlint:none",
		"-Xdoclint:none",
		/*
			"-J-Xmx200000000", // "heap space"
			"-J-Xms5000000",
			"-J-Xss1000000",
		*/
		"-d",
		j.javaClasses,
		"-cp",
		classpath,
		"-encoding",
		"UTF-8",
		"-target",
		*java_release,
		"-source",
		*java_release,
	}
	// return filenames RELATIVE to j.javaSrc
	files, err := j.findChangedJavaFiles()
	if err != nil {
		return j.Errorf("findchangedjavafiles", err)
	}
	// make filenames ABSOLUTE
	var jfiles []*java.JavaFile
	for _, f := range files {
		jfiles = append(jfiles, &java.JavaFile{
			Absolute: fmt.Sprintf("%s/%s", j.javaSrc, f),
			Relative: f,
		},
		)
	}
	for _, f := range jfiles {
		if !utils.FileExists(f.Absolute) {
			fmt.Printf("SourceDir: %s\n", j.javaSrc)
			fmt.Printf("WorkDir:   %s\n", j.WorkDir)
			fmt.Printf("TopDir:    %s\n", j.fl.TopDir())
			panic(fmt.Sprintf("File does not exist: %s\n", f.Absolute))
		}
	}

	j2c := &java.Java2Class{
		SourceDir: j.javaSrc,
		JFiles:    jfiles,
		Command:   cmd,
	}
	err = j2c.Compile()
	if err != nil {
		return err
	}
	j.lastVersion = compiling_version
	return nil
}

/******************************************
* find all .java files whose corresponding .proto files changed
* (relative to javacompiler.javaSrc
******************************************/
func (j *JavaCompiler) findChangedJavaFiles() ([]string, error) {
	dedup := make(map[string]bool)
	protos := j.fl.ChangedProtos(j.lastVersion)
	j.Printf("compiled version: %d, filelayouter version: %d, changed files: %d\n", j.lastVersion, j.fl.CurrentVersionNumber(), len(protos))
	for _, tc := range protos {
		jfs, err := j.protoFileToJavaFile(tc.Protofile())
		if err != nil {
			return nil, err
		}
		j.Printf("Changed: %s\n", tc.Protofile().JavaPackage)
		for _, jf := range jfs {
			dedup[jf] = true
			j.Printf("    %s\n", jf)
		}
	}

	res := make([]string, len(dedup))
	i := 0
	for k, _ := range dedup {
		res[i] = k
		i++
	}
	return res, nil
}

// return filenames(!) of all .java file for this protofile
// (relative to javaCompiler.javaSrc)
func (j *JavaCompiler) protoFileToJavaFile(pf *pr.ProtoFile) ([]string, error) {
	jp := pf.JavaPackage
	xs := strings.ReplaceAll(jp, ".", "/")
	dir := j.javaSrc + "/" + xs
	st, err := AllJava(dir)
	if err != nil {
		return nil, fmt.Errorf("protoFileToJavaFile failed for %s (javapackage %s): %s", pf.Filename, jp, err)
	}
	var res []string
	for _, f := range st {
		res = append(res, xs+"/"+f)
	}
	return res, err
}
func (j *JavaCompiler) Printf(format string, args ...interface{}) {
	prefix := fmt.Sprintf("[java stage:%s] ", j.stage)
	f := prefix + format
	fmt.Printf(f, args...)
}
func (j *JavaCompiler) Errorf(text string, err error) error {
	prefix := fmt.Sprintf("[java stage:%s] ", j.stage)
	if len(text) > 1500 {
		text = text[:1500] + "..."
	}
	es := fmt.Sprintf("%s", err)
	if len(es) > 1500 {
		es = es[:1500] + "..."
	}
	f := prefix + fmt.Sprintf("while executing \"%s\", an error has occured: %s", text, es)
	return fmt.Errorf(f)
}
func (j *JavaCompiler) Files(ctx context.Context, pkg *pr.Package, filetype string) ([]File, error) {
	var res []File
	javadirname := goPackagenameToJavaDir(pkg.Prefix)

	if filetype == ".class" {
		d := j.javaClasses + "/" + javadirname
		//		fmt.Printf("Finding .class files in \"%s\"\n", d)
		fs, err := AllFiles(d, ".class")
		if err != nil {
			return nil, err
		}
		for _, f := range fs {
			res = append(res, &StdFile{
				Filename:   javadirname + "/" + f,
				ctFilename: d + "/" + f,
				version:    j.lastVersion,
			})
		}

	} else {
		return nil, fmt.Errorf("java compiler does not provide files of type \"%s\"", filetype)
	}
	return res, nil
}
func (j *JavaCompiler) GetFile(ctx context.Context, filename string) (File, error) {
	ext := filepath.Ext(filename)
	ff := filename
	if ext == ".class" {
		ff = j.javaClasses + "/" + ff
	} else if ext == ".java" {
		ff = j.javaSrc + "/" + ff
	}
	s := &StdFile{Filename: filename, ctFilename: ff, version: j.lastVersion}
	return s, nil
}

// convert something "normal", like "golang.singingcat.net/apis" into java style, e.g.
// net/singingcat/golang/apis"
func goPackagenameToJavaDir(gopkg string) string {
	reverse := gopkg
	dirs := ""
	i := strings.Index(gopkg, "/")
	if i != -1 {
		reverse = gopkg[:i]
		dirs = gopkg[i:]
	}
	revs := strings.Split(reverse, ".")
	last := len(revs) - 1
	for i := 0; i < len(revs)/2; i++ {
		revs[i], revs[last-i] = revs[last-i], revs[i]
	}
	reverse = strings.Join(revs, "/")
	return reverse + dirs
}

/***************************************************************************************************
** add the 'create' helpers e.g. CreateXXXClient()
***************************************************************************************************/

/*
 Example:
        Context ctx = Token.ContextWithToken();
        HelloWorldGrpc.HelloWorldBlockingStub hsb;
        io.grpc.ManagedChannel mc = ConnectionManager.Connect("helloworld.HelloWorld");
        hsb = HelloWorldGrpc.newBlockingStub(mc)
                .withCallCredentials(ConnectionManager.getCallCredentials(ctx));


*/

// we create "extra" stuff, such as CreateXXXClient()
func (j *JavaCompiler) AddExtra() error {
	t := `
package {{.Protofile.JavaPackage}};

import net.conradwood.javaeasyops.client.ConnectionManager;

public class {{.HelperClassName}} {
public static final String PackageName = "{{.Protofile.JavaPackage}}";
public static final String GoPackageName = "{{.Protofile.GoPackage}}";
public static final String Filename = "{{.Protofile.Filename}}";
{{range .Services}}
private static io.grpc.ManagedChannel {{.GRPCClientName}}_Channel = null;
private static Object {{.GRPCClientName}}_Lock = new Object();
{{end}}

{{range .Services}}
// create helper for {{.RegistryName}}
public static {{.GRPCClientName}}.{{.Name}}BlockingStub Get{{.Name}}() throws Exception {
    return Get{{.Name}}(io.grpc.Context.current());
}
public static {{.GRPCClientName}}.{{.Name}}BlockingStub Get{{.Name}}(io.grpc.Context ctx) throws Exception {
  io.grpc.ManagedChannel mc = {{.GRPCClientName}}_Channel;
  if (mc == null) {
     synchronized({{.GRPCClientName}}_Lock) {
     mc = {{.GRPCClientName}}_Channel;
        if (mc == null) { 
          mc = ConnectionManager.Connect("{{.RegistryName}}");
          {{.GRPCClientName}}_Channel = mc;
        }
     }
  }
  {{.GRPCClientName}}.{{.Name}}BlockingStub gp;
  gp = {{.GRPCClientName}}.newBlockingStub(mc).withCallCredentials(ConnectionManager.getCallCredentials(ctx));
  return gp;
}
{{end}}
}
`
	tm := template.New("javaextra")
	_, err := tm.Parse(t)
	if err != nil {
		return err
	}
	files, err := j.findServices()
	if err != nil {
		return err
	}
	j.Printf("adding 'extra' java classes for %d files\n", len(files))
	for _, ps := range files {
		ps.HelperClassName = "Helper"
		if ps.Protofile.Meta == nil {
			j.Printf("Protofile: %s has no meta information\n", ps.javadir)
			continue
		}
		pid := ps.Protofile.Meta.PackageID
		j.Printf("package #%s - javadir: %s\n", pid, ps.javadir)
		if ps.Protofile.JavaPackage == "1" {
			panic("java package cannot be '1'")
		}
		mpkg := j.cc.GetMetaPackageByID(pid)
		if mpkg == nil {
			fmt.Printf("No meta package for %s\n", pid)
			continue
		}
		for _, s := range mpkg.Services {
			j.Printf("    service: \"%s\"\n", s.Name)
			ps.Services = append(ps.Services, &servicedef{
				Name:           s.Name,
				RegistryName:   mpkg.Name + "." + s.Name,
				GRPCClientName: s.Name + "Grpc",
			})
		}
		s := &bytes.Buffer{}
		err := tm.Execute(s, ps)
		if err != nil {
			return err
		}
		fname := j.javaSrc + "/" + ps.javadir + "/" + ps.HelperClassName + ".java"
		err = utils.WriteFile(fname, s.Bytes())
		if err != nil {
			return err
		}
	}
	return nil
}

// get all protos and services defined in each
func (j *JavaCompiler) findServices() ([]*protoservice, error) {
	var res []*protoservice
	protos := j.fl.ChangedProtos(j.lastVersion)
	j.Printf("compiled version: %d, filelayouter version: %d, changed files: %d\n", j.lastVersion, j.fl.CurrentVersionNumber(), len(protos))
	for _, tc := range protos {
		pf := tc.Protofile()
		jf, err := j.protoFileToJavaFile(pf)
		if err != nil {
			return nil, err
		}
		ps := &protoservice{Protofile: pf}
		if len(jf) != 0 {
			ps.javadir = filepath.Dir(jf[0])
		}
		/*

		 */
		res = append(res, ps)
	}
	return res, nil
}

type protoservice struct {
	Protofile       *pr.ProtoFile
	javadir         string
	Services        []*servicedef
	HelperClassName string
}
type servicedef struct {
	Name           string
	RegistryName   string
	GRPCClientName string
}
















































































