package compiler

import (
	"context"
	"flag"
	"fmt"
	pr "golang.conradwood.net/apis/protorenderer"
)

var (
	debug = flag.Bool("debug_compiler", false, "debug compilers")
)

type ResultTracker interface {
	AddFailed(c Compiler, filename string, errormessage string)
}
type Compiler interface {
	Compile(ResultTracker) error
	// returns most recent error (or nil)
	Error() error
	// filetype is specific for the compiler. e.g. .class for java or .pb.go for go
	Files(ctx context.Context, pkg *pr.Package, filetype string) ([]File, error)
	// get a specific file
	GetFile(ctx context.Context, filename string) (File, error)
	// compiler name
	Name() string
}

type File interface {
	GetVersion() int
	GetFilename() string
	GetContent() ([]byte, error)
}

func Debugf(format string, args ...interface{}) {
	if !*debug {
		return
	}
	fmt.Printf("[compiler ] "+format, args...)
}








































































































