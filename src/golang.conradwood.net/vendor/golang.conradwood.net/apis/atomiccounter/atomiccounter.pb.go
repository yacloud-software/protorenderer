// Code generated by protoc-gen-go.
// source: golang.conradwood.net/apis/atomiccounter/atomiccounter.proto
// DO NOT EDIT!

/*
Package atomiccounter is a generated protocol buffer package.

It is generated from these files:
	golang.conradwood.net/apis/atomiccounter/atomiccounter.proto

It has these top-level messages:
	Counter
	ModifyRequest
	ModifyResponse
	CompareAndModifyRequest
*/
package atomiccounter

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

type Counter struct {
	CounterName string `protobuf:"bytes,1,opt,name=CounterName" json:"CounterName,omitempty"`
}

func (m *Counter) Reset()                    { *m = Counter{} }
func (m *Counter) String() string            { return proto.CompactTextString(m) }
func (*Counter) ProtoMessage()               {}
func (*Counter) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Counter) GetCounterName() string {
	if m != nil {
		return m.CounterName
	}
	return ""
}

type ModifyRequest struct {
	Counter *Counter `protobuf:"bytes,1,opt,name=Counter" json:"Counter,omitempty"`
	Value   int64    `protobuf:"varint,2,opt,name=Value" json:"Value,omitempty"`
}

func (m *ModifyRequest) Reset()                    { *m = ModifyRequest{} }
func (m *ModifyRequest) String() string            { return proto.CompactTextString(m) }
func (*ModifyRequest) ProtoMessage()               {}
func (*ModifyRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ModifyRequest) GetCounter() *Counter {
	if m != nil {
		return m.Counter
	}
	return nil
}

func (m *ModifyRequest) GetValue() int64 {
	if m != nil {
		return m.Value
	}
	return 0
}

type ModifyResponse struct {
	OldValue int64 `protobuf:"varint,1,opt,name=OldValue" json:"OldValue,omitempty"`
	NewValue int64 `protobuf:"varint,2,opt,name=NewValue" json:"NewValue,omitempty"`
}

func (m *ModifyResponse) Reset()                    { *m = ModifyResponse{} }
func (m *ModifyResponse) String() string            { return proto.CompactTextString(m) }
func (*ModifyResponse) ProtoMessage()               {}
func (*ModifyResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *ModifyResponse) GetOldValue() int64 {
	if m != nil {
		return m.OldValue
	}
	return 0
}

func (m *ModifyResponse) GetNewValue() int64 {
	if m != nil {
		return m.NewValue
	}
	return 0
}

type CompareAndModifyRequest struct {
	CompareValue int64    `protobuf:"varint,1,opt,name=CompareValue" json:"CompareValue,omitempty"`
	Counter      *Counter `protobuf:"bytes,2,opt,name=Counter" json:"Counter,omitempty"`
	Value        int64    `protobuf:"varint,3,opt,name=Value" json:"Value,omitempty"`
}

func (m *CompareAndModifyRequest) Reset()                    { *m = CompareAndModifyRequest{} }
func (m *CompareAndModifyRequest) String() string            { return proto.CompactTextString(m) }
func (*CompareAndModifyRequest) ProtoMessage()               {}
func (*CompareAndModifyRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *CompareAndModifyRequest) GetCompareValue() int64 {
	if m != nil {
		return m.CompareValue
	}
	return 0
}

func (m *CompareAndModifyRequest) GetCounter() *Counter {
	if m != nil {
		return m.Counter
	}
	return nil
}

func (m *CompareAndModifyRequest) GetValue() int64 {
	if m != nil {
		return m.Value
	}
	return 0
}

