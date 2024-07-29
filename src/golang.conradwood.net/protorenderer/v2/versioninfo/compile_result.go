package versioninfo

import (
	"fmt"
	pb "golang.conradwood.net/apis/protorenderer2"
	"golang.conradwood.net/protorenderer/v2/interfaces"
	"sync"
)

type compileresult struct {
	sync.Mutex
	vi      *VersionInfo
	results []*compileresult_file
}
type compileresult_file struct {
	comp   interfaces.Compiler
	pf     interfaces.ProtoFile
	fail   bool
	err    error
	output []byte
}

func (cr *compileresult) AddSuccess(c interfaces.Compiler, pf interfaces.ProtoFile) {
	cr.Lock()
	defer cr.Unlock()
	cr.results = append(cr.results, &compileresult_file{comp: c, pf: pf})
}
func (cr *compileresult) AddFailed(c interfaces.Compiler, pf interfaces.ProtoFile, err error, output []byte) {
	cr.Lock()
	defer cr.Unlock()
	cr.results = append(cr.results, &compileresult_file{
		comp:   c,
		pf:     pf,
		fail:   true,
		err:    err,
		output: output,
	})
}
func (cr *compileresult) GetResults(pf interfaces.ProtoFile) []*pb.CompileResult {
	cr.Lock()
	defer cr.Unlock()
	var res []*pb.CompileResult
	for _, crf := range cr.results {
		if crf.pf.GetFilename() == pf.GetFilename() {
			cf := crf.toCompileResultProto()
			res = append(res, cf)
		}
	}
	return res
}
func (cr *compileresult) getresultsforfile(filename string) []*compileresult_file {
	cr.Lock()
	defer cr.Unlock()
	var res []*compileresult_file
	for _, crf := range cr.results {
		if crf.pf.GetFilename() == filename {
			res = append(res, crf)
		}
	}
	return res
}
func (cr *compileresult_file) toCompileResultProto() *pb.CompileResult {
	cf := &pb.CompileResult{
		CompilerName: cr.comp.ShortName(),
		Output:       string(cr.output),
		Success:      !cr.fail,
	}
	if cr.err != nil {
		cf.ErrorMessage = fmt.Sprintf("%s", cr.err)
	}
	return cf
}
