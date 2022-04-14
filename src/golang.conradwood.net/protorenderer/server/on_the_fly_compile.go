package main

import (
	"context"
	"fmt"
	pb "golang.conradwood.net/apis/protorenderer"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/linux"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/protorenderer/common"
	"golang.conradwood.net/protorenderer/compiler"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	ontheflylock sync.Mutex
	ontheflynum  uint64
)

func newOnTheFlyNum() uint64 {
	ontheflylock.Lock()
	defer ontheflylock.Unlock()
	ontheflynum++
	return ontheflynum
}

func normaliseFilename(filename string) string {
	idx := strings.Index(filename, "protos/")
	if idx == -1 {
		return filename
	}
	res := filename[idx:]
	return res
}

type onthefly_compiler struct {
	protofile string             // the full filename of the proto file to compile
	workdir   string             // the directory to place the results into
	req       *pb.CompileRequest // the compile request
	renderer  *protoRenderer
}

func (oc *onthefly_compiler) WorkDir() string {
	return oc.workdir
}

func (e *protoRenderer) CompileFile(ctx context.Context, req *pb.CompileRequest) (*pb.CompileResult, error) {
	// default: golang compiler
	if len(req.Compilers) == 0 {
		req.Compilers = []pb.CompilerType{
			pb.CompilerType_GOLANG,
		}
	}
	dir := TopDir() + fmt.Sprintf("/onthefly/%d", newOnTheFlyNum())
	fmt.Printf("Compilefile in dir \"%s\"\n", dir)
	err := common.RecreateSafely(dir)
	if err != nil {
		return nil, err
	}
	fname := normaliseFilename(req.AddProtoRequest.Name)
	fullfilename := dir + "/src/" + fname
	os.MkdirAll(filepath.Dir(fullfilename), 0777)
	err = utils.WriteFile(fullfilename, []byte(req.AddProtoRequest.Content))
	if err != nil {
		return nil, err
	}
	oc := &onthefly_compiler{
		protofile: fullfilename,
		workdir:   dir,
		renderer:  e,
		req:       req,
	}
	res := &pb.CompileResult{
		SourceFilename: fullfilename,
	}
	for _, compiler := range req.Compilers {
		// add new compilers here
		if compiler == pb.CompilerType_GOLANG {
			cr, err := oc.golang(ctx)
			if err != nil || cr.CompileError != "" {
				return cr, err
			}
			res.Files = append(res.Files, cr.Files...)
		} else if compiler == pb.CompilerType_NANOPB {
			cr, err := oc.nanopb(ctx)
			if err != nil || cr.CompileError != "" {
				return cr, err
			}
			res.Files = append(res.Files, cr.Files...)
		} else {
			return nil, errors.NotImplemented(ctx, fmt.Sprintf("compiler %v not implemented to compile on-the-fly", compiler))
		}
	}
	return res, nil
}

// compile go file
func (oc *onthefly_compiler) golang(ctx context.Context) (*pb.CompileResult, error) {
	dir := oc.workdir
	fname := oc.protofile
	targetdir := dir + "/build/go"
	os.MkdirAll(targetdir, 0777)
	incdir := current.filelayouter.SrcDir()
	fmt.Printf("Incdir    : \"%s\"\n", incdir)

	dir = dir + "/src"
	res := &pb.CompileResult{SourceFilename: fname}

	// compile go, plugin protoc-gen-go
	pcfname := compiler.FindCompiler("protoc-gen-go")
	cmd := []string{
		"/opt/yacloud/ctools/dev/go/current/protoc/protoc",
		fmt.Sprintf("-I%s", incdir),
		fmt.Sprintf("-I%s", dir),
		fmt.Sprintf("--plugin=protoc-gen-go=%s", pcfname),
		fmt.Sprintf("--go_out=plugins=grpc:%s", targetdir),
	}

	err := oc.renderer.CompileGoWithPlugin(cmd, dir, fname, targetdir, res)
	if err != nil {
		return nil, err
	}

	// compile go, plugin protoc-gen-cnw
	pcfname = compiler.FindCompiler("protoc-gen-cnw")
	cmd = []string{
		"/opt/yacloud/ctools/dev/go/current/protoc/protoc",
		fmt.Sprintf("-I%s", incdir),
		fmt.Sprintf("-I%s", dir),
		fmt.Sprintf("--plugin=protoc-gen-cnw=%s", pcfname),
		fmt.Sprintf("--cnw_out=%s", targetdir),
	}

	err = oc.renderer.CompileGoWithPlugin(cmd, dir, fname, targetdir, res)
	if err != nil {
		return nil, err
	}
	return res, nil

}