func init() {
	proto.RegisterType((*Counter)(nil), "atomiccounter.Counter")
	proto.RegisterType((*ModifyRequest)(nil), "atomiccounter.ModifyRequest")
	proto.RegisterType((*ModifyResponse)(nil), "atomiccounter.ModifyResponse")
	proto.RegisterType((*CompareAndModifyRequest)(nil), "atomiccounter.CompareAndModifyRequest")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for AtomicCounterService service

type AtomicCounterServiceClient interface {
	// increment a counter
	Inc(ctx context.Context, in *Counter, opts ...grpc.CallOption) (*ModifyResponse, error)
	// decrement a counter
	Dec(ctx context.Context, in *Counter, opts ...grpc.CallOption) (*ModifyResponse, error)
	// modify a counter (set to an absolute value)
	Modify(ctx context.Context, in *ModifyRequest, opts ...grpc.CallOption) (*ModifyResponse, error)
	// the classic Read-And-Modify
	// it is intented for atomic singular
	// increment of COUNTERS
	// DO NOT USE THIS TO BUILD GLOBAL LOCKS!
	// GLOBAL LOCKS are catastrophic for scaleability
	// and availability.
	// if you think you need to create a global lock
	// with this function, you *must* consult the CTO prior to doing so.
	// it's a mandatory process!
	CompareAndInc(ctx context.Context, in *CompareAndModifyRequest, opts ...grpc.CallOption) (*ModifyResponse, error)
	// same as above but instead of adding the value it'll set it
	CompareAndSet(ctx context.Context, in *CompareAndModifyRequest, opts ...grpc.CallOption) (*ModifyResponse, error)
}

type atomicCounterServiceClient struct {
	cc *grpc.ClientConn
}

func NewAtomicCounterServiceClient(cc *grpc.ClientConn) AtomicCounterServiceClient {
	return &atomicCounterServiceClient{cc}
}

func (c *atomicCounterServiceClient) Inc(ctx context.Context, in *Counter, opts ...grpc.CallOption) (*ModifyResponse, error) {
	out := new(ModifyResponse)
	err := grpc.Invoke(ctx, "/atomiccounter.AtomicCounterService/Inc", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *atomicCounterServiceClient) Dec(ctx context.Context, in *Counter, opts ...grpc.CallOption) (*ModifyResponse, error) {
	out := new(ModifyResponse)
	err := grpc.Invoke(ctx, "/atomiccounter.AtomicCounterService/Dec", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *atomicCounterServiceClient) Modify(ctx context.Context, in *ModifyRequest, opts ...grpc.CallOption) (*ModifyResponse, error) {
	out := new(ModifyResponse)
	err := grpc.Invoke(ctx, "/atomiccounter.AtomicCounterService/Modify", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *atomicCounterServiceClient) CompareAndInc(ctx context.Context, in *CompareAndModifyRequest, opts ...grpc.CallOption) (*ModifyResponse, error) {
	out := new(ModifyResponse)
	err := grpc.Invoke(ctx, "/atomiccounter.AtomicCounterService/CompareAndInc", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *atomicCounterServiceClient) CompareAndSet(ctx context.Context, in *CompareAndModifyRequest, opts ...grpc.CallOption) (*ModifyResponse, error) {
	out := new(ModifyResponse)
	err := grpc.Invoke(ctx, "/atomiccounter.AtomicCounterService/CompareAndSet", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for AtomicCounterService service

type AtomicCounterServiceServer interface {
	// increment a counter
	Inc(context.Context, *Counter) (*ModifyResponse, error)
	// decrement a counter
	Dec(context.Context, *Counter) (*ModifyResponse, error)
	// modify a counter (set to an absolute value)
	Modify(context.Context, *ModifyRequest) (*ModifyResponse, error)
	// the classic Read-And-Modify
	// it is intented for atomic singular
	// increment of COUNTERS
	// DO NOT USE THIS TO BUILD GLOBAL LOCKS!
	// GLOBAL LOCKS are catastrophic for scaleability
	// and availability.
	// if you think you need to create a global lock
	// with this function, you *must* consult the CTO prior to doing so.
	// it's a mandatory process!
	CompareAndInc(context.Context, *CompareAndModifyRequest) (*ModifyResponse, error)
	// same as above but instead of adding the value it'll set it
	CompareAndSet(context.Context, *CompareAndModifyRequest) (*ModifyResponse, error)
}

func RegisterAtomicCounterServiceServer(s *grpc.Server, srv AtomicCounterServiceServer) {
	s.RegisterService(&_AtomicCounterService_serviceDesc, srv)
}

func _AtomicCounterService_Inc_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Counter)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AtomicCounterServiceServer).Inc(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/atomiccounter.AtomicCounterService/Inc",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AtomicCounterServiceServer).Inc(ctx, req.(*Counter))
	}
	return interceptor(ctx, in, info, handler)
}

func _AtomicCounterService_Dec_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Counter)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AtomicCounterServiceServer).Dec(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/atomiccounter.AtomicCounterService/Dec",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AtomicCounterServiceServer).Dec(ctx, req.(*Counter))
	}
	return interceptor(ctx, in, info, handler)
}

func _AtomicCounterService_Modify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ModifyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AtomicCounterServiceServer).Modify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/atomiccounter.AtomicCounterService/Modify",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AtomicCounterServiceServer).Modify(ctx, req.(*ModifyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AtomicCounterService_CompareAndInc_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CompareAndModifyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AtomicCounterServiceServer).CompareAndInc(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/atomiccounter.AtomicCounterService/CompareAndInc",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AtomicCounterServiceServer).CompareAndInc(ctx, req.(*CompareAndModifyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AtomicCounterService_CompareAndSet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CompareAndModifyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AtomicCounterServiceServer).CompareAndSet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/atomiccounter.AtomicCounterService/CompareAndSet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AtomicCounterServiceServer).CompareAndSet(ctx, req.(*CompareAndModifyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _AtomicCounterService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "atomiccounter.AtomicCounterService",
	HandlerType: (*AtomicCounterServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Inc",
			Handler:    _AtomicCounterService_Inc_Handler,
		},
		{
			MethodName: "Dec",
			Handler:    _AtomicCounterService_Dec_Handler,
		},
		{
			MethodName: "Modify",
			Handler:    _AtomicCounterService_Modify_Handler,
		},
		{
			MethodName: "CompareAndInc",
			Handler:    _AtomicCounterService_CompareAndInc_Handler,
		},
		{
			MethodName: "CompareAndSet",
			Handler:    _AtomicCounterService_CompareAndSet_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "golang.conradwood.net/apis/atomiccounter/atomiccounter.proto",
}

func init() {
	proto.RegisterFile("golang.conradwood.net/apis/atomiccounter/atomiccounter.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 321 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xac, 0x53, 0x4d, 0x4b, 0xc3, 0x40,
	0x10, 0x25, 0x0d, 0x56, 0x9d, 0x5a, 0x0f, 0x4b, 0xd1, 0x52, 0x14, 0x4a, 0x0e, 0x52, 0x14, 0x52,
	0xa9, 0xd7, 0x5c, 0x6a, 0x15, 0xf4, 0x60, 0x85, 0x14, 0xea, 0x79, 0x4d, 0xc6, 0x12, 0x48, 0x76,
	0x62, 0xb2, 0xb1, 0xf8, 0x03, 0x3c, 0xf8, 0xaf, 0x25, 0x1f, 0xc6, 0xdd, 0x6a, 0x91, 0x8a, 0xb7,
	0x9d, 0x37, 0x6f, 0xde, 0xee, 0x9b, 0xe5, 0x81, 0xb3, 0xa0, 0x90, 0x8b, 0x85, 0xed, 0x91, 0x48,
	0xb8, 0xbf, 0x24, 0xf2, 0x6d, 0x81, 0x72, 0xc8, 0xe3, 0x20, 0x1d, 0x72, 0x49, 0x51, 0xe0, 0x79,
	0x94, 0x09, 0x89, 0x89, 0x5e, 0xd9, 0x71, 0x42, 0x92, 0x58, 0x5b, 0x03, 0xad, 0x33, 0xd8, 0x9e,
	0x94, 0x47, 0xd6, 0x87, 0x56, 0x75, 0x9c, 0xf2, 0x08, 0xbb, 0x46, 0xdf, 0x18, 0xec, 0xba, 0x2a,
	0x64, 0x3d, 0x40, 0xfb, 0x8e, 0xfc, 0xe0, 0xe9, 0xd5, 0xc5, 0xe7, 0x0c, 0x53, 0xc9, 0xce, 0xeb,
	0xe9, 0x82, 0xde, 0x1a, 0x1d, 0xd8, 0xfa, 0x9d, 0x55, 0xd7, 0xad, 0x2f, 0xe9, 0xc0, 0xd6, 0x9c,
	0x87, 0x19, 0x76, 0x1b, 0x7d, 0x63, 0x60, 0xba, 0x65, 0x61, 0xdd, 0xc0, 0xfe, 0xa7, 0x70, 0x1a,
	0x93, 0x48, 0x91, 0xf5, 0x60, 0xe7, 0x3e, 0xf4, 0x4b, 0xaa, 0x51, 0x50, 0xeb, 0x3a, 0xef, 0x4d,
	0x71, 0xa9, 0xca, 0xd4, 0xb5, 0xf5, 0x66, 0xc0, 0xe1, 0x84, 0xa2, 0x98, 0x27, 0x38, 0x16, 0xbe,
	0xfe, 0x5a, 0x0b, 0xf6, 0xaa, 0x96, 0xaa, 0xab, 0x61, 0xaa, 0xa3, 0xc6, 0x86, 0x8e, 0x4c, 0xc5,
	0xd1, 0xe8, 0xdd, 0x84, 0xce, 0xb8, 0x18, 0xac, 0x78, 0x33, 0x4c, 0x5e, 0x02, 0x0f, 0x99, 0x03,
	0xe6, 0xad, 0xf0, 0xd8, 0x1a, 0xd9, 0xde, 0xf1, 0x0a, 0xbe, 0xb2, 0x16, 0x07, 0xcc, 0x2b, 0xfc,
	0xf3, 0xf4, 0x35, 0x34, 0x4b, 0x84, 0x1d, 0xad, 0x21, 0x16, 0x8b, 0xfa, 0x4d, 0x66, 0x0e, 0xed,
	0xaf, 0x15, 0xe7, 0x66, 0x4e, 0xbe, 0x3d, 0xe7, 0xc7, 0x0f, 0xd8, 0x48, 0x77, 0x86, 0xf2, 0x9f,
	0x74, 0x2f, 0x4f, 0x61, 0x20, 0x50, 0xaa, 0x79, 0xa9, 0x12, 0x94, 0x47, 0x46, 0x1f, 0x7d, 0x6c,
	0x16, 0x29, 0xb9, 0xf8, 0x08, 0x00, 0x00, 0xff, 0xff, 0x4b, 0x8f, 0x00, 0xb7, 0x65, 0x03, 0x00,
	0x00,
}
