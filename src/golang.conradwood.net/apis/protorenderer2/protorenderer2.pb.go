// Code generated by protoc-gen-go.
// source: protos/golang.conradwood.net/apis/protorenderer2/protorenderer2.proto
// DO NOT EDIT!

/*
Package protorenderer2 is a generated protocol buffer package.

It is generated from these files:
	protos/golang.conradwood.net/apis/protorenderer2/protorenderer2.proto

It has these top-level messages:
	ProtocRequest
	FileWithContent
*/
package protorenderer2

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import protorenderer "golang.conradwood.net/apis/protorenderer"
import common "golang.conradwood.net/apis/common"
import _ "golang.conradwood.net/apis/h2gproxy"
import google_protobuf "github.com/golang/protobuf/protoc-gen-go/descriptor"

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

type ProtocRequest struct {
	VerifyToken string                                 `protobuf:"bytes,1,opt,name=VerifyToken" json:"VerifyToken,omitempty"`
	ProtoFiles  []*google_protobuf.FileDescriptorProto `protobuf:"bytes,2,rep,name=ProtoFiles" json:"ProtoFiles,omitempty"`
}

func (m *ProtocRequest) Reset()                    { *m = ProtocRequest{} }
func (m *ProtocRequest) String() string            { return proto.CompactTextString(m) }
func (*ProtocRequest) ProtoMessage()               {}
func (*ProtocRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ProtocRequest) GetVerifyToken() string {
	if m != nil {
		return m.VerifyToken
	}
	return ""
}

func (m *ProtocRequest) GetProtoFiles() []*google_protobuf.FileDescriptorProto {
	if m != nil {
		return m.ProtoFiles
	}
	return nil
}

type FileWithContent struct {
	Filename string `protobuf:"bytes,1,opt,name=Filename" json:"Filename,omitempty"`
	Data     []byte `protobuf:"bytes,2,opt,name=Data,proto3" json:"Data,omitempty"`
}

func (m *FileWithContent) Reset()                    { *m = FileWithContent{} }
func (m *FileWithContent) String() string            { return proto.CompactTextString(m) }
func (*FileWithContent) ProtoMessage()               {}
func (*FileWithContent) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *FileWithContent) GetFilename() string {
	if m != nil {
		return m.Filename
	}
	return ""
}

func (m *FileWithContent) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*ProtocRequest)(nil), "protorenderer2.ProtocRequest")
	proto.RegisterType((*FileWithContent)(nil), "protorenderer2.FileWithContent")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for ProtoRenderer2 service

type ProtoRenderer2Client interface {
	// protoc-gen-meta2 submits files to protorenderer server for analysis
	SubmitSource(ctx context.Context, in *ProtocRequest, opts ...grpc.CallOption) (*common.Void, error)
	// compile
	Compile(ctx context.Context, opts ...grpc.CallOption) (ProtoRenderer2_CompileClient, error)
	// add or update a ".proto" file in the renderers database
	UpdateProto(ctx context.Context, in *protorenderer.AddProtoRequest, opts ...grpc.CallOption) (*protorenderer.AddProtoResponse, error)
}

type protoRenderer2Client struct {
	cc *grpc.ClientConn
}

func NewProtoRenderer2Client(cc *grpc.ClientConn) ProtoRenderer2Client {
	return &protoRenderer2Client{cc}
}

func (c *protoRenderer2Client) SubmitSource(ctx context.Context, in *ProtocRequest, opts ...grpc.CallOption) (*common.Void, error) {
	out := new(common.Void)
	err := grpc.Invoke(ctx, "/protorenderer2.ProtoRenderer2/SubmitSource", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *protoRenderer2Client) Compile(ctx context.Context, opts ...grpc.CallOption) (ProtoRenderer2_CompileClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_ProtoRenderer2_serviceDesc.Streams[0], c.cc, "/protorenderer2.ProtoRenderer2/Compile", opts...)
	if err != nil {
		return nil, err
	}
	x := &protoRenderer2CompileClient{stream}
	return x, nil
}

type ProtoRenderer2_CompileClient interface {
	Send(*FileWithContent) error
	Recv() (*FileWithContent, error)
	grpc.ClientStream
}

type protoRenderer2CompileClient struct {
	grpc.ClientStream
}

func (x *protoRenderer2CompileClient) Send(m *FileWithContent) error {
	return x.ClientStream.SendMsg(m)
}

func (x *protoRenderer2CompileClient) Recv() (*FileWithContent, error) {
	m := new(FileWithContent)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
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
	// protoc-gen-meta2 submits files to protorenderer server for analysis
	SubmitSource(context.Context, *ProtocRequest) (*common.Void, error)
	// compile
	Compile(ProtoRenderer2_CompileServer) error
	// add or update a ".proto" file in the renderers database
	UpdateProto(context.Context, *protorenderer.AddProtoRequest) (*protorenderer.AddProtoResponse, error)
}

func RegisterProtoRenderer2Server(s *grpc.Server, srv ProtoRenderer2Server) {
	s.RegisterService(&_ProtoRenderer2_serviceDesc, srv)
}

func _ProtoRenderer2_SubmitSource_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProtocRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProtoRenderer2Server).SubmitSource(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protorenderer2.ProtoRenderer2/SubmitSource",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProtoRenderer2Server).SubmitSource(ctx, req.(*ProtocRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProtoRenderer2_Compile_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ProtoRenderer2Server).Compile(&protoRenderer2CompileServer{stream})
}

type ProtoRenderer2_CompileServer interface {
	Send(*FileWithContent) error
	Recv() (*FileWithContent, error)
	grpc.ServerStream
}

type protoRenderer2CompileServer struct {
	grpc.ServerStream
}

