package protoparser

import (
	pr "golang.conradwood.net/apis/protorenderer"
)

type ProtoParsed struct {
	GoPackage   string
	JavaPackage string
	Content     string
	Imports     []string
}

func (pp *ProtoParsed) Protofile() *pr.ProtoFile {
	res := &pr.ProtoFile{
		Content:     pp.Content,
		GoPackage:   pp.GoPackage,
		JavaPackage: pp.JavaPackage,
		Imports:     pp.Imports,
	}
	return res
}






























































































