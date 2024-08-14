package srv

import (
	"golang.conradwood.net/protorenderer/v2/interfaces"
)

type CompareResult struct {
	c1 interfaces.CompileResult
	c2 interfaces.CompileResult
}

func NewCompareResult(c1, c2 interfaces.CompileResult) *CompareResult {
	res := &CompareResult{c1: c1, c2: c2}
	return res
}

// TODO: work out the rules for this
func (cr *CompareResult) IsWorse() bool {
	return false
}