func (x *protoRenderer2CompileServer) Send(m *FileWithContent) error {
	return x.ServerStream.SendMsg(m)
}

func (x *protoRenderer2CompileServer) Recv() (*FileWithContent, error) {
	m := new(FileWithContent)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
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
			MethodName: "SubmitSource",
			Handler:    _ProtoRenderer2_SubmitSource_Handler,
		},
		{
			MethodName: "UpdateProto",
			Handler:    _ProtoRenderer2_UpdateProto_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Compile",
			Handler:       _ProtoRenderer2_Compile_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "protos/golang.conradwood.net/apis/protorenderer2/protorenderer2.proto",
}

func init() {
	proto.RegisterFile("protos/golang.conradwood.net/apis/protorenderer2/protorenderer2.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 369 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x8c, 0x52, 0x4d, 0x6f, 0xe2, 0x30,
	0x10, 0x55, 0xd8, 0xd5, 0x7e, 0x18, 0x96, 0x95, 0x7c, 0x42, 0x91, 0x76, 0x89, 0xd0, 0x1e, 0xc2,
	0xc5, 0x59, 0xa5, 0xc7, 0xf6, 0x42, 0xa1, 0x3d, 0xb6, 0x55, 0x68, 0xa9, 0xd4, 0x5b, 0x88, 0x87,
	0x60, 0x35, 0xf1, 0xa4, 0x8e, 0x23, 0xca, 0x9f, 0xee, 0x6f, 0xa8, 0x92, 0x18, 0x44, 0x22, 0x95,
	0x72, 0xca, 0xbc, 0xf1, 0xbc, 0x79, 0x2f, 0xf6, 0x23, 0x57, 0x99, 0x42, 0x8d, 0xb9, 0x17, 0x63,
	0x12, 0xca, 0x98, 0x45, 0x28, 0x55, 0xc8, 0x37, 0x88, 0x9c, 0x49, 0xd0, 0x5e, 0x98, 0x89, 0xdc,
	0xab, 0x26, 0x14, 0x48, 0x0e, 0x0a, 0x94, 0xdf, 0x82, 0xac, 0x82, 0xb4, 0xdf, 0xec, 0xda, 0x17,
	0xa7, 0xee, 0x6b, 0xa2, 0x7a, 0x9b, 0xcd, 0x8e, 0xb0, 0x23, 0x4c, 0x53, 0x94, 0xe6, 0x63, 0xe6,
	0xfd, 0x23, 0xf3, 0x6b, 0x3f, 0xce, 0x14, 0xbe, 0x6e, 0xf7, 0x85, 0xe1, 0x38, 0x31, 0x62, 0x9c,
	0x40, 0xad, 0xbf, 0x2c, 0x56, 0x1e, 0x87, 0x3c, 0x52, 0x22, 0xd3, 0x68, 0x5c, 0x8c, 0x36, 0xe4,
	0xd7, 0x5d, 0x59, 0x44, 0x01, 0xbc, 0x14, 0x90, 0x6b, 0xea, 0x90, 0xee, 0x02, 0x94, 0x58, 0x6d,
	0xef, 0xf1, 0x19, 0xe4, 0xc0, 0x72, 0x2c, 0xf7, 0x67, 0x70, 0xd8, 0xa2, 0x33, 0x42, 0x2a, 0xca,
	0xb5, 0x48, 0x20, 0x1f, 0x74, 0x9c, 0x2f, 0x6e, 0xd7, 0xff, 0xc7, 0x6a, 0x25, 0xb6, 0x53, 0x62,
	0xe5, 0xe9, 0x6c, 0xaf, 0x56, 0x11, 0x82, 0x03, 0xde, 0x68, 0x42, 0x7e, 0x97, 0xc5, 0xa3, 0xd0,
	0xeb, 0x29, 0x4a, 0x0d, 0x52, 0x53, 0x9b, 0xfc, 0x28, 0x5b, 0x32, 0x4c, 0xc1, 0xe8, 0xee, 0x31,
	0xa5, 0xe4, 0xeb, 0x2c, 0xd4, 0xe1, 0xa0, 0xe3, 0x58, 0x6e, 0x2f, 0xa8, 0x6a, 0xff, 0xcd, 0x22,
	0xfd, 0x7a, 0xf1, 0xee, 0x49, 0xe8, 0x39, 0xe9, 0xcd, 0x8b, 0x65, 0x2a, 0xf4, 0x1c, 0x0b, 0x15,
	0x01, 0xfd, 0xc3, 0x5a, 0x2f, 0xd9, 0xf8, 0x59, 0xbb, 0xc7, 0xcc, 0x15, 0x2f, 0x50, 0x70, 0x7a,
	0x4b, 0xbe, 0x4f, 0x31, 0xcd, 0x44, 0x02, 0x74, 0xd8, 0xe6, 0xb5, 0xbc, 0xda, 0x9f, 0x0d, 0xb8,
	0xd6, 0x7f, 0x8b, 0xde, 0x90, 0xee, 0x43, 0xc6, 0x43, 0x0d, 0x95, 0x2a, 0xfd, 0xdb, 0xe4, 0xb0,
	0x09, 0xe7, 0xc6, 0x7e, 0xed, 0x66, 0xf8, 0xe1, 0x79, 0x9e, 0xa1, 0xcc, 0xe1, 0x72, 0x41, 0xc6,
	0x12, 0xf4, 0x61, 0x02, 0x4c, 0x26, 0xca, 0x10, 0xb4, 0x0c, 0x3d, 0x8d, 0x4f, 0x4e, 0xfb, 0xf2,
	0x5b, 0x85, 0xcf, 0xde, 0x03, 0x00, 0x00, 0xff, 0xff, 0x01, 0xd8, 0x9d, 0x2c, 0x28, 0x03, 0x00,
	0x00,
}
