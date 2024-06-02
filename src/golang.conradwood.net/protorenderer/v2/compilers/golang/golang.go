package golang

import (
	"context"
	"golang.conradwood.net/protorenderer/v2/interfaces"
)

type golangCompiler struct{}

func New() interfaces.Compiler {
	return &golangCompiler{}
}
func (gc *golangCompiler) ShortName() string { return "golang" }
func (gc *golangCompiler) Compile(ctx context.Context, ce interfaces.CompilerEnvironment, files []interfaces.ProtoFile) ([]interfaces.CompileResult, error) {
	return nil, nil
}
