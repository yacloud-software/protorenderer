package interfaces

import (
	"context"
)

// a .proto file
type ProtoFile interface {
	GetID() uint64
	GetFilename() string // always relative
	Content() []byte
}

type CompileResult interface {
	GetProtoFile() ProtoFile // result for a given protofile
	GetError() error         // might be nil, if success
}

// e.g. implemented by go, java, python, php etc..
// not all compilers can compile everything unfortunately.
// the MetaCompiler is a different type of compiler though
type Compiler interface {
	ShortName() string // e.g. "java" or "golang" or "nanopb" compilers must put all of their results in a directory with this name
	/*
		The compiler might fail a specific file, e.g. with a Syntax Error, but it might also fail completely, such as OOM or Disk full
		thus it can return an error as well as Compile Resuls
	*/
	Compile(ctx context.Context, ce CompilerEnvironment, files []ProtoFile) ([]CompileResult, error)
}
