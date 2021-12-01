package compiler

import (
	"context"
	"fmt"
	pr "golang.conradwood.net/apis/protorenderer"
	"golang.conradwood.net/go-easyops/linux"
	"golang.conradwood.net/protorenderer/common"
	"golang.conradwood.net/protorenderer/filelayouter"
	"sync"
)

type NanoPBCompiler struct {
	lock sync.Mutex
	err  error
	fl   *filelayouter.FileLayouter
}

func (npb *NanoPBCompiler) WorkDir() string {
	return npb.fl.TopDir() + "build"
}
func (npb *NanoPBCompiler) Compile() error {
	npb.Printf("***************** nanopb compiler ******************\n")
	dir := npb.fl.SrcDir()
	npb.lock.Lock()
	defer npb.lock.Unlock()
	npb.err = nil
	fmt.Printf("Compiling go \"%s\"\n", dir)
	files, err := AllProtos(dir)
	if err != nil {
		npb.err = err
		return err
	}
	l := linux.New()
	l.MyIP()
	targetdir := npb.WorkDir() + "/nanopb"
	err = common.RecreateSafely(targetdir)
	if err != nil {
		npb.err = err
		return err
	}
	npb.Printf("Sourcedir: %s\n", npb.fl.SrcDir())
	npb.Printf("Workdir:   %s\n", npb.WorkDir())
	npb.Printf("Targetdir: %s\n", targetdir)
	for _, f := range files {
		srcname := f
		npb.Printf("File: %s [COMPILING]\n", srcname)
		com := []string{
			"nanopb_generator.py",
			"-D", targetdir, // output dir
			"-Q", `# include "nanopb/%s"`,
			"-L", `# include <nanopb/%s>"`,
			srcname,
		}
		out, err := l.SafelyExecuteWithDir(com, npb.fl.SrcDir(), nil)
		if err != nil {
			npb.Printf("Nanopb failed: %s\n", out)
			npb.Printf("Error: %s\n", err)
			panic("nanopb failure")
		}
		npb.Printf("File: %s [COMPILED]\n", srcname)
	}
	return nil
}
func (npb *NanoPBCompiler) Printf(format string, txt ...interface{}) {
	s := "[nanopb] "
	s = s + fmt.Sprintf(format, txt...)
	fmt.Printf("%s", s)
}
func (npb *NanoPBCompiler) Error() error {
	return npb.err
}
func (npb *NanoPBCompiler) Files(ctx context.Context, pkg *pr.Package, filetype string) ([]File, error) {
	npb.Printf("Want files for package: %v\n", pkg)
	return nil, fmt.Errorf("nanopb compiler: not impl")
}

// get a specific file
func (npb *NanoPBCompiler) GetFile(ctx context.Context, filename string) (File, error) {
	return nil, fmt.Errorf("nanopb compiler: not impl")
}
