package compiler

import (
	fl "golang.conradwood.net/protorenderer/v1/filelayouter"
	"golang.conradwood.net/protorenderer/v1/meta"
)

type CompilerCallback interface {
	GetFileLayouter() *fl.FileLayouter
	GetMetaPackageByID(pkgid string) *meta.Package
}




























