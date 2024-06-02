// Code generated by protoc-gen-go.
// source: protos/golang.conradwood.net/apis/protorenderer2/protorenderer2.proto
// DO NOT EDIT!

/*
Package protorenderer2 is a generated protocol buffer package.

It is generated from these files:
	protos/golang.conradwood.net/apis/protorenderer2/protorenderer2.proto

It has these top-level messages:
*/
package protorenderer2

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import protorenderer "golang.conradwood.net/apis/protorenderer"
import _ "golang.conradwood.net/apis/common"
import _ "golang.conradwood.net/apis/h2gproxy"
import _ "github.com/golang/protobuf/protoc-gen-go/descriptor"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for ProtoRenderer2 service

type ProtoRenderer2Client interface {
	// add or update a ".proto" file in the renderers database
	UpdateProto(ctx context.Context, in *protorenderer.AddProtoRequest, opts ...grpc.CallOption) (*protorenderer.AddProtoResponse, error)
}

type protoRenderer2Client struct {
	cc *grpc.ClientConn
}

func NewProtoRenderer2Client(cc *grpc.ClientConn) ProtoRenderer2Client {
	return &protoRenderer2Client{cc}
}

func (c *protoRenderer2Client) UpdateProto(ctx context.Context, in *protorenderer.AddProtoRequest, opts ...grpc.CallOption) (*protorenderer.AddProtoResponse, error) {
	out := new(protorenderer.AddProtoResponse)
	err := grpc.Invoke(ctx, "/protorenderer2.ProtoRenderer2/UpdateProto", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ProtoRenderer2 service

type ProtoRenderer2Server interface {
	// add or update a ".proto" file in the renderers database
	UpdateProto(context.Context, *protorenderer.AddProtoRequest) (*protorenderer.AddProtoResponse, error)
}

func RegisterProtoRenderer2Server(s *grpc.Server, srv ProtoRenderer2Server) {
	s.RegisterService(&_ProtoRenderer2_serviceDesc, srv)
}

func _ProtoRenderer2_UpdateProto_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(protorenderer.AddProtoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProtoRenderer2Server).UpdateProto(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protorenderer2.ProtoRenderer2/UpdateProto",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProtoRenderer2Server).UpdateProto(ctx, req.(*protorenderer.AddProtoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ProtoRenderer2_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protorenderer2.ProtoRenderer2",
	HandlerType: (*ProtoRenderer2Server)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpdateProto",
			Handler:    _ProtoRenderer2_UpdateProto_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protos/golang.conradwood.net/apis/protorenderer2/protorenderer2.proto",
}

func init() {
	proto.RegisterFile("protos/golang.conradwood.net/apis/protorenderer2/protorenderer2.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 215 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x8c, 0x91, 0x31, 0x6f, 0xc2, 0x30,
	0x10, 0x85, 0xb7, 0x0e, 0xae, 0x94, 0x21, 0x63, 0x86, 0xb6, 0x6b, 0x16, 0x5b, 0x72, 0xd7, 0x2e,
	0xad, 0xd4, 0x15, 0x21, 0x24, 0x18, 0x98, 0x70, 0xe2, 0xc3, 0x44, 0x4a, 0x7c, 0xc6, 0x76, 0x04,
	0xfc, 0x7b, 0x14, 0xc7, 0x20, 0x8c, 0x44, 0x94, 0xc9, 0xf7, 0xce, 0xef, 0x9d, 0xbf, 0x93, 0xc9,
	0xbf, 0xb1, 0xe8, 0xd1, 0x31, 0x85, 0xad, 0xd0, 0x8a, 0xd6, 0xa8, 0xad, 0x90, 0x27, 0x44, 0x49,
	0x35, 0x78, 0x26, 0x4c, 0xe3, 0x58, 0x70, 0x58, 0xd0, 0x12, 0x2c, 0x58, 0xfe, 0x24, 0x69, 0x90,
	0x79, 0x96, 0x76, 0x8b, 0x9f, 0xb9, 0xf3, 0x52, 0x35, 0x4e, 0x2b, 0xe8, 0x44, 0xba, 0xc6, 0xae,
	0x43, 0x1d, 0x8f, 0xe8, 0xe7, 0x13, 0xfe, 0x03, 0x57, 0xc6, 0xe2, 0xf9, 0x72, 0x2f, 0x62, 0xe6,
	0x4b, 0x21, 0xaa, 0x16, 0xc6, 0xf7, 0xab, 0x7e, 0xcf, 0x24, 0xb8, 0xda, 0x36, 0xc6, 0x63, 0xa4,
	0xe0, 0x3b, 0x92, 0x2d, 0x87, 0x62, 0x75, 0xdb, 0x2a, 0x5f, 0x90, 0xf7, 0xb5, 0x91, 0xc2, 0x43,
	0xe8, 0xe7, 0x1f, 0x34, 0x85, 0xff, 0x95, 0x32, 0x06, 0x8e, 0x3d, 0x38, 0x5f, 0x7c, 0xbe, 0xbc,
	0x77, 0x06, 0xb5, 0x83, 0xbf, 0x0d, 0x29, 0x35, 0xf8, 0x47, 0xec, 0xb8, 0xc8, 0x40, 0x9e, 0x86,
	0xf9, 0xb6, 0x9c, 0xfd, 0x45, 0xd5, 0x5b, 0xd0, 0xdf, 0xd7, 0x00, 0x00, 0x00, 0xff, 0xff, 0x35,
	0xf6, 0x38, 0x2b, 0xdd, 0x01, 0x00, 0x00,
}
