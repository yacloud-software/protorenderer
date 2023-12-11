package main

import (
	fl "golang.conradwood.net/protorenderer/filelayouter"
	"golang.conradwood.net/protorenderer/meta"
)

type compilerCallback struct {
	nfly         *fl.FileLayouter
	metacompiler *meta.MetaCompiler
}

func (cc *compilerCallback) GetFileLayouter() *fl.FileLayouter {
	return cc.nfly
}
func (cc *compilerCallback) GetMetaPackageByID(pkgid string) *meta.Package {
	return cc.metacompiler.PackageByID(pkgid)
}



















