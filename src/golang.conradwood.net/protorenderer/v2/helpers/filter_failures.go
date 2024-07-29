package helpers

import (
	pb "golang.conradwood.net/apis/protorenderer2"
)

// if at least one compileresult is a failure, return true
func ContainsFailure(in []*pb.CompileResult) bool {
	for _, result := range in {
		if !result.Success {
			return true
		}
	}
	return false
}
