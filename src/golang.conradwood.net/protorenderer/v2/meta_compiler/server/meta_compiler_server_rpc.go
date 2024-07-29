package server

import (
	"context"
	"fmt"
	google_protobuf "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"golang.conradwood.net/apis/common"
	pb "golang.conradwood.net/apis/protorenderer2"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/protorenderer/cmdline"
	"golang.conradwood.net/protorenderer/v2/helpers"
	mcomp "golang.conradwood.net/protorenderer/v2/meta_compiler"
	"golang.conradwood.net/protorenderer/v2/store"
	"sort"
	"strings"
)

type ServerMetaCompiler struct {
	mc          *mcomp.MetaCompiler
	descriptors map[string]*MessageDescriptor
}
type MessageDescriptor struct {
	ID        uint64
	descproto *google_protobuf.DescriptorProto
}

// called by the protoc plugin
func InternalMetaSubmit(ctx context.Context, req *pb.ProtocRequest) (*common.Void, error) {
	mc, err := mcomp.GetMetaCompilerByID(ctx, req.MetaCompilerID)
	if err != nil {
		fmt.Printf("meta compiler server invoked with an invalid meta compiler id (%s)\n", err)
		return nil, err
	}
	icm := &ServerMetaCompiler{mc: mc, descriptors: make(map[string]*MessageDescriptor)}
	// we want all protofiles in the database, to be able to refer to them by ID
	// we also need a map of all protobuf messages, because within a message we might reference another one
	err = icm.save_request_to_db(ctx, req)
	if err != nil {
		return nil, err
	}
	fmt.Printf("[metacompiler] request with %d protofiles\n", len(req.ProtoFiles))
	files_written := 0
	for _, pf := range req.ProtoFiles {
		was_submitted := icm.mc.WasFilenameSubmitted(*pf.Name)
		was_processed := icm.mc.WasProcessed(*pf.Name)
		fmt.Printf("[metacompiler] Protofile request received: %s (submitted=%v,processed=%v)\n", *pf.Name, was_submitted, was_processed)
		if !was_submitted {
			// don't re-run the meta-compiler for a dependency
			continue
		}
		if was_processed {
			// don't re-run meta for any one file twice
			continue
		}
		info, err := icm.handle_protofile(ctx, pf)
		if err != nil {
			return nil, err
		}
		// create the ProtoFieldInfoText
		err = icm.add_info_text(ctx, info)
		if err != nil {
			return nil, err
		}

		// save the meta information to disk
		y, err := utils.MarshalBytes(info)
		if err != nil {
			return nil, err
		}
		//fmt.Println(string(y))
		save_dir := mc.CompilerEnvironment().CompilerOutDir() + "/info"
		fname := save_dir + "/" + *pf.Name
		fname = strings.TrimSuffix(fname, ".proto")
		fname = fname + ".info"
		err = helpers.WriteFileWithDir(fname, y)
		if err != nil {
			return nil, err
		}
		icm.mc.AddProcessed(*pf.Name)
		files_written++
	}
	fmt.Printf("[metacompiler] files written: %d\n", files_written)

	return &common.Void{}, nil
}
func (smc *ServerMetaCompiler) add_info_text(ctx context.Context, info *pb.ProtoFileInfo) error {
	for _, msg := range info.Messages {
		for _, field := range msg.Fields {
			pft := &pb.ProtoFieldInfoText{
				ModifierString:       fmt.Sprintf("%v", field.Modifier),
				Type1String:          fmt.Sprintf("%v", field.Type1),
				Type2String:          fmt.Sprintf("%v", field.Type2),
				PrimitiveType1String: fmt.Sprintf("%v", field.PrimitiveType1),
				PrimitiveType2String: fmt.Sprintf("%v", field.PrimitiveType2),
			}
			msg := smc.GetMessageDescriptorByID(field.ObjectID1)
			if msg != nil {
				pft.ObjectID1String = *msg.descproto.Name
			}
			msg = smc.GetMessageDescriptorByID(field.ObjectID2)
			if msg != nil {
				pft.ObjectID2String = *msg.descproto.Name
			}
			field.TextualRepresentation = pft
		}
	}
	return nil
}
func (smc *ServerMetaCompiler) GetMessageDescriptorByID(id uint64) *MessageDescriptor {
	for _, msg := range smc.descriptors {
		if msg.ID == id {
			return msg
		}
	}
	return nil
}
func (smc *ServerMetaCompiler) GetMessageDescriptorByFQDN(fqdn string) *MessageDescriptor {
	msg, found := smc.descriptors[fqdn]
	if found {
		return msg
	}
	if cmdline.GetDebugMeta() {
		var names []string
		for k, _ := range smc.descriptors {
			names = append(names, k)
		}
		sort.Slice(names, func(i, j int) bool {
			return names[i] < names[j]
		})
		fmt.Printf("Known proto messages:\n")
		for _, n := range names {
			fmt.Printf(" \"%s\"\n", n)
		}
		fmt.Printf("Not found: \"%s\"\n", fqdn)
	}
	return nil
}
func (icm *ServerMetaCompiler) save_request_to_db(ctx context.Context, req *pb.ProtocRequest) error {
	for _, pf := range req.ProtoFiles {
		if icm.mc.WasProcessed(*pf.Name) {
			continue
		}
		// save the protofile
		proto_file, err := store.GetOrCreateFile(ctx, *pf.Name, 0)
		if err != nil {
			return err
		}

		// save the messages
		for _, ms := range pf.MessageType {
			msg, err := store.GetOrCreateMessage(ctx, proto_file.ID, *ms.Name)
			if err != nil {
				return err
			}
			fqdn := fmt.Sprintf(".%v.%v", *pf.Package, *ms.Name)
			icm.descriptors[fqdn] = &MessageDescriptor{ID: msg.ID, descproto: ms}
		}
	}
	return nil
}
