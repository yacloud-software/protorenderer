package server

import (
	"context"
	"fmt"
	google_protobuf "github.com/golang/protobuf/protoc-gen-go/descriptor"
	pb "golang.conradwood.net/apis/protorenderer2"
	"golang.conradwood.net/protorenderer/cmdline"
	"golang.conradwood.net/protorenderer/v2/store"
)

func (smc *ServerMetaCompiler) handle_protofile(ctx context.Context, pf *google_protobuf.FileDescriptorProto) (*pb.ProtoFileInfo, error) {
	// comments starting with 'CNW_OPTION: '
	cnwopts := smc.parseCNWOptionsFromFile(pf)
	pkg_name := *pf.Package
	file_name := *pf.Name
	proto_file, err := store.FileByName(ctx, file_name)
	if err != nil {
		return nil, err
	}
	if proto_file.Package == "" {
		err = store.UpdatePackageInFile(ctx, proto_file, pkg_name)
		if err != nil {
			return nil, err
		}
	}
	res := &pb.ProtoFileInfo{
		ProtoFile:  proto_file,
		CNWOptions: cnwopts,
	}
	for _, imp := range pf.Dependency {
		pf, err := store.FileByName(ctx, imp)
		if err != nil {
			return nil, err
		}
		res.Imports = append(res.Imports, pf)
	}
	comments, err := smc.dealWithComments(pf)
	if err != nil {
		return nil, err
	}
	smc.debugf("Package \"%s\" in file \"%s\" (#%d), %d cnwopts\n", pkg_name, file_name, proto_file.GetID(), len(cnwopts))
	// now find all messages
	for _, ms := range pf.MessageType {
		msg, err := store.GetOrCreateMessage(ctx, proto_file.ID, *ms.Name)
		if err != nil {
			return nil, err
		}
		smc.debugf("   Message %d: %s\n", msg.ID, msg.Name)
		pmi := &pb.ProtoMessageInfo{
			Message: msg,
			Comment: comments.GetMessageComment(msg.Name),
		}
		res.Messages = append(res.Messages, pmi)
		for _, f := range ms.Field {
			pfi, err := smc.createFieldInfo(msg, f)
			if err != nil {
				return nil, err
			}
			pfi.Comment = comments.GetFieldComment(msg.Name, *f.Name)

			pmi.Fields = append(pmi.Fields, pfi)
		}
	}
	return res, nil
}

