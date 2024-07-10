// Code generated by protoc-gen-go.
// source: protos/golang.conradwood.net/apis/protorenderer2/protorenderer2.proto
// DO NOT EDIT!

/*
Package protorenderer2 is a generated protocol buffer package.

It is generated from these files:
	protos/golang.conradwood.net/apis/protorenderer2/protorenderer2.proto

It has these top-level messages:
	ProtoFileInfo
	ProtoMessageInfo
	ProtoFieldInfo
	DBProtoFile
	SQLMessage
	ProtocRequest
	FileTransfer
	FileResult
	CompileFailure
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

type FieldModifier int32

const (
	FieldModifier_FM_UNDEFINED FieldModifier = 0
	FieldModifier_SINGLE       FieldModifier = 1
	FieldModifier_MAP          FieldModifier = 2
	FieldModifier_ARRAY        FieldModifier = 3
)

var FieldModifier_name = map[int32]string{
	0: "FM_UNDEFINED",
	1: "SINGLE",
	2: "MAP",
	3: "ARRAY",
}
var FieldModifier_value = map[string]int32{
	"FM_UNDEFINED": 0,
	"SINGLE":       1,
	"MAP":          2,
	"ARRAY":        3,
}

func (x FieldModifier) String() string {
	return proto.EnumName(FieldModifier_name, int32(x))
}
func (FieldModifier) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type FieldType int32

const (
	FieldType_ft_UNDEFINED FieldType = 0
	FieldType_PRIMITIVE    FieldType = 1
	FieldType_OBJECT       FieldType = 2
)

var FieldType_name = map[int32]string{
	0: "ft_UNDEFINED",
	1: "PRIMITIVE",
	2: "OBJECT",
}
var FieldType_value = map[string]int32{
	"ft_UNDEFINED": 0,
	"PRIMITIVE":    1,
	"OBJECT":       2,
}

func (x FieldType) String() string {
	return proto.EnumName(FieldType_name, int32(x))
}
func (FieldType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type ProtoFieldPrimitive int32

const (
	ProtoFieldPrimitive_PFT_UNDEFINED ProtoFieldPrimitive = 0
	ProtoFieldPrimitive_STRING        ProtoFieldPrimitive = 1
	ProtoFieldPrimitive_UINT64        ProtoFieldPrimitive = 2
	ProtoFieldPrimitive_UINT32        ProtoFieldPrimitive = 3
	ProtoFieldPrimitive_BYTES         ProtoFieldPrimitive = 4
	ProtoFieldPrimitive_BOOL          ProtoFieldPrimitive = 5
	ProtoFieldPrimitive_INT32         ProtoFieldPrimitive = 6
	ProtoFieldPrimitive_ENUM          ProtoFieldPrimitive = 7
	ProtoFieldPrimitive_DOUBLE        ProtoFieldPrimitive = 8
	ProtoFieldPrimitive_FLOAT         ProtoFieldPrimitive = 9
	ProtoFieldPrimitive_INT64         ProtoFieldPrimitive = 10
)

var ProtoFieldPrimitive_name = map[int32]string{
	0:  "PFT_UNDEFINED",
	1:  "STRING",
	2:  "UINT64",
	3:  "UINT32",
	4:  "BYTES",
	5:  "BOOL",
	6:  "INT32",
	7:  "ENUM",
	8:  "DOUBLE",
	9:  "FLOAT",
	10: "INT64",
}
var ProtoFieldPrimitive_value = map[string]int32{
	"PFT_UNDEFINED": 0,
	"STRING":        1,
	"UINT64":        2,
	"UINT32":        3,
	"BYTES":         4,
	"BOOL":          5,
	"INT32":         6,
	"ENUM":          7,
	"DOUBLE":        8,
	"FLOAT":         9,
	"INT64":         10,
}

func (x ProtoFieldPrimitive) String() string {
	return proto.EnumName(ProtoFieldPrimitive_name, int32(x))
}
func (ProtoFieldPrimitive) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

// information parsed out of a .proto file, enriched with references to database IDs
type ProtoFileInfo struct {
	ProtoFile  *DBProtoFile        `protobuf:"bytes,1,opt,name=ProtoFile" json:"ProtoFile,omitempty"`
	Imports    []*DBProtoFile      `protobuf:"bytes,2,rep,name=Imports" json:"Imports,omitempty"`
	CNWOptions map[string]string   `protobuf:"bytes,3,rep,name=CNWOptions" json:"CNWOptions,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Messages   []*ProtoMessageInfo `protobuf:"bytes,4,rep,name=Messages" json:"Messages,omitempty"`
}

func (m *ProtoFileInfo) Reset()                    { *m = ProtoFileInfo{} }
func (m *ProtoFileInfo) String() string            { return proto.CompactTextString(m) }
func (*ProtoFileInfo) ProtoMessage()               {}
func (*ProtoFileInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ProtoFileInfo) GetProtoFile() *DBProtoFile {
	if m != nil {
		return m.ProtoFile
	}
	return nil
}

func (m *ProtoFileInfo) GetImports() []*DBProtoFile {
	if m != nil {
		return m.Imports
	}
	return nil
}

func (m *ProtoFileInfo) GetCNWOptions() map[string]string {
	if m != nil {
		return m.CNWOptions
	}
	return nil
}

func (m *ProtoFileInfo) GetMessages() []*ProtoMessageInfo {
	if m != nil {
		return m.Messages
	}
	return nil
}

type ProtoMessageInfo struct {
	Message *SQLMessage       `protobuf:"bytes,1,opt,name=Message" json:"Message,omitempty"`
	Comment string            `protobuf:"bytes,2,opt,name=Comment" json:"Comment,omitempty"`
	Fields  []*ProtoFieldInfo `protobuf:"bytes,3,rep,name=Fields" json:"Fields,omitempty"`
}

func (m *ProtoMessageInfo) Reset()                    { *m = ProtoMessageInfo{} }
func (m *ProtoMessageInfo) String() string            { return proto.CompactTextString(m) }
func (*ProtoMessageInfo) ProtoMessage()               {}
func (*ProtoMessageInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ProtoMessageInfo) GetMessage() *SQLMessage {
	if m != nil {
		return m.Message
	}
	return nil
}

func (m *ProtoMessageInfo) GetComment() string {
	if m != nil {
		return m.Comment
	}
	return ""
}

func (m *ProtoMessageInfo) GetFields() []*ProtoFieldInfo {
	if m != nil {
		return m.Fields
	}
	return nil
}

type ProtoFieldInfo struct {
	Name           string              `protobuf:"bytes,1,opt,name=Name" json:"Name,omitempty"`
	Comment        string              `protobuf:"bytes,2,opt,name=Comment" json:"Comment,omitempty"`
	Modifier       FieldModifier       `protobuf:"varint,3,opt,name=Modifier,enum=protorenderer2.FieldModifier" json:"Modifier,omitempty"`
	Type1          FieldType           `protobuf:"varint,4,opt,name=Type1,enum=protorenderer2.FieldType" json:"Type1,omitempty"`
	Type2          FieldType           `protobuf:"varint,5,opt,name=Type2,enum=protorenderer2.FieldType" json:"Type2,omitempty"`
	PrimitiveType1 ProtoFieldPrimitive `protobuf:"varint,6,opt,name=PrimitiveType1,enum=protorenderer2.ProtoFieldPrimitive" json:"PrimitiveType1,omitempty"`
	PrimitiveType2 ProtoFieldPrimitive `protobuf:"varint,7,opt,name=PrimitiveType2,enum=protorenderer2.ProtoFieldPrimitive" json:"PrimitiveType2,omitempty"`
	ObjectID1      uint64              `protobuf:"varint,8,opt,name=ObjectID1" json:"ObjectID1,omitempty"`
	ObjectID2      uint64              `protobuf:"varint,9,opt,name=ObjectID2" json:"ObjectID2,omitempty"`
}

func (m *ProtoFieldInfo) Reset()                    { *m = ProtoFieldInfo{} }
func (m *ProtoFieldInfo) String() string            { return proto.CompactTextString(m) }
func (*ProtoFieldInfo) ProtoMessage()               {}
func (*ProtoFieldInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *ProtoFieldInfo) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ProtoFieldInfo) GetComment() string {
	if m != nil {
		return m.Comment
	}
	return ""
}

func (m *ProtoFieldInfo) GetModifier() FieldModifier {
	if m != nil {
		return m.Modifier
	}
	return FieldModifier_FM_UNDEFINED
}

func (m *ProtoFieldInfo) GetType1() FieldType {
	if m != nil {
		return m.Type1
	}
	return FieldType_ft_UNDEFINED
}

func (m *ProtoFieldInfo) GetType2() FieldType {
	if m != nil {
		return m.Type2
	}
	return FieldType_ft_UNDEFINED
}

func (m *ProtoFieldInfo) GetPrimitiveType1() ProtoFieldPrimitive {
	if m != nil {
		return m.PrimitiveType1
	}
	return ProtoFieldPrimitive_PFT_UNDEFINED
}

func (m *ProtoFieldInfo) GetPrimitiveType2() ProtoFieldPrimitive {
	if m != nil {
		return m.PrimitiveType2
	}
	return ProtoFieldPrimitive_PFT_UNDEFINED
}

func (m *ProtoFieldInfo) GetObjectID1() uint64 {
	if m != nil {
		return m.ObjectID1
	}
	return 0
}

func (m *ProtoFieldInfo) GetObjectID2() uint64 {
	if m != nil {
		return m.ObjectID2
	}
	return 0
}

// keeping a file in database with metadata
type DBProtoFile struct {
	ID           uint64 `protobuf:"varint,1,opt,name=ID" json:"ID,omitempty"`
	Name         string `protobuf:"bytes,2,opt,name=Name" json:"Name,omitempty"`
	RepositoryID uint64 `protobuf:"varint,3,opt,name=RepositoryID" json:"RepositoryID,omitempty"`
	Package      string `protobuf:"bytes,4,opt,name=Package" json:"Package,omitempty"`
}

func (m *DBProtoFile) Reset()                    { *m = DBProtoFile{} }
func (m *DBProtoFile) String() string            { return proto.CompactTextString(m) }
func (*DBProtoFile) ProtoMessage()               {}
func (*DBProtoFile) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *DBProtoFile) GetID() uint64 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *DBProtoFile) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *DBProtoFile) GetRepositoryID() uint64 {
	if m != nil {
		return m.RepositoryID
	}
	return 0
}

func (m *DBProtoFile) GetPackage() string {
	if m != nil {
		return m.Package
	}
	return ""
}

type SQLMessage struct {
	ID        uint64       `protobuf:"varint,1,opt,name=ID" json:"ID,omitempty"`
	ProtoFile *DBProtoFile `protobuf:"bytes,2,opt,name=ProtoFile" json:"ProtoFile,omitempty"`
	Name      string       `protobuf:"bytes,3,opt,name=Name" json:"Name,omitempty"`
}

func (m *SQLMessage) Reset()                    { *m = SQLMessage{} }
func (m *SQLMessage) String() string            { return proto.CompactTextString(m) }
func (*SQLMessage) ProtoMessage()               {}
func (*SQLMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *SQLMessage) GetID() uint64 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *SQLMessage) GetProtoFile() *DBProtoFile {
	if m != nil {
		return m.ProtoFile
	}
	return nil
}

func (m *SQLMessage) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type ProtocRequest struct {
	MetaCompilerID string                                 `protobuf:"bytes,1,opt,name=MetaCompilerID" json:"MetaCompilerID,omitempty"`
	ProtoFiles     []*google_protobuf.FileDescriptorProto `protobuf:"bytes,2,rep,name=ProtoFiles" json:"ProtoFiles,omitempty"`
}

func (m *ProtocRequest) Reset()                    { *m = ProtocRequest{} }
func (m *ProtocRequest) String() string            { return proto.CompactTextString(m) }
func (*ProtocRequest) ProtoMessage()               {}
func (*ProtocRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *ProtocRequest) GetMetaCompilerID() string {
	if m != nil {
		return m.MetaCompilerID
	}
	return ""
}

func (m *ProtocRequest) GetProtoFiles() []*google_protobuf.FileDescriptorProto {
	if m != nil {
		return m.ProtoFiles
	}
	return nil
}

type FileTransfer struct {
	Filename         string      `protobuf:"bytes,1,opt,name=Filename" json:"Filename,omitempty"`
	Data             []byte      `protobuf:"bytes,2,opt,name=Data,proto3" json:"Data,omitempty"`
	TransferComplete bool        `protobuf:"varint,3,opt,name=TransferComplete" json:"TransferComplete,omitempty"`
	RepositoryID     uint32      `protobuf:"varint,4,opt,name=RepositoryID" json:"RepositoryID,omitempty"`
	Result           *FileResult `protobuf:"bytes,5,opt,name=Result" json:"Result,omitempty"`
}

func (m *FileTransfer) Reset()                    { *m = FileTransfer{} }
func (m *FileTransfer) String() string            { return proto.CompactTextString(m) }
func (*FileTransfer) ProtoMessage()               {}
func (*FileTransfer) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *FileTransfer) GetFilename() string {
	if m != nil {
		return m.Filename
	}
	return ""
}

func (m *FileTransfer) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *FileTransfer) GetTransferComplete() bool {
	if m != nil {
		return m.TransferComplete
	}
	return false
}

func (m *FileTransfer) GetRepositoryID() uint32 {
	if m != nil {
		return m.RepositoryID
	}
	return 0
}

func (m *FileTransfer) GetResult() *FileResult {
	if m != nil {
		return m.Result
	}
	return nil
}

type FileResult struct {
	Failed   bool              `protobuf:"varint,1,opt,name=Failed" json:"Failed,omitempty"`
	Filename string            `protobuf:"bytes,2,opt,name=Filename" json:"Filename,omitempty"`
	Failures []*CompileFailure `protobuf:"bytes,3,rep,name=Failures" json:"Failures,omitempty"`
}

func (m *FileResult) Reset()                    { *m = FileResult{} }
func (m *FileResult) String() string            { return proto.CompactTextString(m) }
func (*FileResult) ProtoMessage()               {}
func (*FileResult) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *FileResult) GetFailed() bool {
	if m != nil {
		return m.Failed
	}
	return false
}

func (m *FileResult) GetFilename() string {
	if m != nil {
		return m.Filename
	}
	return ""
}

func (m *FileResult) GetFailures() []*CompileFailure {
	if m != nil {
		return m.Failures
	}
	return nil
}

type CompileFailure struct {
	CompilerName string `protobuf:"bytes,1,opt,name=CompilerName" json:"CompilerName,omitempty"`
	ErrorMessage string `protobuf:"bytes,2,opt,name=ErrorMessage" json:"ErrorMessage,omitempty"`
	Output       []byte `protobuf:"bytes,3,opt,name=Output,proto3" json:"Output,omitempty"`
}

func (m *CompileFailure) Reset()                    { *m = CompileFailure{} }
func (m *CompileFailure) String() string            { return proto.CompactTextString(m) }
func (*CompileFailure) ProtoMessage()               {}
func (*CompileFailure) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *CompileFailure) GetCompilerName() string {
	if m != nil {
		return m.CompilerName
	}
	return ""
}

func (m *CompileFailure) GetErrorMessage() string {
	if m != nil {
		return m.ErrorMessage
	}
	return ""
}

func (m *CompileFailure) GetOutput() []byte {
	if m != nil {
		return m.Output
	}
	return nil
}

func init() {
	proto.RegisterType((*ProtoFileInfo)(nil), "protorenderer2.ProtoFileInfo")
	proto.RegisterType((*ProtoMessageInfo)(nil), "protorenderer2.ProtoMessageInfo")
	proto.RegisterType((*ProtoFieldInfo)(nil), "protorenderer2.ProtoFieldInfo")
	proto.RegisterType((*DBProtoFile)(nil), "protorenderer2.DBProtoFile")
	proto.RegisterType((*SQLMessage)(nil), "protorenderer2.SQLMessage")
	proto.RegisterType((*ProtocRequest)(nil), "protorenderer2.ProtocRequest")
	proto.RegisterType((*FileTransfer)(nil), "protorenderer2.FileTransfer")
	proto.RegisterType((*FileResult)(nil), "protorenderer2.FileResult")
	proto.RegisterType((*CompileFailure)(nil), "protorenderer2.CompileFailure")
	proto.RegisterEnum("protorenderer2.FieldModifier", FieldModifier_name, FieldModifier_value)
	proto.RegisterEnum("protorenderer2.FieldType", FieldType_name, FieldType_value)
	proto.RegisterEnum("protorenderer2.ProtoFieldPrimitive", ProtoFieldPrimitive_name, ProtoFieldPrimitive_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for ProtoRenderer2 service

type ProtoRenderer2Client interface {
	// protoc-gen-meta2 submits files to protorenderer server for analysis. not intented to be used by clients other than protoc-gen-meta2
	InternalMetaSubmit(ctx context.Context, in *ProtocRequest, opts ...grpc.CallOption) (*common.Void, error)
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

func (c *protoRenderer2Client) InternalMetaSubmit(ctx context.Context, in *ProtocRequest, opts ...grpc.CallOption) (*common.Void, error) {
	out := new(common.Void)
	err := grpc.Invoke(ctx, "/protorenderer2.ProtoRenderer2/InternalMetaSubmit", in, out, c.cc, opts...)
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
	Send(*FileTransfer) error
	Recv() (*FileTransfer, error)
	grpc.ClientStream
}

type protoRenderer2CompileClient struct {
	grpc.ClientStream
}

func (x *protoRenderer2CompileClient) Send(m *FileTransfer) error {
	return x.ClientStream.SendMsg(m)
}

func (x *protoRenderer2CompileClient) Recv() (*FileTransfer, error) {
	m := new(FileTransfer)
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
	// protoc-gen-meta2 submits files to protorenderer server for analysis. not intented to be used by clients other than protoc-gen-meta2
	InternalMetaSubmit(context.Context, *ProtocRequest) (*common.Void, error)
	// compile
	Compile(ProtoRenderer2_CompileServer) error
	// add or update a ".proto" file in the renderers database
	UpdateProto(context.Context, *protorenderer.AddProtoRequest) (*protorenderer.AddProtoResponse, error)
}

func RegisterProtoRenderer2Server(s *grpc.Server, srv ProtoRenderer2Server) {
	s.RegisterService(&_ProtoRenderer2_serviceDesc, srv)
}

func _ProtoRenderer2_InternalMetaSubmit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProtocRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProtoRenderer2Server).InternalMetaSubmit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protorenderer2.ProtoRenderer2/InternalMetaSubmit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProtoRenderer2Server).InternalMetaSubmit(ctx, req.(*ProtocRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProtoRenderer2_Compile_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ProtoRenderer2Server).Compile(&protoRenderer2CompileServer{stream})
}

type ProtoRenderer2_CompileServer interface {
	Send(*FileTransfer) error
	Recv() (*FileTransfer, error)
	grpc.ServerStream
}

type protoRenderer2CompileServer struct {
	grpc.ServerStream
}

func (x *protoRenderer2CompileServer) Send(m *FileTransfer) error {
	return x.ServerStream.SendMsg(m)
}

func (x *protoRenderer2CompileServer) Recv() (*FileTransfer, error) {
	m := new(FileTransfer)
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
			MethodName: "InternalMetaSubmit",
			Handler:    _ProtoRenderer2_InternalMetaSubmit_Handler,
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
	// 1070 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x9c, 0x56, 0xdf, 0x72, 0xdb, 0xc4,
	0x17, 0xfe, 0x49, 0x76, 0xfc, 0xe7, 0xd4, 0xf1, 0x4f, 0x2c, 0x0c, 0xa3, 0x9a, 0x14, 0x3c, 0x86,
	0x61, 0xd2, 0xcc, 0xa0, 0x50, 0x35, 0x64, 0xda, 0x4e, 0xb9, 0xb0, 0x23, 0xbb, 0x23, 0x88, 0xed,
	0xb0, 0x71, 0xd2, 0x29, 0x37, 0x8c, 0x62, 0xad, 0x8d, 0xa8, 0xac, 0x15, 0xd2, 0xba, 0x90, 0x19,
	0xb8, 0xe9, 0x1d, 0x6f, 0xc0, 0xf4, 0x92, 0x37, 0xe0, 0x15, 0xb8, 0xe7, 0x55, 0x78, 0x06, 0x66,
	0x57, 0x2b, 0xd9, 0x52, 0x9c, 0x34, 0xc3, 0x95, 0xf7, 0x9c, 0xfd, 0xbe, 0xb3, 0xe7, 0x8f, 0xf6,
	0xf3, 0x42, 0x3f, 0x8c, 0x28, 0xa3, 0xf1, 0xfe, 0x9c, 0xfa, 0x4e, 0x30, 0x37, 0xa6, 0x34, 0x88,
	0x1c, 0xf7, 0x27, 0x4a, 0x5d, 0x23, 0x20, 0x6c, 0xdf, 0x09, 0xbd, 0x78, 0x5f, 0x20, 0x22, 0x12,
	0xb8, 0x24, 0x22, 0x91, 0x59, 0x30, 0x0d, 0x61, 0xa2, 0x66, 0xde, 0xdb, 0x7a, 0x7a, 0xdb, 0x78,
	0x79, 0x2b, 0x89, 0xd6, 0x32, 0x6e, 0x60, 0x4f, 0xe9, 0x62, 0x41, 0x03, 0xf9, 0x23, 0xf1, 0xe6,
	0x0d, 0xf8, 0xef, 0xcd, 0x79, 0x18, 0xd1, 0x9f, 0x2f, 0xb3, 0x85, 0xe4, 0xb4, 0xe7, 0x94, 0xce,
	0x7d, 0x92, 0x9c, 0x7f, 0xb1, 0x9c, 0xed, 0xbb, 0x24, 0x9e, 0x46, 0x5e, 0xc8, 0xa8, 0xcc, 0xa2,
	0xf3, 0xb7, 0x0a, 0xdb, 0x27, 0x7c, 0x35, 0xf0, 0x7c, 0x62, 0x07, 0x33, 0x8a, 0x1e, 0x43, 0x3d,
	0x73, 0xe8, 0x4a, 0x5b, 0xd9, 0xbd, 0x63, 0x7e, 0x60, 0x14, 0xfa, 0x61, 0xf5, 0x32, 0x08, 0x5e,
	0xa1, 0xd1, 0x17, 0x50, 0xb5, 0x17, 0x21, 0x8d, 0x58, 0xac, 0xab, 0xed, 0xd2, 0xdb, 0x88, 0x29,
	0x16, 0x0d, 0x01, 0x8e, 0x46, 0xcf, 0xc7, 0x21, 0xf3, 0x68, 0x10, 0xeb, 0x25, 0xc1, 0xfc, 0xac,
	0xc8, 0xcc, 0x25, 0x69, 0xac, 0xf0, 0xfd, 0x80, 0x45, 0x97, 0x78, 0x2d, 0x00, 0x7a, 0x0a, 0xb5,
	0x21, 0x89, 0x63, 0x67, 0x4e, 0x62, 0xbd, 0x2c, 0x82, 0xb5, 0x37, 0x06, 0x93, 0x20, 0x1e, 0x0f,
	0x67, 0x8c, 0xd6, 0x97, 0xf0, 0xff, 0x42, 0x70, 0xa4, 0x41, 0xe9, 0x25, 0xb9, 0x14, 0xbd, 0xa8,
	0x63, 0xbe, 0x44, 0xef, 0xc1, 0xd6, 0x2b, 0xc7, 0x5f, 0x12, 0x5d, 0x15, 0xbe, 0xc4, 0x78, 0xa2,
	0x3e, 0x52, 0x3a, 0x6f, 0x14, 0xd0, 0x8a, 0xd1, 0xd1, 0x01, 0x54, 0xa5, 0x29, 0x1b, 0xda, 0x2a,
	0x26, 0x74, 0xfa, 0xcd, 0xb1, 0x44, 0xe0, 0x14, 0x8a, 0x74, 0xa8, 0x1e, 0xd1, 0xc5, 0x82, 0x04,
	0x4c, 0x1e, 0x93, 0x9a, 0xe8, 0x10, 0x2a, 0x03, 0x8f, 0xf8, 0x6e, 0xda, 0xac, 0x0f, 0xaf, 0x69,
	0x16, 0xf1, 0x5d, 0x51, 0x9d, 0x44, 0x77, 0xfe, 0x2c, 0x41, 0x33, 0xbf, 0x85, 0x10, 0x94, 0x47,
	0xce, 0x82, 0xc8, 0xe2, 0xc4, 0xfa, 0x86, 0x83, 0x1f, 0x43, 0x6d, 0x48, 0x5d, 0x6f, 0xe6, 0x91,
	0x48, 0x2f, 0xb5, 0x95, 0xdd, 0xa6, 0x79, 0xaf, 0x78, 0xb4, 0x08, 0x9d, 0x82, 0x70, 0x06, 0x47,
	0xfb, 0xb0, 0x35, 0xb9, 0x0c, 0xc9, 0x03, 0xbd, 0x2c, 0x78, 0x77, 0x37, 0xf2, 0x38, 0x02, 0x27,
	0xb8, 0x94, 0x60, 0xea, 0x5b, 0xb7, 0x22, 0x98, 0xe8, 0x6b, 0x5e, 0x9c, 0xb7, 0xf0, 0x98, 0xf7,
	0x8a, 0x24, 0x47, 0x55, 0x04, 0xf3, 0xe3, 0xeb, 0xbb, 0x93, 0xe1, 0x71, 0x81, 0x7a, 0x25, 0x98,
	0xa9, 0x57, 0xff, 0x6b, 0x30, 0x13, 0xed, 0x40, 0x7d, 0x7c, 0xf1, 0x03, 0x99, 0x32, 0xdb, 0x7a,
	0xa0, 0xd7, 0xda, 0xca, 0x6e, 0x19, 0xaf, 0x1c, 0xeb, 0xbb, 0xa6, 0x5e, 0xcf, 0xef, 0x9a, 0x1d,
	0x0a, 0x77, 0xd6, 0x2e, 0x0d, 0x6a, 0x82, 0x6a, 0x5b, 0x62, 0x5a, 0x65, 0xac, 0xda, 0x56, 0x36,
	0x3f, 0x75, 0x6d, 0x7e, 0x1d, 0x68, 0x60, 0x12, 0xd2, 0xd8, 0x63, 0x34, 0xba, 0xb4, 0x2d, 0x31,
	0xa9, 0x32, 0xce, 0xf9, 0xf8, 0x8c, 0x4f, 0x9c, 0xe9, 0x4b, 0xfe, 0x49, 0x96, 0x93, 0x19, 0x4b,
	0xb3, 0xf3, 0x9b, 0x02, 0xb0, 0xfa, 0x1c, 0xaf, 0x1c, 0xf8, 0x7c, 0x5d, 0x1e, 0xd4, 0xb7, 0xca,
	0x43, 0x6f, 0xe7, 0xcd, 0xeb, 0xbb, 0x95, 0xa5, 0x17, 0xb0, 0xc3, 0x83, 0x3f, 0x5e, 0xdf, 0x6d,
	0xba, 0x17, 0x02, 0x3b, 0xf3, 0x7c, 0x62, 0x78, 0xee, 0xba, 0x78, 0xa4, 0x95, 0x94, 0x56, 0x95,
	0x74, 0x7e, 0x95, 0xe2, 0x34, 0xc5, 0xe4, 0xc7, 0x25, 0x89, 0x19, 0xfa, 0x14, 0x9a, 0x43, 0xc2,
	0x9c, 0x23, 0xba, 0x08, 0x3d, 0x9f, 0x44, 0x32, 0xb3, 0x3a, 0x2e, 0x78, 0x91, 0x05, 0x90, 0x45,
	0x4e, 0xc5, 0xe8, 0x13, 0x23, 0x51, 0x43, 0x23, 0x55, 0x43, 0x83, 0xef, 0x5a, 0x99, 0x22, 0x0a,
	0x02, 0x5e, 0xe3, 0x75, 0xfe, 0x52, 0xa0, 0xc1, 0x57, 0x93, 0xc8, 0x09, 0xe2, 0x19, 0x89, 0x50,
	0x0b, 0x6a, 0xdc, 0x0e, 0x56, 0x37, 0x26, 0xb3, 0x79, 0xfe, 0x96, 0xc3, 0x1c, 0xd1, 0x93, 0x06,
	0x16, 0x6b, 0xb4, 0x07, 0x5a, 0xca, 0xe5, 0xc9, 0xf9, 0x84, 0x25, 0xf5, 0xd5, 0xf0, 0x15, 0xff,
	0x95, 0xa9, 0xf1, 0xb1, 0x6c, 0x17, 0xa6, 0x66, 0x42, 0x05, 0x93, 0x78, 0xe9, 0x33, 0x71, 0x29,
	0x36, 0xe8, 0x88, 0x10, 0x56, 0x81, 0xc0, 0x12, 0xd9, 0xf9, 0x05, 0x60, 0xe5, 0x45, 0xef, 0x43,
	0x65, 0xe0, 0x78, 0x3e, 0x71, 0x45, 0xfe, 0x35, 0x2c, 0xad, 0x5c, 0x65, 0x6a, 0xa1, 0xb2, 0x27,
	0x50, 0xe3, 0xa8, 0x65, 0x44, 0xae, 0x15, 0x1c, 0xd9, 0x7a, 0x09, 0xc3, 0x19, 0xbe, 0x13, 0x42,
	0x33, 0xbf, 0xc7, 0xeb, 0x4c, 0x07, 0xb5, 0xa6, 0x3c, 0x39, 0x1f, 0xc7, 0xf4, 0xa3, 0x88, 0x46,
	0xa9, 0x6a, 0x26, 0x19, 0xe5, 0x7c, 0xbc, 0x92, 0xf1, 0x92, 0x85, 0x4b, 0x26, 0x3a, 0xda, 0xc0,
	0xd2, 0xda, 0xeb, 0xc2, 0x76, 0x4e, 0x83, 0x90, 0x06, 0x8d, 0xc1, 0xf0, 0xbb, 0xb3, 0x91, 0xd5,
	0x1f, 0xd8, 0xa3, 0xbe, 0xa5, 0xfd, 0x0f, 0x01, 0x54, 0x4e, 0xed, 0xd1, 0xb3, 0xe3, 0xbe, 0xa6,
	0xa0, 0x2a, 0x94, 0x86, 0xdd, 0x13, 0x4d, 0x45, 0x75, 0xd8, 0xea, 0x62, 0xdc, 0x7d, 0xa1, 0x95,
	0xf6, 0x1e, 0x41, 0x3d, 0x53, 0x17, 0x4e, 0x9f, 0xb1, 0x1c, 0x7d, 0x1b, 0xea, 0x27, 0xd8, 0x1e,
	0xda, 0x13, 0xfb, 0x9c, 0x47, 0x00, 0xa8, 0x8c, 0x7b, 0x5f, 0xf5, 0x8f, 0x26, 0x9a, 0xba, 0xf7,
	0xbb, 0x02, 0xef, 0x6e, 0x50, 0x04, 0xf4, 0x0e, 0x6c, 0x9f, 0x0c, 0x26, 0x57, 0x92, 0x98, 0x60,
	0x7b, 0xf4, 0x2c, 0x09, 0x71, 0x66, 0x8f, 0x26, 0x87, 0x07, 0x9a, 0x9a, 0xae, 0x1f, 0x9a, 0x5a,
	0x89, 0xe7, 0xd4, 0x7b, 0x31, 0xe9, 0x9f, 0x6a, 0x65, 0x54, 0x83, 0x72, 0x6f, 0x3c, 0x3e, 0xd6,
	0xb6, 0xb8, 0x33, 0xd9, 0xaf, 0x70, 0x67, 0x7f, 0x74, 0x36, 0xd4, 0xaa, 0x9c, 0x65, 0x8d, 0xcf,
	0x7a, 0xc7, 0x7d, 0xad, 0xc6, 0x01, 0x83, 0xe3, 0x71, 0x77, 0xa2, 0xd5, 0x25, 0xf6, 0xf0, 0x40,
	0x03, 0xf3, 0x1f, 0x45, 0x8a, 0x3f, 0x4e, 0xa7, 0x86, 0xba, 0x80, 0xec, 0x80, 0x91, 0x28, 0x70,
	0x7c, 0x7e, 0x7f, 0x4e, 0x97, 0x17, 0x0b, 0x8f, 0xa1, 0x7b, 0x1b, 0x25, 0x2e, 0xbd, 0x82, 0xad,
	0x86, 0x21, 0x9f, 0x25, 0xe7, 0xd4, 0x73, 0x91, 0x2d, 0xfe, 0x2b, 0xf8, 0xe4, 0xd0, 0xce, 0xa6,
	0x8f, 0x31, 0xfd, 0xcc, 0x5b, 0x37, 0xee, 0xee, 0x2a, 0x9f, 0x2b, 0x68, 0x04, 0x77, 0xce, 0x42,
	0xd7, 0x61, 0x44, 0x9c, 0x87, 0x0a, 0xdf, 0x98, 0xd1, 0x75, 0x5d, 0x99, 0x7e, 0x92, 0xc7, 0x47,
	0xd7, 0xee, 0xc7, 0x21, 0x0d, 0x62, 0xd2, 0x3b, 0x87, 0xfb, 0x01, 0x61, 0xeb, 0xef, 0x25, 0xf9,
	0x82, 0xe2, 0x4f, 0xa6, 0x42, 0x36, 0xdf, 0xde, 0xbf, 0xf5, 0xdb, 0xf0, 0xa2, 0x22, 0xec, 0x87,
	0xff, 0x06, 0x00, 0x00, 0xff, 0xff, 0x0e, 0x40, 0x82, 0xcc, 0x56, 0x0a, 0x00, 0x00,
}
