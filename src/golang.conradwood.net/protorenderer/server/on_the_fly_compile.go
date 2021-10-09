package main

import (
	"context"
	"fmt"
	pb "golang.conradwood.net/apis/protorenderer"
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

func (e *protoRenderer) CompileFile(ctx context.Context, req *pb.CompileRequest) (*pb.CompileResult, error) {
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

	targetdir := dir + "/build/go"
	os.MkdirAll(targetdir, 0777)
	incdir := current.filelayouter.SrcDir()
	fmt.Printf("Incdir    : \"%s\"\n", incdir)

	dir = dir + "/src"
	res := &pb.CompileResult{SourceFilename: fname}

	// compile go, plugin protoc-gen-go
	pcfname := compiler.FindCompiler("protoc-gen-go")
	cmd := []string{
		"/opt/cnw/ctools/dev/go/current/protoc/protoc",
		fmt.Sprintf("-I%s", incdir),
		fmt.Sprintf("-I%s", dir),
		fmt.Sprintf("--plugin=protoc-gen-go=%s", pcfname),
		fmt.Sprintf("--go_out=plugins=grpc:%s", targetdir),
	}

	err = e.CompileGoWithPlugin(cmd, dir, fname, targetdir, res)
	if err != nil {
		return nil, err
	}

	// compile go, plugin protoc-gen-cnw
	pcfname = compiler.FindCompiler("protoc-gen-cnw")
	cmd = []string{
		"/opt/cnw/ctools/dev/go/current/protoc/protoc",
		fmt.Sprintf("-I%s", incdir),
		fmt.Sprintf("-I%s", dir),
		fmt.Sprintf("--plugin=protoc-gen-cnw=%s", pcfname),
		fmt.Sprintf("--cnw_out=%s", targetdir),
	}

	err = e.CompileGoWithPlugin(cmd, dir, fname, targetdir, res)
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
