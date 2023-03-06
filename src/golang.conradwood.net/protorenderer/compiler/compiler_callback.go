package compiler

import (
	fl "golang.conradwood.net/protorenderer/filelayouter"
	"golang.conradwood.net/protorenderer/meta"
)

type CompilerCallback interface {
	GetFileLayouter() *fl.FileLayouter
	GetMetaPackageByID(pkgid string) *meta.Package
}
