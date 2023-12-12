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
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	nanopb_version = ""
	nanopb_gen     = flag.String("nanopb_binary", "", "executable name for generator")
	nanopb_flags   = flag.String("nanopb_options", "", "comma delimited key=value options that will be passed to nanopb_generator with -s option")
)

type NanoPBCompiler struct {
	lock      sync.Mutex
	err       error
	Layouter  *filelayouter.FileLayouter
	TargetDir string
}

func (npb *NanoPBCompiler) WorkDir() string {
	return npb.Layouter.TopDir() + "build"
}
func (g *NanoPBCompiler) Name() string { return "nanopb" }
func (npb *NanoPBCompiler) Compile(rt ResultTracker) error {
	npb.Printf("***************** nanopb compiler ******************\n")
	dir := npb.Layouter.SrcDir()
	npb.lock.Lock()
	defer npb.lock.Unlock()
	npb.err = nil
	fmt.Printf("Compiling nanopb \"%s\"\n", dir)
	files, err := AllProtos(dir)
	if err != nil {
		npb.err = err
		return err
	}
	//	l.MyIP()
	npb.TargetDir = npb.WorkDir() + "/nanopb"
	err = common.RecreateSafely(npb.TargetDir)
	if err != nil {
		npb.err = err
		return err
	}
	npb.Printf("Sourcedir: %s\n", npb.Layouter.SrcDir())
	npb.Printf("Workdir:   %s\n", npb.WorkDir())
	npb.Printf("Targetdir: %s\n", npb.TargetDir)
	for _, f := range files {
		srcname := f
		npb.Printf("File: %s [COMPILING]\n", srcname)
		com := []string{
			Nanopb_binary(),
			"-D", npb.TargetDir, // output dir
			"-Q", `#include "nanopb/%s"`,
			"-L", `#include <nanopb/%s>`,
			"--strip-path",
		}
		com = AddNanoPBOptions(com)
		com = append(com, srcname)
		l := linux.New()
		out, err := l.SafelyExecuteWithDir(com, npb.Layouter.SrcDir(), nil)
		if err != nil {
			npb.Printf("Nanopb failed: %s\n", out)
			npb.Printf("File %s [Error: %s]\n", srcname, err)
			// ignore errors for now..
			continue
		}
		npb.Printf("File: %s [COMPILED]\n", srcname)
		err = addCustomFiles(srcname, npb.TargetDir)
		if err != nil {
			npb.Printf("Custom files failed: %s\n", err)
			continue
		}
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
	ds := npb.WorkDir() + "/nanopb/" + pkg.Prefix
	fmt.Printf("Targetdir: %s\n", npb.TargetDir)
	// examples of paths:
	// targetdir: /tmp/wd//1438/build/nanopb
	// pkg.Prefix: golang.conradwood.net/apis/weblogin
	fpath := npb.TargetDir + "/" + pkg.Prefix
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
	fn := npb.TargetDir + "/" + filename
	npb.Printf("Want file \"%s\" => \"%s\"\n", filename, fn)
	fl := &StdFile{Filename: filename, version: 1, ctFilename: fn}
	return fl, nil
}
func AddNanoPBOptions(com []string) []string {
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

// write some custom files
func addCustomFiles(srcfile string, targetdir string) error {
	if !strings.HasSuffix(srcfile, ".proto") {
		return nil
	}
	outfile := strings.TrimSuffix(srcfile, ".proto")
	outfile = targetdir + "/" + outfile + "_def.h"
	templ := `/*
Compile-Options: %s
Compile-Timestamp: %d
Compile-Time: %s
Compiler-Version: %s
/*`
	t := time.Now()
	s := fmt.Sprintf(templ, *nanopb_flags, t.Unix(), utils.TimeString(t), get_nanopb_version())
	err := utils.WriteFile(outfile, []byte(s))
	return err
}

func get_nanopb_version() string {
	if nanopb_version != "" {
		return nanopb_version
	}
	l := linux.New()
	out, err := l.SafelyExecuteWithDir([]string{Nanopb_binary(), "--version"}, "/", nil)
	if err != nil {
		fmt.Printf("Failure: %s\n", out)
		fmt.Printf("Failure: %s\n", err)
		panic("failed to execute nanopb_generator.py")
	}
	out = strings.Trim(out, "\n")
	out = strings.Trim(out, "\r")
	out = strings.Trim(out, " ")
	nanopb_version = out
	return out
}

func Nanopb_binary() string {
	res := find_nanopb_binary()
	if !utils.FileExists(res) {
		panic("nanopb_generator.py (" + res + ") is not executable")
	}
	return res
}
func find_nanopb_binary() string {
	if *nanopb_gen != "" {
		return *nanopb_gen
	}
	paths := []string{"/sbin", "/usr/sbin", "/bin", "/usr/bin", "/usr/local/bin"}
	for _, p := range paths {
		filename := p + "/nanopb_generator.py"
		if utils.FileExists(filename) {
			return filename
		}
	}
	f, err := utils.FindFile(fmt.Sprintf("extra/compilers/%s/nanopb/nanopb_generator.py", common.GetCompilerVersion()))
	if err == nil {
		f, err = filepath.Abs(f)
		if err == nil {
			return f
		}
	}
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("failed to get current working directory: %s\n", err)
	}
	return pwd + "/" + fmt.Sprintf("extra/compilers/%s/nanopb/nanopb_generator.py", common.GetCompilerVersion())

}


























































































