package nanopb

import (
	"context"
	"fmt"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/linux"
	"golang.conradwood.net/go-easyops/utils"
	//	"golang.conradwood.net/protorenderer/v2/helpers"
	"golang.conradwood.net/protorenderer/cmdline"
	"golang.conradwood.net/protorenderer/v2/interfaces"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	nanopb_version = ""
)

type NanoPBCompiler struct {
}

func New() interfaces.Compiler {
	res := &NanoPBCompiler{}
	return res
}
func (g *NanoPBCompiler) ShortName() string { return "nanopb" }
func (npb *NanoPBCompiler) Compile(ctx context.Context, ce interfaces.CompilerEnvironment, files []interfaces.ProtoFile, outdir string, cr interfaces.CompileResult) error {
	npb.Printf("***************** nanopb compiler ******************\n")
	dir := ce.NewProtosDir()
	srcdir := ce.NewProtosDir()
	fmt.Printf("Compiling nanopb \"%s\"\n", dir)
	targetDir := ce.CompilerOutDir() + "/nanopb"
	err := utils.RecreateSafely(targetDir)
	if err != nil {
		return err
	}

	for _, f := range files {
		srcname := f.GetFilename()
		npb.Printf("File: %s [COMPILING]\n", srcname)
		com := []string{
			Nanopb_binary(),
			"-D", targetDir, // output dir
			"-Q", `#include "nanopb/%s"`,
			"-L", `#include <nanopb/%s>`,
			"--strip-path",
		}
		com = AddNanoPBOptions(com)
		com = append(com, srcname)
		l := linux.New()
		out, err := l.SafelyExecuteWithDir(com, srcdir, nil)
		if err != nil {
			npb.Printf("Nanopb failed: %s\n", out)
			npb.Printf("File %s [Error: %s]\n", srcname, err)
			// ignore errors for now..
			continue
		}
		npb.Printf("File: %s [COMPILED]\n", srcname)
		err = addCustomFiles(srcname, targetDir)
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

func AddNanoPBOptions(com []string) []string {
	if cmdline.GetNanoPBFlags() == "" {
		return com
	}
	res := com
	cl := strings.Split(cmdline.GetNanoPBFlags(), ",")
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
	s := fmt.Sprintf(templ, cmdline.GetNanoPBFlags(), t.Unix(), utils.TimeString(t), get_nanopb_version())
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
	paths := []string{"/sbin", "/usr/sbin", "/bin", "/usr/bin", "/usr/local/bin"}
	for _, p := range paths {
		filename := p + "/nanopb_generator.py"
		if utils.FileExists(filename) {
			return filename
		}
	}
	f, err := utils.FindFile(fmt.Sprintf("extra/compilers/%s/nanopb/nanopb_generator.py", cmdline.GetCompilerVersion()))
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
	return pwd + "/" + fmt.Sprintf("extra/compilers/%s/nanopb/nanopb_generator.py", cmdline.GetCompilerVersion())

}
func (c *NanoPBCompiler) DirsForPackage(ctx context.Context, package_name string) ([]string, error) {
	return nil, errors.NotImplemented(ctx, c.ShortName()+".DirsForPackage")
}
