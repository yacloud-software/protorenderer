package compiler

import (
	"context"
	"flag"
	pr "golang.conradwood.net/apis/protorenderer"
)

var (
	debug = flag.Bool("debug_compiler", false, "debug compilers")
)

type Compiler interface {
	Compile() error
	// returns most recent error (or nil)
	Error() error
	// filetype is specific for the compiler. e.g. .class for java or .pb.go for go
	Files(ctx context.Context, pkg *pr.Package, filetype string) ([]File, error)
	// get a specific file
	GetFile(ctx context.Context, filename string) (File, error)
}

type File interface {
	GetVersion() int
	GetFilename() string
	GetContent() ([]byte, error)
}
