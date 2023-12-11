package compiler

import (
	"context"
	"fmt"
	pr "golang.conradwood.net/apis/protorenderer"
	"golang.conradwood.net/go-easyops/linux"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/protorenderer/common"
	"golang.conradwood.net/protorenderer/filelayouter"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

/*
python seems a bit like java :) we need to hunt down various files.
specifically we need "protoc-gen-grpc_python" which apparently exists as a native library and some pip/python extension stuff.
Since we _definitely_ want to use standard protoc rather than the python parser
we'll get the lib from here and github.com/grpc/grpc and compile it from source
you'll need to use recursive:
  git clone --recursive github.com/grpc/grpc
and then
  make grpc_python_plugin

most annoyingly this is far from generic..
from "foo.bar/myproto.proto" the plugins generate this:
python (protobuf): foo/bar/myproto.bar
python (grpc)    : foo.bar/myproto.bar

to make matters worse, the 'python (grpc)' plugin way is the more proto way,
but the 'python (protobuf)' way is the more python-y way.

After a bit of digging, the python-y way introduces further challenges in the build, since we
then can no longer match files to protos (e.g. the conversion is lossy)
Thus we stick to the grpc way and move the protobuf files to the grpc directory

*/
// I reckon this might just work for php as well
type GenericCompiler struct {
	plugin_args []string // what's the name of the --xxx_out parameter?
	subdir      string   // where will the result be stored?
	fl          *filelayouter.FileLayouter
	lock        sync.Mutex
	WorkDir     string
	err         error
}

func NewNanoPBCompiler(cc CompilerCallback) Compiler {
	res := &NanoPBCompiler{
		Layouter: cc.GetFileLayouter(),
	}
	return res
}

func NewPythonCompiler(cc CompilerCallback) Compiler {
	res := &GenericCompiler{
		plugin_args: []string{"python_out", "grpc_python_out"},
		subdir:      "python",
		fl:          cc.GetFileLayouter(),
		WorkDir:     cc.GetFileLayouter().TopDir() + "build",
	}
	return res
}
func (g *GenericCompiler) Name() string { return "generic" }
func (pc *GenericCompiler) Compile(rt ResultTracker) error {
	dir := pc.fl.SrcDir()
	pc.lock.Lock()
	defer pc.lock.Unlock()
	pc.err = nil
	pc.Printf("Compiling \"%s\"\n", dir)
	files, err := AllProtos(dir)
	if err != nil {
		pc.err = err
		return err
	}
	l := linux.New()
	targetdir := pc.WorkDir + "/" + pc.subdir
	err = common.RecreateSafely(targetdir)
	if err != nil {
		pc.err = err
		return err
	}

	dirfiles, err := newDirFiles(dir, files)
	if err != nil {
		pc.err = err
		return err
	}
	/***************************** compile .proto -> .pb2.py or .php ******************************/
	cmd := []string{
		"/opt/cnw/ctools/dev/go/current/protoc/protoc",
	}
	for _, pa := range pc.plugin_args {
		cmd = append(cmd, fmt.Sprintf("--%s=%s", pa, targetdir))

	}
	if *debug {
		fmt.Printf("Compiler workdir: %s\n", dir)
	}
	pc.Printf("Compiler output dir: %s\n", targetdir)
	for _, f := range dirfiles {
		//		pc.Printf("Compiler working dir: %s, compiling %s\n", dir, f)
		cmdandfile := append(cmd, f...)
		if *debug {
			fmt.Printf("compiler command: %s\n", strings.Join(cmdandfile, " "))
		}
		out, err := l.SafelyExecuteWithDir(cmdandfile, dir, nil)
		if err != nil {
			pc.Printf("Failed to compile: %s: %s\n", f, err)
			pc.Printf("Compiler output: %s\n", out)
			pc.err = err
		}
		err = pc.MoveFile(targetdir, f)
		if err != nil {
			pc.err = err
		}
	}
	pc.Printf("Compile of %d files completed\n", len(dirfiles))
	return nil
}
func (pc *GenericCompiler) Error() error {
	return pc.err
}
func (pc *GenericCompiler) Files(ctx context.Context, pkg *pr.Package, filetype string) ([]File, error) {
	pc.Printf("Want compiled %s files for package: ID=%s, Name=%s, Prefix=%s\n", filetype, pkg.ID, pkg.Name, pkg.Prefix)
	ds := pc.WorkDir + "/" + pc.subdir + "/" + pkg.Prefix
	pc.Printf("Want compiled %s files for package: ID=%s, Name=%s, Prefix=%s. Looking in dir \"%s\"\n", filetype, pkg.ID, pkg.Name, pkg.Prefix, ds)
	fs, err := AllFiles(ds, filetype)
	if err != nil {
		return nil, err
	}
	var res []File
	for _, f := range fs {
		fn := pkg.Prefix + "/" + f
		fl := &StdFile{Filename: fn, version: 1, ctFilename: f}
		res = append(res, fl)
	}
	return res, nil
}
func (pc *GenericCompiler) GetFile(ctx context.Context, filename string) (File, error) {
	fn := pc.WorkDir + "/" + pc.subdir + "/" + filename
	fl := &StdFile{Filename: filename, version: 1, ctFilename: fn}
	return fl, nil
}
func (pc *GenericCompiler) Printf(format string, data ...interface{}) {
	s := fmt.Sprintf("[%s] ", pc.plugin_args[0])
	fmt.Printf(s+format, data...)
}

/*
e.g. from this proto:
golang.conradwood.net/apis/common/common.proto

python compiler(s) create this:
/tmp/protorenderer/29/build/python/golang.conradwood.net/apis/common/common_pb2_grpc.py
/tmp/protorenderer/29/build/python/golang/conradwood/net/apis/common/common_pb2.py

we need to move
/tmp/protorenderer/29/build/python/golang.conradwood.net/apis/common/common_pb2_grpc.py
to
/tmp/protorenderer/29/build/python/golang/conradwood/net/apis/common/common_pb2_grpc.py

*/

func (pc *GenericCompiler) MoveFile(targetdir string, filenames []string) error {
	for _, f := range filenames {
		potential_dir := targetdir + "/" + strings.Replace(filepath.Dir(f), ".", "/", -1)
		if *debug {
			fmt.Printf("moving: Checking for \"%s\"\n", potential_dir)
		}
		if !utils.FileExists(potential_dir) {
			continue
		}
		files_to_move, err := ioutil.ReadDir(potential_dir)
		if err != nil {
			return err
		}
		tdir := targetdir + "/" + filepath.Dir(f)
		os.MkdirAll(tdir, 0777) // ignore error, might exist already, bit sloppy, error comes later when moving
		for _, fm := range files_to_move {
			fname := potential_dir + "/" + fm.Name()
			tname := tdir + "/" + fm.Name()
			if *debug {
				pc.Printf("%s generated %s (moving to %s)\n", f, fname, tname)
			}
			err = os.Rename(fname, tname)
			if err != nil {
				return err
			}
		}
		err = os.Remove(potential_dir)
		if err != nil {
			return err
		}
	}
	return nil
}



















































