// creates field info (but does not add comment)
func (smc *ServerMetaCompiler) createFieldInfo(m *pb.SQLMessage, fd *google_protobuf.FieldDescriptorProto) (*pb.ProtoFieldInfo, error) {
	pfi := &pb.ProtoFieldInfo{
		Name: *fd.Name,
	}

	is_primitive := true
	if fd.Type != nil && *fd.Type == google_protobuf.FieldDescriptorProto_TYPE_MESSAGE {
		is_primitive = false
	}
	pfi.Modifier = pb.FieldModifier_SINGLE

	smc.debugf("Field %s primitive: %v, label=%v\n", *fd.Name, is_primitive, *fd.Label)

	// check if it is an arry
	if fd.Label != nil {
		fdl := *fd.Label
		if fdl == google_protobuf.FieldDescriptorProto_LABEL_REPEATED {
			pfi.Modifier = pb.FieldModifier_ARRAY
		}
	}

	// handle the easy part, the string, uint etc..
	if is_primitive {
		typ, err := protoc_type_to_protorenderer_type(fd.Type)
		if err != nil {
			return nil, err
		}
		pfi.Type1 = pb.FieldType_PRIMITIVE
		pfi.PrimitiveType1 = typ
		return pfi, nil
	}

	msgfqdn := *fd.TypeName

	if pfi.Modifier == pb.FieldModifier_SINGLE {
		// is is a single (neither array nor map) field of complex type
		msg := smc.GetMessageDescriptorByFQDN(msgfqdn)
		if msg == nil {
			return nil, fmt.Errorf("No descriptor for message \"%s\"", msgfqdn)
		}
		if msg.ID == 0 {
			return nil, fmt.Errorf("message \"%s\" must not have ID 0", msgfqdn)
		}
		pfi.Type1 = pb.FieldType_OBJECT
		pfi.ObjectID1 = msg.ID
		return pfi, nil
	}
	// it is either a array of messages or a map. to determine which one use getTypeName()
	// https://github.com/protocolbuffers/protobuf-javascript/issues/13
	// unknown stuff:
	// if the messagetype refers to an unresolvable protobuf we assume it is a map (it's a bit rubbish, but found no better way)
	// could not find a way to resolve the typename to the internal MapEntry protobuf
	// there is a protoc generator https://github.com/golang/protobuf/blob/master/protoc-gen-go/generator/generator.go#L68
	// but it is deprecated.

	ismap := false
	msg := smc.GetMessageDescriptorByFQDN(msgfqdn)
	if msg == nil {
		ismap = true
	}

	smc.debugf("Field %s.%s map or complex array (type=%v, typename=%v, label=%v)\n", m.Name, *fd.Name, *fd.Type, *fd.TypeName, *fd.Label)

	if !ismap {
		// it is an _array_ of complex types (not a map)
		pfi.Type1 = pb.FieldType_OBJECT
		pfi.ObjectID1 = msg.ID
		return pfi, nil
	}
	if fd.TypeName != nil {
		smc.debugf("            Typename: %s\n", *fd.TypeName)
	}
	if fd.Type != nil {
		smc.debugf("            Type    : %v\n", *fd.Type)
	}

	// it is a map.

	// TODO: figure out who we can get to the elements of the map
	// once that is figured out, the IDs of the messages need to be set in ObjectID1 and ObjectID2
	pfi.Modifier = pb.FieldModifier_MAP
	pfi.Type1 = pb.FieldType_ft_UNDEFINED
	pfi.Type2 = pb.FieldType_ft_UNDEFINED
	return pfi, nil
}
func protoc_type_to_protorenderer_type(typ *google_protobuf.FieldDescriptorProto_Type) (pb.ProtoFieldPrimitive, error) {
	t := *typ
	if t == google_protobuf.FieldDescriptorProto_TYPE_STRING {
		return pb.ProtoFieldPrimitive_STRING, nil
	} else if t == google_protobuf.FieldDescriptorProto_TYPE_INT32 {
		return pb.ProtoFieldPrimitive_INT32, nil
	} else if t == google_protobuf.FieldDescriptorProto_TYPE_ENUM {
		return pb.ProtoFieldPrimitive_ENUM, nil
	} else if t == google_protobuf.FieldDescriptorProto_TYPE_UINT64 {
		return pb.ProtoFieldPrimitive_UINT64, nil
	} else if t == google_protobuf.FieldDescriptorProto_TYPE_UINT32 {
		return pb.ProtoFieldPrimitive_UINT32, nil
	} else if t == google_protobuf.FieldDescriptorProto_TYPE_BYTES {
		return pb.ProtoFieldPrimitive_BYTES, nil
	} else if t == google_protobuf.FieldDescriptorProto_TYPE_BOOL {
		return pb.ProtoFieldPrimitive_BOOL, nil
	} else if t == google_protobuf.FieldDescriptorProto_TYPE_DOUBLE {
		return pb.ProtoFieldPrimitive_DOUBLE, nil
	} else if t == google_protobuf.FieldDescriptorProto_TYPE_FLOAT {
		return pb.ProtoFieldPrimitive_FLOAT, nil
	} else if t == google_protobuf.FieldDescriptorProto_TYPE_INT64 {
		return pb.ProtoFieldPrimitive_INT64, nil
	} else {
		return 0, fmt.Errorf("unknown protoc field type (%v)", t)
	}
}
func (smc *ServerMetaCompiler) debugf(format string, args ...interface{}) {
	if !cmdline.GetDebugMeta() {
		return
	}
	prefix := fmt.Sprintf("[servermetacompiler] ")
	text := fmt.Sprintf(format, args...)
	fmt.Printf("%s%s", prefix, text)
}
