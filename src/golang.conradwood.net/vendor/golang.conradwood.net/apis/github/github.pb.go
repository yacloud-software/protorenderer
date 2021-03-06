// Code generated by protoc-gen-go.
// source: golang.conradwood.net/apis/github/github.proto
// DO NOT EDIT!

/*
Package github is a generated protocol buffer package.

It is generated from these files:
	golang.conradwood.net/apis/github/github.proto

It has these top-level messages:
	PingResponse
*/
package github

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import common "golang.conradwood.net/apis/common"
import h2gproxy "golang.conradwood.net/apis/h2gproxy"

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

// comment: message pingresponse
type PingResponse struct {
	// comment: field pingresponse.response
	Response string `protobuf:"bytes,1,opt,name=Response" json:"Response,omitempty"`
}

func (m *PingResponse) Reset()                    { *m = PingResponse{} }
func (m *PingResponse) String() string            { return proto.CompactTextString(m) }
func (*PingResponse) ProtoMessage()               {}
func (*PingResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *PingResponse) GetResponse() string {
	if m != nil {
		return m.Response
	}
	return ""
}

func init() {
	proto.RegisterType((*PingResponse)(nil), "github.PingResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for GitHub service

type GitHubClient interface {
	// comment: rpc ping
	Ping(ctx context.Context, in *common.Void, opts ...grpc.CallOption) (*PingResponse, error)
	ServeHTML(ctx context.Context, in *h2gproxy.ServeRequest, opts ...grpc.CallOption) (*h2gproxy.ServeResponse, error)
}

type gitHubClient struct {
	cc *grpc.ClientConn
}

func NewGitHubClient(cc *grpc.ClientConn) GitHubClient {
	return &gitHubClient{cc}
}

func (c *gitHubClient) Ping(ctx context.Context, in *common.Void, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := grpc.Invoke(ctx, "/github.GitHub/Ping", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gitHubClient) ServeHTML(ctx context.Context, in *h2gproxy.ServeRequest, opts ...grpc.CallOption) (*h2gproxy.ServeResponse, error) {
	out := new(h2gproxy.ServeResponse)
	err := grpc.Invoke(ctx, "/github.GitHub/ServeHTML", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for GitHub service

type GitHubServer interface {
	// comment: rpc ping
	Ping(context.Context, *common.Void) (*PingResponse, error)
	ServeHTML(context.Context, *h2gproxy.ServeRequest) (*h2gproxy.ServeResponse, error)
}

func RegisterGitHubServer(s *grpc.Server, srv GitHubServer) {
	s.RegisterService(&_GitHub_serviceDesc, srv)
}

func _GitHub_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(common.Void)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GitHubServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/github.GitHub/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GitHubServer).Ping(ctx, req.(*common.Void))
	}
	return interceptor(ctx, in, info, handler)
}

func _GitHub_ServeHTML_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(h2gproxy.ServeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GitHubServer).ServeHTML(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/github.GitHub/ServeHTML",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GitHubServer).ServeHTML(ctx, req.(*h2gproxy.ServeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _GitHub_serviceDesc = grpc.ServiceDesc{
	ServiceName: "github.GitHub",
	HandlerType: (*GitHubServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _GitHub_Ping_Handler,
		},
		{
			MethodName: "ServeHTML",
			Handler:    _GitHub_ServeHTML_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "golang.conradwood.net/apis/github/github.proto",
}

func init() { proto.RegisterFile("golang.conradwood.net/apis/github/github.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 214 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xd2, 0x4b, 0xcf, 0xcf, 0x49,
	0xcc, 0x4b, 0xd7, 0x4b, 0xce, 0xcf, 0x2b, 0x4a, 0x4c, 0x29, 0xcf, 0xcf, 0x4f, 0xd1, 0xcb, 0x4b,
	0x2d, 0xd1, 0x4f, 0x2c, 0xc8, 0x2c, 0xd6, 0x4f, 0xcf, 0x2c, 0xc9, 0x28, 0x4d, 0x82, 0x52, 0x7a,
	0x05, 0x45, 0xf9, 0x25, 0xf9, 0x42, 0x6c, 0x10, 0x9e, 0x14, 0x3e, 0x7d, 0xc9, 0xf9, 0xb9, 0xb9,
	0xf9, 0x79, 0x50, 0x0a, 0xa2, 0x4f, 0xca, 0x08, 0x8f, 0xfa, 0x0c, 0xa3, 0xf4, 0x82, 0xa2, 0xfc,
	0x8a, 0x4a, 0x38, 0x03, 0xa2, 0x47, 0x49, 0x8b, 0x8b, 0x27, 0x20, 0x33, 0x2f, 0x3d, 0x28, 0xb5,
	0xb8, 0x20, 0x3f, 0xaf, 0x38, 0x55, 0x48, 0x8a, 0x8b, 0x03, 0xc6, 0x96, 0x60, 0x54, 0x60, 0xd4,
	0xe0, 0x0c, 0x82, 0xf3, 0x8d, 0x8a, 0xb8, 0xd8, 0xdc, 0x33, 0x4b, 0x3c, 0x4a, 0x93, 0x84, 0xb4,
	0xb8, 0x58, 0x40, 0xba, 0x84, 0x78, 0xf4, 0xa0, 0x0e, 0x08, 0xcb, 0xcf, 0x4c, 0x91, 0x12, 0xd1,
	0x83, 0x7a, 0x03, 0xc5, 0x44, 0x1b, 0x2e, 0xce, 0xe0, 0xd4, 0xa2, 0xb2, 0x54, 0x8f, 0x10, 0x5f,
	0x1f, 0x21, 0x31, 0x3d, 0xb8, 0xfd, 0x60, 0xc1, 0xa0, 0xd4, 0xc2, 0xd2, 0xd4, 0xe2, 0x12, 0x29,
	0x71, 0x0c, 0x71, 0x88, 0x6e, 0x27, 0x19, 0x2e, 0xa9, 0xbc, 0xd4, 0x12, 0x64, 0x2f, 0x81, 0xbc,
	0x03, 0xb5, 0x28, 0x89, 0x0d, 0xec, 0x09, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x9e, 0x8a,
	0x8e, 0xf6, 0x62, 0x01, 0x00, 0x00,
}
