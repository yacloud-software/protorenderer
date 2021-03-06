// Code generated by protoc-gen-go.
// source: yacloud.eu/apis/urlmapper/urlmapper.proto
// DO NOT EDIT!

/*
Package urlmapper is a generated protocol buffer package.

It is generated from these files:
	yacloud.eu/apis/urlmapper/urlmapper.proto

It has these top-level messages:
	JsonMapping
	GetJsonMappingRequest
	JsonMappingResponse
	JsonMappingResponseList
	DomainList
	ServiceID
*/
package urlmapper

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import common "golang.conradwood.net/apis/common"

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

type JsonMapping struct {
	ID        uint64 `protobuf:"varint,1,opt,name=ID" json:"ID,omitempty"`
	Domain    string `protobuf:"bytes,2,opt,name=Domain" json:"Domain,omitempty"`
	Path      string `protobuf:"bytes,3,opt,name=Path" json:"Path,omitempty"`
	ServiceID string `protobuf:"bytes,4,opt,name=ServiceID" json:"ServiceID,omitempty"`
	GroupID   string `protobuf:"bytes,5,opt,name=GroupID" json:"GroupID,omitempty"`
}

func (m *JsonMapping) Reset()                    { *m = JsonMapping{} }
func (m *JsonMapping) String() string            { return proto.CompactTextString(m) }
func (*JsonMapping) ProtoMessage()               {}
func (*JsonMapping) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *JsonMapping) GetID() uint64 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *JsonMapping) GetDomain() string {
	if m != nil {
		return m.Domain
	}
	return ""
}

func (m *JsonMapping) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *JsonMapping) GetServiceID() string {
	if m != nil {
		return m.ServiceID
	}
	return ""
}

func (m *JsonMapping) GetGroupID() string {
	if m != nil {
		return m.GroupID
	}
	return ""
}

type GetJsonMappingRequest struct {
	Domain string `protobuf:"bytes,1,opt,name=Domain" json:"Domain,omitempty"`
	Path   string `protobuf:"bytes,2,opt,name=Path" json:"Path,omitempty"`
}

func (m *GetJsonMappingRequest) Reset()                    { *m = GetJsonMappingRequest{} }
func (m *GetJsonMappingRequest) String() string            { return proto.CompactTextString(m) }
func (*GetJsonMappingRequest) ProtoMessage()               {}
func (*GetJsonMappingRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *GetJsonMappingRequest) GetDomain() string {
	if m != nil {
		return m.Domain
	}
	return ""
}

func (m *GetJsonMappingRequest) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

type JsonMappingResponse struct {
	Mapping     *JsonMapping `protobuf:"bytes,1,opt,name=Mapping" json:"Mapping,omitempty"`
	GRPCService string       `protobuf:"bytes,2,opt,name=GRPCService" json:"GRPCService,omitempty"`
}

func (m *JsonMappingResponse) Reset()                    { *m = JsonMappingResponse{} }
func (m *JsonMappingResponse) String() string            { return proto.CompactTextString(m) }
func (*JsonMappingResponse) ProtoMessage()               {}
func (*JsonMappingResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *JsonMappingResponse) GetMapping() *JsonMapping {
	if m != nil {
		return m.Mapping
	}
	return nil
}

func (m *JsonMappingResponse) GetGRPCService() string {
	if m != nil {
		return m.GRPCService
	}
	return ""
}

type JsonMappingResponseList struct {
	Responses []*JsonMappingResponse `protobuf:"bytes,1,rep,name=Responses" json:"Responses,omitempty"`
}

func (m *JsonMappingResponseList) Reset()                    { *m = JsonMappingResponseList{} }
func (m *JsonMappingResponseList) String() string            { return proto.CompactTextString(m) }
func (*JsonMappingResponseList) ProtoMessage()               {}
func (*JsonMappingResponseList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *JsonMappingResponseList) GetResponses() []*JsonMappingResponse {
	if m != nil {
		return m.Responses
	}
	return nil
}

type DomainList struct {
	Domains []string `protobuf:"bytes,1,rep,name=Domains" json:"Domains,omitempty"`
}

func (m *DomainList) Reset()                    { *m = DomainList{} }
func (m *DomainList) String() string            { return proto.CompactTextString(m) }
func (*DomainList) ProtoMessage()               {}
func (*DomainList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *DomainList) GetDomains() []string {
	if m != nil {
		return m.Domains
	}
	return nil
}

type ServiceID struct {
	ID string `protobuf:"bytes,1,opt,name=ID" json:"ID,omitempty"`
}

