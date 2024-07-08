package server

import (
	"fmt"
	google_protobuf "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

func (smc *ServerMetaCompiler) handle_protofile(pf *google_protobuf.FileDescriptorProto) error {
	// comments starting with 'CNW_OPTION: '
	cnwopts := smc.parseCNWOptionsFromFile(pf)
	pkg_name := *pf.Package
	file_name := *pf.Name
	smc.debugf("Package \"%s\" in file \"%s\", %d cnwopts\n", pkg_name, file_name, len(cnwopts))
	// now find all messages
	for _, ms := range pf.MessageType {
		smc.debugf("   Message: %s\n", *ms.Name)
	}
	return nil
}

func (smc *ServerMetaCompiler) debugf(format string, args ...interface{}) {
	prefix := fmt.Sprintf("[servermetacompiler] ")
	text := fmt.Sprintf(format, args...)
	fmt.Printf("%s%s", prefix, text)
}
