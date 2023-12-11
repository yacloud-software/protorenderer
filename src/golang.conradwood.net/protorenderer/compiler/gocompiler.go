package compiler

/*
 this compiles the following:
 1. .proto->pb.go
 2. .proto->create_n.go (the create client stubs, see protoc-gen-cnw)
*/
import (
	"context"
	"fmt"
	pr "golang.conradwood.net/apis/protorenderer"
	"golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/linux"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/protorenderer/common"
	"golang.conradwood.net/protorenderer/filelayouter"
	"golang.conradwood.net/protorenderer/protoparser"
	"path/filepath"
	"strings"
	"sync"
)

var ()

type GoCompiler struct {
	lock    sync.Mutex
	WorkDir string
	err     error
	cc      CompilerCallback
	fl      *filelayouter.FileLayouter
}

func NewGoCompiler(cc CompilerCallback) Compiler {
	res := &GoCompiler{
		cc:      cc,
		fl:      cc.GetFileLayouter(),
		WorkDir: cc.GetFileLayouter().TopDir() + "build",
	}
	return res
}
func (g *GoCompiler) Error() error {
	return g.err
}
func (g *GoCompiler) Files(ctx context.Context, pkg *pr.Package, filetype string) ([]File, error) {
	fmt.Printf("Want compiled go files for package: ID=%s, Name=%s, Prefix=%s\n", pkg.ID, pkg.Name, pkg.Prefix)
	ds := g.WorkDir + "/go/" + pkg.Prefix
	fmt.Printf("Want compiled go files for package: ID=%s, Name=%s, Prefix=%s. Looking in dir \"%s\"\n", pkg.ID, pkg.Name, pkg.Prefix, ds)
	fs, err := AllFiles(ds, filetype)
	if err != nil {
		return nil, err
	}
	var res []File
	for _, f := range fs {
		fn := pkg.Prefix + "/" + f
		fl := &StdFile{Filename: fn, version: 1, ctFilename: ds + "/" + f}
		res = append(res, fl)
	}
	return res, nil
}
func (g *GoCompiler) GetFile(ctx context.Context, filename string) (File, error) {
	fn := g.WorkDir + "/go/" + filename
	fl := &StdFile{Filename: filename, version: 1, ctFilename: fn}
	return fl, nil
}
func (g *GoCompiler) Name() string { return "go" }
func (g *GoCompiler) Compile(rt ResultTracker) error {
	dir := g.fl.SrcDir()
	g.lock.Lock()
	defer g.lock.Unlock()
	g.err = nil
	files, err := AllProtos(dir)
	if err != nil {
		fmt.Printf("Failed to compile go files: %s\n", err)
		g.err = err
		return err
	}
	fmt.Printf("Compiling %d .proto files to .pb.go \"%s\"\n", len(files), dir)
	targetdir := g.WorkDir + "/go"
	err = common.RecreateSafely(targetdir)
	if err != nil {
		g.err = err
		return err
	}

	dirfiles, err := newDirFiles(dir, files)
	if err != nil {
		g.err = err
		return err
	}
	/***************************** compile .proto -> .pb.go ******************************/
	pcfname := FindCompiler("protoc-gen-go")
	fmt.Printf("Using: %s\n", pcfname)
	cmd := []string{
		cmdline.GetYACloudDir() + "/ctools/dev/go/current/protoc/protoc",
		fmt.Sprintf("--plugin=protoc-gen-go=%s", pcfname),
		fmt.Sprintf("--go_out=plugins=grpc:%s", targetdir),
	}

	//	fmt.Printf("Compiler working dir: %s\n", dir)
	for _, f := range dirfiles {
		Debugf("Compiler working dir: %s, compiling %s\n", dir, f)
		cmdandfile := append(cmd, f...)
		l := linux.New()
		out, err := l.SafelyExecuteWithDir(cmdandfile, dir, nil)
		if err != nil {
			for _, af := range f {
				rt.AddFailed(g, af, fmt.Sprintf("%s", out))
			}
			fmt.Printf("Failed to compile: %s: %s\n", f, err)
			fmt.Printf("Compiler output: %s\n", out)
			g.err = err
		}
	}
	/***************************** compile create.go ******************************/
	pcfname = FindCompiler("protoc-gen-cnw")
	fmt.Printf("Using: %s\n", pcfname)
	cmd = []string{
		cmdline.GetYACloudDir() + "/ctools/dev/go/current/protoc/protoc",
		fmt.Sprintf("--plugin=protoc-gen-cnw=%s", pcfname),
		fmt.Sprintf("--cnw_out=%s", targetdir),
	}

	//	fmt.Printf("Compiler working dir: %s\n", dir)
	for _, f := range files {
		//		fmt.Printf("Compiler working dir: %s, compiling %s\n", dir, f)
		cmdandfile := append(cmd, f)
		l := linux.New()
		out, err := l.SafelyExecuteWithDir(cmdandfile, dir, nil)
		if err != nil {
			fmt.Printf("Failed to compile: %s: %s\n", f, err)
			fmt.Printf("Compiler output: %s\n", out)
			g.err = err
		}
	}

	fmt.Printf("Compiling go completed\n")
	return nil
}
func FindCompiler(cname string) string {
	check := []string{
		"dist/linux/amd64/",
		fmt.Sprintf("extra/compilers/%s/", common.GetCompilerVersion()),
		"linux/amd64/",
		cmdline.GetYACloudDir() + "/ctools/dev/go/current/go/bin/",
		"/opt/cnw/ctools/dev/go/current/go/bin/",
		"/home/cnw/go/bin/",
		"/home/cnw/devel/go/protorenderer/dist/linux/amd64/",
	}
	var err error
	for _, d := range check {
		c := d + cname
		if !strings.HasPrefix(d, "/") {
			c, err = utils.FindFile(d + cname)
			if err != nil {
				continue
			}
		}

		if !utils.FileExists(c) {
			continue
		}
		if c[0] == '/' {
			return c
		}
		cs, err := filepath.Abs(c)
		if err != nil {
			fmt.Printf("Unable to absolutise \"%s\": %s\n", cs, err)
			return c
		}
		return cs
	}
	e := fmt.Sprintf("%s not found\n", cname)
	fmt.Println(e)
	panic(e)
}

// protofiles per directory
func newDirFiles(dir string, files []string) (map[string][]string, error) {
	res := make(map[string][]string)
	for _, f := range files {
		b, err := utils.ReadFile(dir + "/" + f)
		if err != nil {
			return nil, err
		}
		pp, err := protoparser.Parse(string(b))
		key := filepath.Dir(f) + "_" + pp.GoPackage
		l := append(res[key], f)
		res[key] = l
	}
	return res, nil
}




















































































