package common

import (
	"fmt"
	pb "golang.conradwood.net/apis/protorenderer2"
	"golang.conradwood.net/protorenderer/v2/interfaces"
)

type StandardCompileResult struct {
	errors []*standard_compile_failure
}

type standard_compile_failure struct {
	pf   interfaces.ProtoFile
	comp interfaces.Compiler
	err  error
	out  []byte
}

func (scr *StandardCompileResult) AddFailed(c interfaces.Compiler, pf interfaces.ProtoFile, err error, output []byte) {
	scr.errors = append(scr.errors, &standard_compile_failure{
		out:  output,
		comp: c,
		pf:   pf,
		err:  err,
	})
}
func (scr *StandardCompileResult) GetFailures(pf interfaces.ProtoFile) []*pb.CompileFailure {
	var res []*pb.CompileFailure
	for _, cf := range scr.errors {
		if cf.pf.GetFilename() == pf.GetFilename() {
			res = append(res, cf.getCompileFailure())
		}
	}
	return res
}
func (scf *standard_compile_failure) getCompileFailure() *pb.CompileFailure {
	res := &pb.CompileFailure{
		CompilerName: "no compiler",
		ErrorMessage: fmt.Sprintf("%s", scf.err),
		Output:       scf.out,
	}
	if scf.comp != nil {
		res.CompilerName = scf.comp.ShortName()
	}
	return res
}
