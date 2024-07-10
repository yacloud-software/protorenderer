package server

import (
	"context"
	"fmt"
	google_protobuf "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"golang.conradwood.net/apis/common"
	pb "golang.conradwood.net/apis/protorenderer2"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/protorenderer/v2/helpers"
	mcomp "golang.conradwood.net/protorenderer/v2/meta_compiler"
	"sort"
	"strings"
)

type ServerMetaCompiler struct {
	mc          *mcomp.MetaCompiler
	descriptors map[string]*MessageDescriptor
}
type MessageDescriptor struct {
	id        uint64
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
	for _, pf := range req.ProtoFiles {
		fmt.Printf("meta compiler: %#v\n", pf)
	}
	// we need a map of all protobuf messages, because within a message we might reference another one
	// this needs to be resolved first so that we can later resolve fields properly
	for _, pf := range req.ProtoFiles {
		for _, ms := range pf.MessageType {
			fqdn := fmt.Sprintf(".%v.%v", *pf.Package, *ms.Name)
			icm.descriptors[fqdn] = &MessageDescriptor{descproto: ms}
		}
	}
	for _, pf := range req.ProtoFiles {
		fmt.Printf("Protofile: %s\n", *pf.Name)
		info, err := icm.handle_protofile(ctx, pf)
		if err != nil {
			return nil, err
		}
		y, err := utils.MarshalYaml(info)
		if err != nil {
			return nil, err
		}
		//fmt.Println(string(y))
		save_dir := mc.CompilerEnvironment().ResultsDir() + "/info"
		fname := save_dir + "/" + *pf.Name
		fname = strings.TrimSuffix(fname, ".proto")
		fname = fname + ".info"
		err = helpers.WriteFileWithDir(fname, y)
		if err != nil {
			return nil, err
		}
	}

	return &common.Void{}, nil
}

func (smc *ServerMetaCompiler) GetMessageDescriptorByFQDN(fqdn string) *MessageDescriptor {
	msg, found := smc.descriptors[fqdn]
	if found {
		return msg
	}
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
	return nil
}
