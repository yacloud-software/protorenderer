package main

import (
	"context"
	"fmt"
	pr "golang.conradwood.net/apis/protorenderer"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/protorenderer/compiler"
	"golang.conradwood.net/protorenderer/meta"
	"path/filepath"
)

// this is different to python/java/go because it's not through a compiler
func (p *protoRenderer) GetFilesProto(ctx context.Context, req *pr.ID) (*pr.FilenameList, error) {
	result := completeVersion.metaCompiler.GetMostRecentResult()
	if result == nil {
		return nil, errors.Unavailable(ctx, "GetPackages")
	}
	id := req.ID
	if id == "" {
		return nil, errors.InvalidArgs(ctx, "missing packageid", "missing packageid")
	}
	fl := &pr.FilenameList{}
	for _, p := range result.Packages {
		if p.Proto.ID != req.ID {
			continue
		}
		for _, pf := range p.Protofiles {
			fl.Files = append(fl.Files, pf.Filename)
		}
	}
	return fl, nil
}
func (p *protoRenderer) GetFilesPython(ctx context.Context, req *pr.ID) (*pr.FilenameList, error) {
	filetype := ".py"
	compiler := completeVersion.pythonCompiler

	result := completeVersion.metaCompiler.GetMostRecentResult()
	if result == nil {
		return nil, errors.Unavailable(ctx, "GetPackages")
	}
	var pkg *meta.Package
	id := req.ID
	if id == "" {
		return nil, errors.InvalidArgs(ctx, "missing packageid", "missing packageid")
	}
	for _, pm := range result.Packages {
		if pm.Proto.ID == id {
			pkg = pm
			break
		}
	}
	if pkg == nil {
		return nil, errors.InvalidArgs(ctx, "no such package", "no package \"%s\"", id)
	}
	//	fmt.Printf("Need java class files for %s\n", pkg.FQDN)
	rf, err := compiler.Files(ctx, pkg.Proto, filetype)
	if err != nil {
		return nil, err
	}
	fl := &pr.FilenameList{}
	for _, f := range rf {
		s := f.GetFilename()
		fl.Files = append(fl.Files, s)
	}
	return fl, nil
}
func (p *protoRenderer) GetFilesJavaClass(ctx context.Context, req *pr.ID) (*pr.FilenameList, error) {
	filetype := ".class"
	compiler := completeVersion.javaCompiler

	result := completeVersion.metaCompiler.GetMostRecentResult()
	if result == nil {
		return nil, errors.Unavailable(ctx, "GetPackages")
	}
	var pkg *meta.Package
	id := req.ID
	if id == "" {
		return nil, errors.InvalidArgs(ctx, "missing packageid", "missing packageid")
	}
	for _, pm := range result.Packages {
		if pm.Proto.ID == id {
			pkg = pm
			break
		}
	}
	if pkg == nil {
		return nil, errors.InvalidArgs(ctx, "no such package", "no package \"%s\"", id)
	}
	//	fmt.Printf("Need java class files for %s\n", pkg.FQDN)
	rf, err := compiler.Files(ctx, pkg.Proto, filetype)
	if err != nil {
		return nil, err
	}
	fl := &pr.FilenameList{}
	for _, f := range rf {
		s := f.GetFilename()
		fl.Files = append(fl.Files, s)
	}
	return fl, nil

}
func (p *protoRenderer) GetFilesGO(ctx context.Context, req *pr.ID) (*pr.FilenameList, error) {
	filetype := ".go"
	compiler := completeVersion.goCompiler

	result := completeVersion.metaCompiler.GetMostRecentResult()
	if result == nil {
		return nil, errors.Unavailable(ctx, "GetPackages")
	}
	var pkg *meta.Package
	id := req.ID
	if id == "" {
		return nil, errors.InvalidArgs(ctx, "missing packageid", "missing packageid")
	}
	for _, pm := range result.Packages {
		if pm.Proto.ID == id {
			pkg = pm
			break
		}
	}
	if pkg == nil {
		return nil, errors.InvalidArgs(ctx, "no such package", "no package \"%s\"", id)
	}
	//	fmt.Printf("Need java class files for %s\n", pkg.FQDN)
	rf, err := compiler.Files(ctx, pkg.Proto, filetype)
	if err != nil {
		return nil, err
	}
	fl := &pr.FilenameList{}
	for _, f := range rf {
		s := f.GetFilename()
		fl.Files = append(fl.Files, s)
	}
	return fl, nil

}
func (p *protoRenderer) GetFilesNanoPB(ctx context.Context, req *pr.ID) (*pr.FilenameList, error) {
	filetype := "."
	compiler := completeVersion.nanopbCompiler

	result := completeVersion.metaCompiler.GetMostRecentResult()
	if result == nil {
		return nil, errors.Unavailable(ctx, "GetPackages")
	}
	var pkg *meta.Package
	id := req.ID
	if id == "" {
		return nil, errors.InvalidArgs(ctx, "missing packageid", "missing packageid")
	}
	for _, pm := range result.Packages {
		if pm.Proto.ID == id {
			pkg = pm
			break
		}
	}
	if pkg == nil {
		return nil, errors.InvalidArgs(ctx, "no such package", "no package \"%s\"", id)
	}

	rf, err := compiler.Files(ctx, pkg.Proto, filetype)
	if err != nil {
		return nil, err
	}
	fl := &pr.FilenameList{}
	for _, f := range rf {
		s := f.GetFilename()
		fl.Files = append(fl.Files, s)
	}
	return fl, nil
}

// get a specific file
func (p *protoRenderer) GetFile(ctx context.Context, req *pr.FileRequest) (*pr.File, error) {
	cv := completeVersion
	if cv == nil {
		return nil, errors.Unavailable(ctx, "getfile")
	}
	var compiler compiler.Compiler
	suffix := filepath.Ext(req.Filename)
	if suffix == ".proto" {
		// special case, get it from cache
		pf, err := protocache.GetFile(ctx, req.Filename)
		if err != nil {
			return nil, err
		}
		if pf == nil {
			return nil, errors.NotFound(ctx, "file not found", "file %s not found", req.Filename)
		}
		return &pr.File{
			Content:      []byte(pf.Content),
			RepositoryID: pf.RepositoryID,
		}, nil
	}

	// do we need to guess the compiler?
	if req.Compiler == pr.CompilerType_UNDEFINED {
		if suffix == ".class" || suffix == ".java" {
			compiler = cv.javaCompiler
		} else if suffix == ".py" {
			compiler = cv.pythonCompiler
		} else if suffix == ".go" || suffix == ".pb.go" {
			compiler = cv.goCompiler
		} else {
			return nil, fmt.Errorf("unknown file \"%s\"", suffix)
		}
	} else {
		if req.Compiler == pr.CompilerType_GOLANG {
			compiler = cv.goCompiler
		} else if req.Compiler == pr.CompilerType_JAVA {
			compiler = cv.javaCompiler
		} else if req.Compiler == pr.CompilerType_PYTHON {
			compiler = cv.pythonCompiler
		} else if req.Compiler == pr.CompilerType_NANOPB {
			compiler = cv.nanopbCompiler
		} else {
			return nil, errors.InvalidArgs(ctx, "invalid compilertype", "invalid compilertype (%v)", req.Compiler)
		}
	}
	f, err := compiler.GetFile(ctx, req.Filename)
	if err != nil {
		return nil, err
	}
	ct, err := f.GetContent()
	if err != nil {
		return nil, err
	}
	res := &pr.File{Content: ct}
	return res, nil
}