func (e *protoRenderer) CompileGoWithPlugin(cmd []string, curdir string, fname string, targetdir string, res *pb.CompileResult) error {
	l := linux.New()
	cmdandfile := append(cmd, fname)

	fmt.Printf("Targetdir : \"%s\"\n", targetdir)
	fmt.Printf("Currentdir: \"%s\"\n", curdir)
	out, err := l.SafelyExecuteWithDir(cmdandfile, curdir, nil)
	if err != nil {
		s1 := fmt.Sprintf("Failed to compile: %s: %s", fname, err)
		s2 := fmt.Sprintf("Compiler output: %s", out)
		fmt.Println(s1)
		fmt.Println(s2)
		res.CompileError = fmt.Sprintf("protoc failed to compile: %s\nCompiler Output:\n%s", fname, out)
		return nil
	}
	if *debug {
		fmt.Printf("Compiler output: %s\n", out)
	}

	// pick up the files and return them
	results, err := compiler.AllFiles(targetdir, ".go")
	if err != nil {
		return err
	}

	for _, f := range results {
		ct, err := utils.ReadFile(targetdir + "/" + f)
		if err != nil {
			res.CompileError = fmt.Sprintf("Binary file not loadable (%s)", err)
			return err
		}
		if *debug {
			fmt.Printf("Created file %s\n", f)
		}
		cf := &pb.CompiledFile{
			Filename: f,
			Content:  ct,
		}
		res.Files = append(res.Files, cf)
	}
	return nil
}

// compile go file
func (oc *onthefly_compiler) nanopb(ctx context.Context) (*pb.CompileResult, error) {
	l := linux.New()
	targetDir := oc.WorkDir() + "/nanopb"
	err := common.RecreateSafely(targetDir)
	if err != nil {
		return nil, err
	}
	layouter := current.filelayouter
	fmt.Printf("Sourcedir: %s\n", layouter.SrcDir())
	fmt.Printf("Workdir:   %s\n", oc.WorkDir())
	fmt.Printf("Targetdir: %s\n", targetDir)
	srcname := oc.protofile
	fmt.Printf("File: %s [COMPILING]\n", srcname)
	com := []string{
		compiler.Nanopb_binary(),
		"-D", targetDir, // output dir
		"-Q", `#include "nanopb/%s"`,
		"-L", `#include <nanopb/%s>`,
		"-I", oc.WorkDir() + "/src/protos",
		"--strip-path",
	}
	com = compiler.AddNanoPBOptions(com)
	com = append(com, srcname)
	out, err := l.SafelyExecuteWithDir(com, layouter.SrcDir(), nil)
	if err != nil {
		fmt.Printf("Nanopb failed: %s\n", out)
		fmt.Printf("File %s [Error: %s]\n", srcname, err)
		// ignore errors for now..
		return nil, err
	}
	fmt.Printf("File: %s [COMPILED]\n", srcname)
	/*
		err = addCustomFiles(srcname, npb.TargetDir)
		if err != nil {
			fmt.Printf("Custom files failed: %s\n", err)
			continue
		}
	*/

	// pick up the files and return them
	results, err := compiler.AllFiles(targetDir, "")
	if err != nil {
		return nil, err
	}

	res := &pb.CompileResult{}

	for _, f := range results {
		ct, err := utils.ReadFile(targetDir + "/" + f)
		if err != nil {
			res.CompileError = fmt.Sprintf("Binary file not loadable (%s)", err)
			return nil, err
		}
		if *debug {
			fmt.Printf("Created file %s\n", f)
		}
		cf := &pb.CompiledFile{
			Filename: f,
			Content:  ct,
		}
		res.Files = append(res.Files, cf)
	}
	return res, nil

}
