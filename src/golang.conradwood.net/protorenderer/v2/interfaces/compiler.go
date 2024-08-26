package interfaces

import (
	"context"
	pb "golang.conradwood.net/apis/protorenderer2"
)

// this is a .proto file. usually instantiated via the database (e.g. a DBProtoFile implements this)
type ProtoFile interface {
	GetID() uint64
	GetFilename() string // always relative
}

// e.g. implemented by go, java, python, php etc..
// not all compilers can compile everything unfortunately.
// the MetaCompiler is a different type of compiler though
type Compiler interface {
	ShortName() string // e.g. "java" or "golang" or "nanopb" compilers must put all of their results in a directory with this name
	/*
		The compiler might fail a specific file, e.g. with a Syntax Error, but it might also fail completely, such as OOM or Disk full
		thus it can return an error as well as Compile Result
		All dirs are relative to CompilerEnvironment.WorkDir()
		The compiler is expected to:
		1. use CompilerEnvironment.AllKnownProtosDir() as an include search path
		2. compile any .proto files as specified by ProtoFile
		3. and put the output into outdir.
	*/
	Compile(ctx context.Context, ce CompilerEnvironment, files []ProtoFile, outdir string, cr CompileResult) error
	DirsForPackage(ctx context.Context, package_name string) ([]string, error) // given a package, e.g. golang.conradwood.net/apis/registry will return the dirs with files created by this, e.g. []string{ "java/classes/net/conradwood/golang/apis/h2gproxy",}

}

type CompileResult interface {
	AddSuccess(c Compiler, pf ProtoFile)
	AddFailed(c Compiler, pf ProtoFile, err error, output []byte)
	GetResults(ProtoFile) []*pb.CompileResult // get results for all compilers
	GetFailed() []*pb.CompileResult
}