func (m *ServiceID) Reset()                    { *m = ServiceID{} }
func (m *ServiceID) String() string            { return proto.CompactTextString(m) }
func (*ServiceID) ProtoMessage()               {}
func (*ServiceID) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *ServiceID) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func init() {
	proto.RegisterType((*JsonMapping)(nil), "urlmapper.JsonMapping")
	proto.RegisterType((*GetJsonMappingRequest)(nil), "urlmapper.GetJsonMappingRequest")
	proto.RegisterType((*JsonMappingResponse)(nil), "urlmapper.JsonMappingResponse")
	proto.RegisterType((*JsonMappingResponseList)(nil), "urlmapper.JsonMappingResponseList")
	proto.RegisterType((*DomainList)(nil), "urlmapper.DomainList")
	proto.RegisterType((*ServiceID)(nil), "urlmapper.ServiceID")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for URLMapper service

type URLMapperClient interface {
	// get json mapping for URL and user
	GetJsonMappingWithUser(ctx context.Context, in *GetJsonMappingRequest, opts ...grpc.CallOption) (*JsonMappingResponse, error)
	// get json mapping for URL
	GetJsonMapping(ctx context.Context, in *GetJsonMappingRequest, opts ...grpc.CallOption) (*JsonMappingResponse, error)
	// get all json mappings
	GetJsonMappings(ctx context.Context, in *common.Void, opts ...grpc.CallOption) (*JsonMappingResponseList, error)
	// add a json mapping
	AddJsonMapping(ctx context.Context, in *JsonMapping, opts ...grpc.CallOption) (*common.Void, error)
	// get all domains "mapped" to jsonmultiplexer
	GetJsonDomains(ctx context.Context, in *common.Void, opts ...grpc.CallOption) (*DomainList, error)
	// get mappings for a service
	GetServiceMappings(ctx context.Context, in *ServiceID, opts ...grpc.CallOption) (*JsonMappingResponseList, error)
}

type uRLMapperClient struct {
	cc *grpc.ClientConn
}

func NewURLMapperClient(cc *grpc.ClientConn) URLMapperClient {
	return &uRLMapperClient{cc}
}

func (c *uRLMapperClient) GetJsonMappingWithUser(ctx context.Context, in *GetJsonMappingRequest, opts ...grpc.CallOption) (*JsonMappingResponse, error) {
	out := new(JsonMappingResponse)
	err := grpc.Invoke(ctx, "/urlmapper.URLMapper/GetJsonMappingWithUser", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLMapperClient) GetJsonMapping(ctx context.Context, in *GetJsonMappingRequest, opts ...grpc.CallOption) (*JsonMappingResponse, error) {
	out := new(JsonMappingResponse)
	err := grpc.Invoke(ctx, "/urlmapper.URLMapper/GetJsonMapping", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLMapperClient) GetJsonMappings(ctx context.Context, in *common.Void, opts ...grpc.CallOption) (*JsonMappingResponseList, error) {
	out := new(JsonMappingResponseList)
	err := grpc.Invoke(ctx, "/urlmapper.URLMapper/GetJsonMappings", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLMapperClient) AddJsonMapping(ctx context.Context, in *JsonMapping, opts ...grpc.CallOption) (*common.Void, error) {
	out := new(common.Void)
	err := grpc.Invoke(ctx, "/urlmapper.URLMapper/AddJsonMapping", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLMapperClient) GetJsonDomains(ctx context.Context, in *common.Void, opts ...grpc.CallOption) (*DomainList, error) {
	out := new(DomainList)
	err := grpc.Invoke(ctx, "/urlmapper.URLMapper/GetJsonDomains", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLMapperClient) GetServiceMappings(ctx context.Context, in *ServiceID, opts ...grpc.CallOption) (*JsonMappingResponseList, error) {
	out := new(JsonMappingResponseList)
	err := grpc.Invoke(ctx, "/urlmapper.URLMapper/GetServiceMappings", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for URLMapper service

type URLMapperServer interface {
	// get json mapping for URL and user
	GetJsonMappingWithUser(context.Context, *GetJsonMappingRequest) (*JsonMappingResponse, error)
	// get json mapping for URL
	GetJsonMapping(context.Context, *GetJsonMappingRequest) (*JsonMappingResponse, error)
	// get all json mappings
	GetJsonMappings(context.Context, *common.Void) (*JsonMappingResponseList, error)
	// add a json mapping
	AddJsonMapping(context.Context, *JsonMapping) (*common.Void, error)
	// get all domains "mapped" to jsonmultiplexer
	GetJsonDomains(context.Context, *common.Void) (*DomainList, error)
	// get mappings for a service
	GetServiceMappings(context.Context, *ServiceID) (*JsonMappingResponseList, error)
}

func RegisterURLMapperServer(s *grpc.Server, srv URLMapperServer) {
	s.RegisterService(&_URLMapper_serviceDesc, srv)
}

func _URLMapper_GetJsonMappingWithUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetJsonMappingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLMapperServer).GetJsonMappingWithUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/urlmapper.URLMapper/GetJsonMappingWithUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLMapperServer).GetJsonMappingWithUser(ctx, req.(*GetJsonMappingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLMapper_GetJsonMapping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetJsonMappingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLMapperServer).GetJsonMapping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/urlmapper.URLMapper/GetJsonMapping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLMapperServer).GetJsonMapping(ctx, req.(*GetJsonMappingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLMapper_GetJsonMappings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(common.Void)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLMapperServer).GetJsonMappings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/urlmapper.URLMapper/GetJsonMappings",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLMapperServer).GetJsonMappings(ctx, req.(*common.Void))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLMapper_AddJsonMapping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JsonMapping)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLMapperServer).AddJsonMapping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/urlmapper.URLMapper/AddJsonMapping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLMapperServer).AddJsonMapping(ctx, req.(*JsonMapping))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLMapper_GetJsonDomains_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(common.Void)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLMapperServer).GetJsonDomains(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/urlmapper.URLMapper/GetJsonDomains",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLMapperServer).GetJsonDomains(ctx, req.(*common.Void))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLMapper_GetServiceMappings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ServiceID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLMapperServer).GetServiceMappings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/urlmapper.URLMapper/GetServiceMappings",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLMapperServer).GetServiceMappings(ctx, req.(*ServiceID))
	}
	return interceptor(ctx, in, info, handler)
}

var _URLMapper_serviceDesc = grpc.ServiceDesc{
	ServiceName: "urlmapper.URLMapper",
	HandlerType: (*URLMapperServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetJsonMappingWithUser",
			Handler:    _URLMapper_GetJsonMappingWithUser_Handler,
		},
		{
			MethodName: "GetJsonMapping",
			Handler:    _URLMapper_GetJsonMapping_Handler,
		},
		{
			MethodName: "GetJsonMappings",
			Handler:    _URLMapper_GetJsonMappings_Handler,
		},
		{
			MethodName: "AddJsonMapping",
			Handler:    _URLMapper_AddJsonMapping_Handler,
		},
		{
			MethodName: "GetJsonDomains",
			Handler:    _URLMapper_GetJsonDomains_Handler,
		},
		{
			MethodName: "GetServiceMappings",
			Handler:    _URLMapper_GetServiceMappings_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "yacloud.eu/apis/urlmapper/urlmapper.proto",
}

func init() { proto.RegisterFile("yacloud.eu/apis/urlmapper/urlmapper.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 441 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xac, 0x53, 0x41, 0x6f, 0xd3, 0x30,
	0x14, 0x56, 0xda, 0xb2, 0x2a, 0xaf, 0xa8, 0x48, 0x0f, 0x56, 0x42, 0x87, 0x50, 0x94, 0x03, 0x2a,
	0x17, 0x0f, 0x15, 0xc1, 0x89, 0x0b, 0x34, 0x52, 0x55, 0xb4, 0xa1, 0xc9, 0x68, 0x4c, 0xe2, 0x66,
	0x1a, 0xab, 0xb3, 0xd4, 0xda, 0x26, 0x76, 0x40, 0x1c, 0xf9, 0x65, 0xfc, 0x35, 0x34, 0xc7, 0x69,
	0x1c, 0xd4, 0x01, 0x07, 0x4e, 0xf1, 0x7b, 0xf6, 0xfb, 0xbe, 0xcf, 0xdf, 0x17, 0xc3, 0xb3, 0xef,
	0x6c, 0xbd, 0x55, 0x55, 0x41, 0x78, 0x75, 0xca, 0xb4, 0x30, 0xa7, 0x55, 0xb9, 0xdd, 0x31, 0xad,
	0x79, 0xd9, 0xae, 0x88, 0x2e, 0x95, 0x55, 0x18, 0xef, 0x1b, 0x53, 0xb2, 0x51, 0x5b, 0x26, 0x37,
	0x64, 0xad, 0x64, 0xc9, 0x8a, 0x6f, 0x4a, 0x15, 0x44, 0x72, 0x5b, 0x03, 0xac, 0xd5, 0x6e, 0xa7,
	0xa4, 0xff, 0xd4, 0xa3, 0xd9, 0x8f, 0x08, 0x46, 0xef, 0x8c, 0x92, 0xe7, 0x4c, 0x6b, 0x21, 0x37,
	0x38, 0x86, 0xde, 0x2a, 0x4f, 0xa2, 0x34, 0x9a, 0x0d, 0x68, 0x6f, 0x95, 0xe3, 0x04, 0x8e, 0x72,
	0xb5, 0x63, 0x42, 0x26, 0xbd, 0x34, 0x9a, 0xc5, 0xd4, 0x57, 0x88, 0x30, 0xb8, 0x60, 0xf6, 0x3a,
	0xe9, 0xbb, 0xae, 0x5b, 0xe3, 0x63, 0x88, 0x3f, 0xf0, 0xf2, 0xab, 0x58, 0xf3, 0x55, 0x9e, 0x0c,
	0xdc, 0x46, 0xdb, 0xc0, 0x04, 0x86, 0xcb, 0x52, 0x55, 0x7a, 0x95, 0x27, 0x77, 0xdc, 0x5e, 0x53,
	0x66, 0x0b, 0x38, 0x5e, 0x72, 0x1b, 0xa8, 0xa0, 0xfc, 0x4b, 0xc5, 0x8d, 0x0d, 0xc8, 0xa3, 0x83,
	0xe4, 0xbd, 0x96, 0x3c, 0x13, 0x70, 0xbf, 0x83, 0x60, 0xb4, 0x92, 0x86, 0xe3, 0x73, 0x18, 0xfa,
	0x96, 0xc3, 0x18, 0xcd, 0x27, 0xa4, 0x75, 0x2f, 0x1c, 0x68, 0x8e, 0x61, 0x0a, 0xa3, 0x25, 0xbd,
	0x58, 0x78, 0xe1, 0x9e, 0x23, 0x6c, 0x65, 0x57, 0xf0, 0xf0, 0x00, 0xd5, 0x99, 0x30, 0x16, 0x5f,
	0x43, 0xdc, 0xd4, 0x26, 0x89, 0xd2, 0xfe, 0x6c, 0x34, 0x7f, 0x72, 0x0b, 0xa1, 0x3f, 0x46, 0xdb,
	0x81, 0xec, 0x29, 0x40, 0x7d, 0x43, 0x87, 0x95, 0xc0, 0xb0, 0xae, 0x6a, 0xa4, 0x98, 0x36, 0x65,
	0x76, 0x12, 0x18, 0x1d, 0x24, 0x16, 0xdf, 0x24, 0x36, 0xff, 0xd9, 0x87, 0xf8, 0x92, 0x9e, 0x9d,
	0x3b, 0x46, 0xfc, 0x04, 0x93, 0xae, 0xb7, 0x57, 0xc2, 0x5e, 0x5f, 0x1a, 0x5e, 0x62, 0x1a, 0xe8,
	0x3a, 0x68, 0xff, 0xf4, 0x2f, 0xca, 0x91, 0xc2, 0xb8, 0x3b, 0xf8, 0x1f, 0x30, 0x17, 0x70, 0xaf,
	0x3b, 0x68, 0xf0, 0x2e, 0xf1, 0x7f, 0xec, 0x47, 0x25, 0x8a, 0x69, 0xf6, 0x67, 0x00, 0xe7, 0xdc,
	0x2b, 0x18, 0xbf, 0x29, 0x8a, 0x50, 0xd8, 0x2d, 0xa9, 0x4f, 0x3b, 0xd8, 0xf8, 0x72, 0x7f, 0x21,
	0xef, 0xf4, 0x6f, 0xdc, 0xc7, 0x01, 0x4a, 0x10, 0xd4, 0x7b, 0xc0, 0x25, 0xb7, 0x3e, 0x91, 0xbd,
	0xec, 0x07, 0xc1, 0xe1, 0x7d, 0x5a, 0xff, 0x22, 0xff, 0xed, 0x09, 0x3c, 0xe2, 0x15, 0x69, 0x9e,
	0xff, 0xcd, 0xd3, 0x6d, 0x87, 0x3e, 0x1f, 0xb9, 0x77, 0xfb, 0xe2, 0x57, 0x00, 0x00, 0x00, 0xff,
	0xff, 0xce, 0x8b, 0x99, 0xd9, 0x1f, 0x04, 0x00, 0x00,
}
