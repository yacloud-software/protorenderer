package server

import (
	"context"
	"fmt"
	google_protobuf "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"golang.conradwood.net/protorenderer/v2/store"
)

func (smc *ServerMetaCompiler) handle_protofile(ctx context.Context, pf *google_protobuf.FileDescriptorProto) error {
	// comments starting with 'CNW_OPTION: '
	cnwopts := smc.parseCNWOptionsFromFile(pf)
	pkg_name := *pf.Package
	file_name := *pf.Name
	proto_file, err := store.FileByName(ctx, file_name)
	if err != nil {
		return err
	}
	if proto_file.Package == "" {
		err = store.UpdatePackageInFile(ctx, proto_file, pkg_name)
		if err != nil {
			return err
		}
	}
	smc.debugf("Package \"%s\" in file \"%s\" (#%d), %d cnwopts\n", pkg_name, file_name, proto_file.GetID(), len(cnwopts))
	// now find all messages
	for _, ms := range pf.MessageType {
		msg, err := store.GetOrCreateMessage(ctx, proto_file.ID, *ms.Name)
		if err != nil {
			return err
		}
		smc.debugf("   Message %d: %s\n", msg.ID, msg.Name)
	}
	return nil
}

func (smc *ServerMetaCompiler) debugf(format string, args ...interface{}) {
	prefix := fmt.Sprintf("[servermetacompiler] ")
	text := fmt.Sprintf(format, args...)
	fmt.Printf("%s%s", prefix, text)
}
