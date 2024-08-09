package helpers

import (
	"strings"
)

const (
	proto_file_template = `syntax = "proto3";

package PACKAGE;

option go_package = "GO_PACKAGE";
option java_package = "JAVA_PACKAGE";

message Foo {
  string Bar=1;
}

`
)

type protoFileBuilder struct {
	pkg               string
	goPackage         string
	javaPackage       string
	customGoPackage   bool
	customJavaPackage bool
	filename          string
}

func NewProtoFileBuilder(pkg string) *protoFileBuilder {
	return &protoFileBuilder{pkg: pkg}
}
func (pfb *protoFileBuilder) GetGoPackage() string {
	if pfb.customGoPackage {
		return pfb.goPackage
	}
	return "golang.conradwood.net/apis/" + pfb.pkg

}
func (pfb *protoFileBuilder) GetJavaPackage() string {
	if pfb.customJavaPackage {
		return pfb.javaPackage
	}
	return "net.conradwood.golang.apis." + pfb.pkg

}
func (pfb *protoFileBuilder) GetFilename() string {
	if pfb.filename != "" {
		return pfb.filename
	}
	return "golang.conradwood.net/apis/" + pfb.pkg + "/" + pfb.pkg + ".proto"
}
func (pfb *protoFileBuilder) SetGoPackage(pkg string) {
	pfb.goPackage = pkg
	pfb.customGoPackage = true
}
func (pfb *protoFileBuilder) SetJavaPackage(pkg string) {
	pfb.javaPackage = pkg
	pfb.customJavaPackage = true
}
func (pfb *protoFileBuilder) SetFilename(filename string) {
	pfb.filename = filename
}

func (pfb *protoFileBuilder) Bytes() []byte {
	bs := proto_file_template
	bs = strings.ReplaceAll(bs, "GO_PACKAGE", pfb.GetGoPackage())
	bs = strings.ReplaceAll(bs, "JAVA_PACKAGE", pfb.GetJavaPackage())
	bs = strings.ReplaceAll(bs, "PACKAGE", pfb.pkg)
	return []byte(bs)
}
