// Code generated by protoc-gen-go.
// source: golang.conradwood.net/apis/ipexporter/ipexporter.proto
// DO NOT EDIT!

/*
Package ipexporter is a generated protocol buffer package.

It is generated from these files:
	golang.conradwood.net/apis/ipexporter/ipexporter.proto

It has these top-level messages:
	PingRequest
	PingResponse
*/
package ipexporter

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

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

// comment: message pingrequest
type PingRequest struct {
	// comment: payload
	Payload string `protobuf:"bytes,2,opt,name=Payload" json:"Payload,omitempty"`
	// comment: sequencenumber
	SequenceNumber uint32 `protobuf:"varint,1,opt,name=SequenceNumber" json:"SequenceNumber,omitempty"`
}

func (m *PingRequest) Reset()                    { *m = PingRequest{} }
func (m *PingRequest) String() string            { return proto.CompactTextString(m) }
func (*PingRequest) ProtoMessage()               {}
func (*PingRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *PingRequest) GetPayload() string {
	if m != nil {
		return m.Payload
	}
	return ""
}

func (m *PingRequest) GetSequenceNumber() uint32 {
	if m != nil {
		return m.SequenceNumber
	}
	return 0
}

// comment: message pingresponse
type PingResponse struct {
	// comment: field pingresponse.response
	Response *PingRequest `protobuf:"bytes,1,opt,name=Response" json:"Response,omitempty"`
}

func (m *PingResponse) Reset()                    { *m = PingResponse{} }
func (m *PingResponse) String() string            { return proto.CompactTextString(m) }
func (*PingResponse) ProtoMessage()               {}
func (*PingResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *PingResponse) GetResponse() *PingRequest {
	if m != nil {
		return m.Response
	}
	return nil
}

