package versioninfo

import (
	pb "golang.conradwood.net/apis/protorenderer2"
	"golang.conradwood.net/protorenderer/v2/interfaces"
)

type compileresult struct {
}

func (cr *compileresult) AddFailed(c interfaces.Compiler, pf interfaces.ProtoFile, err error, output []byte) {
}
func (cr *compileresult) GetFailures(pf interfaces.ProtoFile) []*pb.CompileFailure {
	return nil
}
