package srv

import (
	"context"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/protorenderer/v2/helpers"
	"golang.conradwood.net/protorenderer/v2/interfaces"
)

type DependencyCompiler struct {
	ce  interfaces.CompilerEnvironment
	pfs []interfaces.ProtoFile
}

func NewDependencyCompiler(ce interfaces.CompilerEnvironment, pfs []interfaces.ProtoFile) *DependencyCompiler {
	res := &DependencyCompiler{ce: ce, pfs: pfs}
	return res
}

/*
this recompiles each and every dependency of the submitted files. recursively.
it throws an error if something goes wrong during compile.
A faulty file is not an error. each file is reported as a result.
a result is either SUCCESS (has been compiled) OR FAIL (failed to compile)

Actually.. this needs to be thought about a bit. And I am not certain we actually want to do this.
so for now it is a NO-OP
*/
func (dc *DependencyCompiler) Recompile(ctx context.Context) error {
	return nil
}

// recompile all the dependencies on the given file(s)...
func recompile_dependencies_with_err(ctx context.Context, ce interfaces.CompilerEnvironment, pfs []interfaces.ProtoFile, compilers []interfaces.Compiler) error {
	scr := currentVersionInfo.CompileResult()
	for _, pf := range pfs {
		err := recompile_dependencies(ctx, ce, scr, pf, compilers)
		if err != nil {
			return errors.Wrap(err)
		}
	}
	return nil
}

// recompile any files that directly or indirectly import "pf"
func recompile_dependencies(ctx context.Context, ce interfaces.CompilerEnvironment, scr interfaces.CompileResult, pf interfaces.ProtoFile, compilers []interfaces.Compiler) error {
	pfs, err := ce.MetaCache().AllWithDependencyOn(pf.GetFilename(), 0)
	if err != nil {
		return errors.Wrap(err)
	}
	var cpfs []interfaces.ProtoFile
	for _, npf := range pfs {
		spf := &helpers.StandardProtoFile{Filename: npf.ProtoFile.Filename}
		cpfs = append(cpfs, spf)
	}

	err = compile_all_compilersWithCompilerArray(ctx, ce, scr, cpfs, compilers)
	if err != nil {
		return errors.Wrap(err)
	}
	for _, cpf := range cpfs {
		if helpers.ContainsFailure(scr.GetResults(cpf)) {
			return errors.Errorf("failed to compile dependency \"%s\"\n", pf.GetFilename())
		}
	}
	return nil
}