func init() {
	proto.RegisterType((*PingRequest)(nil), "ipexporter.PingRequest")
	proto.RegisterType((*PingResponse)(nil), "ipexporter.PingResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for IPExporterService service

type IPExporterServiceClient interface {
	// comment: rpc ping
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
}

type iPExporterServiceClient struct {
	cc *grpc.ClientConn
}

func NewIPExporterServiceClient(cc *grpc.ClientConn) IPExporterServiceClient {
	return &iPExporterServiceClient{cc}
}

func (c *iPExporterServiceClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := grpc.Invoke(ctx, "/ipexporter.IPExporterService/Ping", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for IPExporterService service

type IPExporterServiceServer interface {
	// comment: rpc ping
	Ping(context.Context, *PingRequest) (*PingResponse, error)
}

func RegisterIPExporterServiceServer(s *grpc.Server, srv IPExporterServiceServer) {
	s.RegisterService(&_IPExporterService_serviceDesc, srv)
}

func _IPExporterService_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IPExporterServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ipexporter.IPExporterService/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IPExporterServiceServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _IPExporterService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ipexporter.IPExporterService",
	HandlerType: (*IPExporterServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _IPExporterService_Ping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "golang.conradwood.net/apis/ipexporter/ipexporter.proto",
}

// Client API for EchoStreamService service

type EchoStreamServiceClient interface {
	// comment: rpc SendToServer
	SendToServer(ctx context.Context, opts ...grpc.CallOption) (EchoStreamService_SendToServerClient, error)
}

type echoStreamServiceClient struct {
	cc *grpc.ClientConn
}

func NewEchoStreamServiceClient(cc *grpc.ClientConn) EchoStreamServiceClient {
	return &echoStreamServiceClient{cc}
}

func (c *echoStreamServiceClient) SendToServer(ctx context.Context, opts ...grpc.CallOption) (EchoStreamService_SendToServerClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_EchoStreamService_serviceDesc.Streams[0], c.cc, "/ipexporter.EchoStreamService/SendToServer", opts...)
	if err != nil {
		return nil, err
	}
	x := &echoStreamServiceSendToServerClient{stream}
	return x, nil
}

type EchoStreamService_SendToServerClient interface {
	Send(*PingRequest) error
	CloseAndRecv() (*PingResponse, error)
	grpc.ClientStream
}

type echoStreamServiceSendToServerClient struct {
	grpc.ClientStream
}

func (x *echoStreamServiceSendToServerClient) Send(m *PingRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *echoStreamServiceSendToServerClient) CloseAndRecv() (*PingResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(PingResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for EchoStreamService service

type EchoStreamServiceServer interface {
	// comment: rpc SendToServer
	SendToServer(EchoStreamService_SendToServerServer) error
}

func RegisterEchoStreamServiceServer(s *grpc.Server, srv EchoStreamServiceServer) {
	s.RegisterService(&_EchoStreamService_serviceDesc, srv)
}

func _EchoStreamService_SendToServer_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(EchoStreamServiceServer).SendToServer(&echoStreamServiceSendToServerServer{stream})
}

type EchoStreamService_SendToServerServer interface {
	SendAndClose(*PingResponse) error
	Recv() (*PingRequest, error)
	grpc.ServerStream
}

type echoStreamServiceSendToServerServer struct {
	grpc.ServerStream
}

func (x *echoStreamServiceSendToServerServer) SendAndClose(m *PingResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *echoStreamServiceSendToServerServer) Recv() (*PingRequest, error) {
	m := new(PingRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _EchoStreamService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ipexporter.EchoStreamService",
	HandlerType: (*EchoStreamServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SendToServer",
			Handler:       _EchoStreamService_SendToServer_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "golang.conradwood.net/apis/ipexporter/ipexporter.proto",
}

func init() {
	proto.RegisterFile("golang.conradwood.net/apis/ipexporter/ipexporter.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 250 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x9c, 0x91, 0xcd, 0x4b, 0xc3, 0x40,
	0x10, 0xc5, 0x89, 0x88, 0x1f, 0xd3, 0x2a, 0x74, 0x2f, 0x06, 0x4f, 0xa5, 0xa0, 0xe6, 0x94, 0x42,
	0x0b, 0x82, 0x57, 0x4b, 0x0f, 0x5e, 0x6a, 0x48, 0x3c, 0x78, 0xdd, 0x26, 0x43, 0x0c, 0xb4, 0x33,
	0xeb, 0xec, 0xd6, 0x8f, 0xff, 0x5e, 0x36, 0x26, 0xba, 0x28, 0x5e, 0xbc, 0xcd, 0xbe, 0x7d, 0xef,
	0xb7, 0xcb, 0x3c, 0xb8, 0xae, 0x79, 0xa3, 0xa9, 0x4e, 0x4b, 0x26, 0xd1, 0xd5, 0x2b, 0x73, 0x95,
	0x12, 0xba, 0xa9, 0x36, 0x8d, 0x9d, 0x36, 0x06, 0xdf, 0x0c, 0x8b, 0x43, 0x09, 0xc6, 0xd4, 0x08,
	0x3b, 0x56, 0xf0, 0xad, 0x4c, 0xee, 0x61, 0x90, 0x35, 0x54, 0xe7, 0xf8, 0xbc, 0x43, 0xeb, 0x54,
	0x0c, 0x87, 0x99, 0x7e, 0xdf, 0xb0, 0xae, 0xe2, 0xbd, 0x71, 0x94, 0x1c, 0xe7, 0xfd, 0x51, 0x5d,
	0xc2, 0x69, 0xe1, 0x4d, 0x54, 0xe2, 0x6a, 0xb7, 0x5d, 0xa3, 0xc4, 0xd1, 0x38, 0x4a, 0x4e, 0xf2,
	0x1f, 0xea, 0x64, 0x01, 0xc3, 0x4f, 0xa0, 0x35, 0x4c, 0x16, 0xd5, 0x1c, 0x8e, 0xfa, 0xb9, 0x4d,
	0x0c, 0x66, 0x67, 0x69, 0xf0, 0xa3, 0xe0, 0xf1, 0xfc, 0xcb, 0x38, 0x5b, 0xc1, 0xe8, 0x2e, 0x5b,
	0x76, 0x9e, 0x02, 0xe5, 0xa5, 0x29, 0x51, 0xdd, 0xc0, 0xbe, 0x77, 0xab, 0xbf, 0xf2, 0xe7, 0xf1,
	0xef, 0x8b, 0x8e, 0xf7, 0x08, 0xa3, 0x65, 0xf9, 0xc4, 0x85, 0x13, 0xd4, 0xdb, 0x9e, 0xb7, 0x80,
	0x61, 0x81, 0x54, 0x3d, 0xb0, 0x17, 0x50, 0xfe, 0xc1, 0x4d, 0xa2, 0xdb, 0x2b, 0xb8, 0x20, 0x74,
	0x61, 0x05, 0x5d, 0x29, 0xbe, 0x85, 0x20, 0xb7, 0x3e, 0x68, 0x77, 0x3f, 0xff, 0x08, 0x00, 0x00,
	0xff, 0xff, 0x14, 0x9c, 0x2a, 0x44, 0xb5, 0x01, 0x00, 0x00,
}
