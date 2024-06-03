package main

import (
	"golang.conradwood.net/protorenderer/v2/interfaces"
)

type StandardCompileResult struct {
}

func (scr *StandardCompileResult) AddFailed(pf interfaces.ProtoFile, err error) {
}
func (scr *StandardCompileResult) HasFailed(pf interfaces.ProtoFile) bool {
	return false
}
