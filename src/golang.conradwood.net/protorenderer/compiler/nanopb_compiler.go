package compiler

import (
	"context"
	"flag"
	"fmt"
	pr "golang.conradwood.net/apis/protorenderer"
	"golang.conradwood.net/go-easyops/linux"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/protorenderer/common"
	"golang.conradwood.net/protorenderer/filelayouter"
	"strings"
	"sync"
)

var (
	nanopb_flags = flag.String("nanopb_options", "", "comma delimited key=value options that will be passed to nanopb_generator with -s option")
)

type NanoPBCompiler struct {
	lock      sync.Mutex
	err       error
	fl        *filelayouter.FileLayouter
	targetdir string
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
	npb.targetdir = npb.WorkDir() + "/nanopb"
	err = common.RecreateSafely(npb.targetdir)
	if err != nil {
		npb.err = err
		return err
	}
	npb.Printf("Sourcedir: %s\n", npb.fl.SrcDir())
	npb.Printf("Workdir:   %s\n", npb.WorkDir())
	npb.Printf("Targetdir: %s\n", npb.targetdir)
	for _, f := range files {
		srcname := f
		npb.Printf("File: %s [COMPILING]\n", srcname)
		com := []string{
			"nanopb_generator.py",
			"-D", npb.targetdir, // output dir
			"-Q", `#include "nanopb/%s"`,
			"-L", `#include <nanopb/%s>`,
		}
		com = addNanoPBOptions(com)
		com = append(com, srcname)
		out, err := l.SafelyExecuteWithDir(com, npb.fl.SrcDir(), nil)
		if err != nil {
			npb.Printf("Nanopb failed: %s\n", out)
			npb.Printf("File %s [Error: %s]\n", srcname, err)
			// ignore errors for now..
			continue
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
	ds := npb.WorkDir() + "/go/" + pkg.Prefix
	fmt.Printf("Targetdir: %s\n", npb.targetdir)
	// examples of paths:
	// targetdir: /tmp/wd//1438/build/nanopb
	// pkg.Prefix: golang.conradwood.net/apis/weblogin
	fpath := npb.targetdir + "/" + pkg.Prefix
	if !utils.FileExists(fpath) {
		// have 0 files matching this
		return []File{}, nil
	}
	fnames, err := AllFiles(fpath, "")
	if err != nil {
		return nil, err
	}
	var res []File
	for _, f := range fnames {
		npb.Printf("File: \"%s\"\n", f)
		fn := pkg.Prefix + "/" + f
		fl := &StdFile{Filename: fn, version: 1, ctFilename: ds + "/" + f}
		res = append(res, fl)
	}
	return res, nil
}

// get a specific file
func (npb *NanoPBCompiler) GetFile(ctx context.Context, filename string) (File, error) {
	fn := npb.targetdir + "/" + filename
	npb.Printf("Want file \"%s\" => \"%s\"\n", filename, fn)
	fl := &StdFile{Filename: filename, version: 1, ctFilename: fn}
	return fl, nil
}
func addNanoPBOptions(com []string) []string {
	if *nanopb_flags == "" {
		return com
	}
	res := com
	cl := strings.Split(*nanopb_flags, ",")
	for _, opt := range cl {
		kvs := strings.Split(opt, "=")
		if len(kvs) != 2 {
			fmt.Printf("[nanopb] opt \"%s\" does not split into 2 parts, but %d\n", opt, len(kvs))
			return com
		}
		k := kvs[0]
		v := kvs[1]
		k = strings.Trim(k, " ")
		v = strings.Trim(v, " ")
		res = append(res, "-s")
		res = append(res, fmt.Sprintf("%s:%s", k, v))
	}
	return res
}
