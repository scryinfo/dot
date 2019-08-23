// Code generated by protoc-gen-go. DO NOT EDIT.
// source: hi.proto

package go_out

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

//--------------------------------------------------------------------
//findDot
type ReqDirs struct {
	Dirs                 []string `protobuf:"bytes,1,rep,name=dirs,proto3" json:"dirs,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReqDirs) Reset()         { *m = ReqDirs{} }
func (m *ReqDirs) String() string { return proto.CompactTextString(m) }
func (*ReqDirs) ProtoMessage()    {}
func (*ReqDirs) Descriptor() ([]byte, []int) {
	return fileDescriptor_d092a8920edeec73, []int{0}
}

func (m *ReqDirs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReqDirs.Unmarshal(m, b)
}
func (m *ReqDirs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReqDirs.Marshal(b, m, deterministic)
}
func (m *ReqDirs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReqDirs.Merge(m, src)
}
func (m *ReqDirs) XXX_Size() int {
	return xxx_messageInfo_ReqDirs.Size(m)
}
func (m *ReqDirs) XXX_DiscardUnknown() {
	xxx_messageInfo_ReqDirs.DiscardUnknown(m)
}

var xxx_messageInfo_ReqDirs proto.InternalMessageInfo

func (m *ReqDirs) GetDirs() []string {
	if m != nil {
		return m.Dirs
	}
	return nil
}

type ResDots struct {
	DotsInfo             string   `protobuf:"bytes,1,opt,name=dotsInfo,proto3" json:"dotsInfo,omitempty"`
	NoExistDirs          []string `protobuf:"bytes,2,rep,name=noExistDirs,proto3" json:"noExistDirs,omitempty"`
	Error                string   `protobuf:"bytes,3,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ResDots) Reset()         { *m = ResDots{} }
func (m *ResDots) String() string { return proto.CompactTextString(m) }
func (*ResDots) ProtoMessage()    {}
func (*ResDots) Descriptor() ([]byte, []int) {
	return fileDescriptor_d092a8920edeec73, []int{1}
}

func (m *ResDots) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResDots.Unmarshal(m, b)
}
func (m *ResDots) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResDots.Marshal(b, m, deterministic)
}
func (m *ResDots) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResDots.Merge(m, src)
}
func (m *ResDots) XXX_Size() int {
	return xxx_messageInfo_ResDots.Size(m)
}
func (m *ResDots) XXX_DiscardUnknown() {
	xxx_messageInfo_ResDots.DiscardUnknown(m)
}

var xxx_messageInfo_ResDots proto.InternalMessageInfo

func (m *ResDots) GetDotsInfo() string {
	if m != nil {
		return m.DotsInfo
	}
	return ""
}

func (m *ResDots) GetNoExistDirs() []string {
	if m != nil {
		return m.NoExistDirs
	}
	return nil
}

func (m *ResDots) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

//---------------------------------------------------------------------
//importfile
type ReqImport struct {
	Filepath             string   `protobuf:"bytes,1,opt,name=filepath,proto3" json:"filepath,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReqImport) Reset()         { *m = ReqImport{} }
func (m *ReqImport) String() string { return proto.CompactTextString(m) }
func (*ReqImport) ProtoMessage()    {}
func (*ReqImport) Descriptor() ([]byte, []int) {
	return fileDescriptor_d092a8920edeec73, []int{2}
}

func (m *ReqImport) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReqImport.Unmarshal(m, b)
}
func (m *ReqImport) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReqImport.Marshal(b, m, deterministic)
}
func (m *ReqImport) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReqImport.Merge(m, src)
}
func (m *ReqImport) XXX_Size() int {
	return xxx_messageInfo_ReqImport.Size(m)
}
func (m *ReqImport) XXX_DiscardUnknown() {
	xxx_messageInfo_ReqImport.DiscardUnknown(m)
}

var xxx_messageInfo_ReqImport proto.InternalMessageInfo

func (m *ReqImport) GetFilepath() string {
	if m != nil {
		return m.Filepath
	}
	return ""
}

type ResImport struct {
	Json                 string   `protobuf:"bytes,1,opt,name=json,proto3" json:"json,omitempty"`
	Error                string   `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ResImport) Reset()         { *m = ResImport{} }
func (m *ResImport) String() string { return proto.CompactTextString(m) }
func (*ResImport) ProtoMessage()    {}
func (*ResImport) Descriptor() ([]byte, []int) {
	return fileDescriptor_d092a8920edeec73, []int{3}
}

func (m *ResImport) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResImport.Unmarshal(m, b)
}
func (m *ResImport) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResImport.Marshal(b, m, deterministic)
}
func (m *ResImport) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResImport.Merge(m, src)
}
func (m *ResImport) XXX_Size() int {
	return xxx_messageInfo_ResImport.Size(m)
}
func (m *ResImport) XXX_DiscardUnknown() {
	xxx_messageInfo_ResImport.DiscardUnknown(m)
}

var xxx_messageInfo_ResImport proto.InternalMessageInfo

func (m *ResImport) GetJson() string {
	if m != nil {
		return m.Json
	}
	return ""
}

func (m *ResImport) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

//----------------------------------------------------------------------
//exportfile
type ReqExport struct {
	Configdata           string   `protobuf:"bytes,1,opt,name=configdata,proto3" json:"configdata,omitempty"`
	Filename             []string `protobuf:"bytes,2,rep,name=filename,proto3" json:"filename,omitempty"`
	Dotdata              string   `protobuf:"bytes,3,opt,name=dotdata,proto3" json:"dotdata,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReqExport) Reset()         { *m = ReqExport{} }
func (m *ReqExport) String() string { return proto.CompactTextString(m) }
func (*ReqExport) ProtoMessage()    {}
func (*ReqExport) Descriptor() ([]byte, []int) {
	return fileDescriptor_d092a8920edeec73, []int{4}
}

func (m *ReqExport) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReqExport.Unmarshal(m, b)
}
func (m *ReqExport) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReqExport.Marshal(b, m, deterministic)
}
func (m *ReqExport) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReqExport.Merge(m, src)
}
func (m *ReqExport) XXX_Size() int {
	return xxx_messageInfo_ReqExport.Size(m)
}
func (m *ReqExport) XXX_DiscardUnknown() {
	xxx_messageInfo_ReqExport.DiscardUnknown(m)
}

var xxx_messageInfo_ReqExport proto.InternalMessageInfo

func (m *ReqExport) GetConfigdata() string {
	if m != nil {
		return m.Configdata
	}
	return ""
}

func (m *ReqExport) GetFilename() []string {
	if m != nil {
		return m.Filename
	}
	return nil
}

func (m *ReqExport) GetDotdata() string {
	if m != nil {
		return m.Dotdata
	}
	return ""
}

type ResExport struct {
	Error                string   `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ResExport) Reset()         { *m = ResExport{} }
func (m *ResExport) String() string { return proto.CompactTextString(m) }
func (*ResExport) ProtoMessage()    {}
func (*ResExport) Descriptor() ([]byte, []int) {
	return fileDescriptor_d092a8920edeec73, []int{5}
}

func (m *ResExport) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResExport.Unmarshal(m, b)
}
func (m *ResExport) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResExport.Marshal(b, m, deterministic)
}
func (m *ResExport) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResExport.Merge(m, src)
}
func (m *ResExport) XXX_Size() int {
	return xxx_messageInfo_ResExport.Size(m)
}
func (m *ResExport) XXX_DiscardUnknown() {
	xxx_messageInfo_ResExport.DiscardUnknown(m)
}

var xxx_messageInfo_ResExport proto.InternalMessageInfo

func (m *ResExport) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func init() {
	proto.RegisterType((*ReqDirs)(nil), "go_out.ReqDirs")
	proto.RegisterType((*ResDots)(nil), "go_out.ResDots")
	proto.RegisterType((*ReqImport)(nil), "go_out.ReqImport")
	proto.RegisterType((*ResImport)(nil), "go_out.ResImport")
	proto.RegisterType((*ReqExport)(nil), "go_out.ReqExport")
	proto.RegisterType((*ResExport)(nil), "go_out.ResExport")
}

func init() { proto.RegisterFile("hi.proto", fileDescriptor_d092a8920edeec73) }

var fileDescriptor_d092a8920edeec73 = []byte{
	// 324 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x52, 0xcd, 0x4a, 0xf3, 0x40,
	0x14, 0x25, 0xfd, 0xcf, 0xed, 0xc7, 0x27, 0x0e, 0x2e, 0x42, 0x41, 0x89, 0xd9, 0x58, 0x10, 0x02,
	0x5a, 0xf5, 0x01, 0x6a, 0x2a, 0x76, 0x9b, 0xa5, 0x20, 0x12, 0x4d, 0xd2, 0x8e, 0xd8, 0xb9, 0xe9,
	0xcc, 0x08, 0xf5, 0x15, 0x7c, 0x6a, 0x99, 0xbf, 0x38, 0x8a, 0x82, 0xbb, 0xfb, 0x77, 0xce, 0xb9,
	0xe7, 0xce, 0xc0, 0x68, 0x4d, 0xd3, 0x86, 0xa3, 0x44, 0x32, 0x58, 0xe1, 0x03, 0xbe, 0xca, 0xe4,
	0x10, 0x86, 0x79, 0xb5, 0xcd, 0x28, 0x17, 0x84, 0x40, 0xaf, 0xa4, 0x5c, 0x44, 0x41, 0xdc, 0x9d,
	0x86, 0xb9, 0x8e, 0x93, 0x7b, 0xd5, 0x16, 0x19, 0x4a, 0x41, 0x26, 0x30, 0x2a, 0x51, 0x8a, 0x25,
	0xab, 0x31, 0x0a, 0xe2, 0x60, 0x1a, 0xe6, 0x6d, 0x4e, 0x62, 0x18, 0x33, 0x5c, 0xec, 0xa8, 0x90,
	0x8a, 0x29, 0xea, 0x68, 0x06, 0xbf, 0x44, 0x0e, 0xa0, 0x5f, 0x71, 0x8e, 0x3c, 0xea, 0x6a, 0xa8,
	0x49, 0x92, 0x13, 0x08, 0xf3, 0x6a, 0xbb, 0xdc, 0x34, 0xc8, 0xa5, 0x12, 0xa8, 0xe9, 0x4b, 0xd5,
	0x14, 0x72, 0xed, 0x04, 0x5c, 0x9e, 0x5c, 0xaa, 0x41, 0x61, 0x07, 0x09, 0xf4, 0x9e, 0x05, 0x32,
	0x3b, 0xa4, 0xe3, 0x4f, 0xfe, 0x8e, 0xcf, 0x5f, 0x68, 0xfe, 0xc5, 0x4e, 0xc3, 0x8e, 0x00, 0x9e,
	0x90, 0xd5, 0x74, 0x55, 0x16, 0xb2, 0xb0, 0x60, 0xaf, 0xe2, 0xf4, 0x59, 0xb1, 0xa9, 0xac, 0x83,
	0x36, 0x27, 0x11, 0x0c, 0x4b, 0x94, 0x1a, 0x68, 0x0c, 0xb8, 0x34, 0x39, 0xd6, 0x9b, 0x59, 0x89,
	0x76, 0x8b, 0xc0, 0xdb, 0xe2, 0xfc, 0xbd, 0x03, 0xfd, 0x5b, 0x9a, 0xa1, 0x24, 0xa7, 0x30, 0xbc,
	0xa1, 0xac, 0x54, 0xe1, 0x5e, 0x6a, 0x5e, 0x20, 0xb5, 0xe7, 0x9f, 0x78, 0x05, 0x73, 0xf0, 0x2b,
	0xf8, 0x6f, 0x0c, 0xcf, 0xdf, 0xae, 0xf5, 0x96, 0x64, 0xdf, 0xc3, 0x98, 0xd6, 0xc4, 0x2b, 0xb9,
	0xf3, 0xcc, 0x60, 0xec, 0x70, 0x4a, 0xe8, 0x6f, 0xa0, 0x0b, 0xf8, 0x67, 0x3c, 0xfc, 0x20, 0x65,
	0x1a, 0x5f, 0x50, 0xd6, 0xef, 0x19, 0x84, 0x26, 0xfa, 0x2e, 0xf4, 0x2b, 0x64, 0x3e, 0xba, 0xb3,
	0x5f, 0xef, 0x71, 0xa0, 0x7f, 0xe2, 0xec, 0x23, 0x00, 0x00, 0xff, 0xff, 0xa8, 0xbd, 0xb1, 0x7d,
	0x95, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// HiDotClient is the client API for HiDot service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type HiDotClient interface {
	FindDot(ctx context.Context, in *ReqDirs, opts ...grpc.CallOption) (*ResDots, error)
	//文件导入
	ImportByConfig(ctx context.Context, in *ReqImport, opts ...grpc.CallOption) (*ResImport, error)
	ImportByDot(ctx context.Context, in *ReqImport, opts ...grpc.CallOption) (*ResImport, error)
	//导出文件
	ExportConfig(ctx context.Context, in *ReqExport, opts ...grpc.CallOption) (*ResExport, error)
	ExportDot(ctx context.Context, in *ReqExport, opts ...grpc.CallOption) (*ResExport, error)
}

type hiDotClient struct {
	cc *grpc.ClientConn
}

func NewHiDotClient(cc *grpc.ClientConn) HiDotClient {
	return &hiDotClient{cc}
}

func (c *hiDotClient) FindDot(ctx context.Context, in *ReqDirs, opts ...grpc.CallOption) (*ResDots, error) {
	out := new(ResDots)
	err := c.cc.Invoke(ctx, "/go_out.HiDot/FindDot", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hiDotClient) ImportByConfig(ctx context.Context, in *ReqImport, opts ...grpc.CallOption) (*ResImport, error) {
	out := new(ResImport)
	err := c.cc.Invoke(ctx, "/go_out.HiDot/ImportByConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hiDotClient) ImportByDot(ctx context.Context, in *ReqImport, opts ...grpc.CallOption) (*ResImport, error) {
	out := new(ResImport)
	err := c.cc.Invoke(ctx, "/go_out.HiDot/ImportByDot", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hiDotClient) ExportConfig(ctx context.Context, in *ReqExport, opts ...grpc.CallOption) (*ResExport, error) {
	out := new(ResExport)
	err := c.cc.Invoke(ctx, "/go_out.HiDot/ExportConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hiDotClient) ExportDot(ctx context.Context, in *ReqExport, opts ...grpc.CallOption) (*ResExport, error) {
	out := new(ResExport)
	err := c.cc.Invoke(ctx, "/go_out.HiDot/ExportDot", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HiDotServer is the server API for HiDot service.
type HiDotServer interface {
	FindDot(context.Context, *ReqDirs) (*ResDots, error)
	//文件导入
	ImportByConfig(context.Context, *ReqImport) (*ResImport, error)
	ImportByDot(context.Context, *ReqImport) (*ResImport, error)
	//导出文件
	ExportConfig(context.Context, *ReqExport) (*ResExport, error)
	ExportDot(context.Context, *ReqExport) (*ResExport, error)
}

// UnimplementedHiDotServer can be embedded to have forward compatible implementations.
type UnimplementedHiDotServer struct {
}

func (*UnimplementedHiDotServer) FindDot(ctx context.Context, req *ReqDirs) (*ResDots, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindDot not implemented")
}
func (*UnimplementedHiDotServer) ImportByConfig(ctx context.Context, req *ReqImport) (*ResImport, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ImportByConfig not implemented")
}
func (*UnimplementedHiDotServer) ImportByDot(ctx context.Context, req *ReqImport) (*ResImport, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ImportByDot not implemented")
}
func (*UnimplementedHiDotServer) ExportConfig(ctx context.Context, req *ReqExport) (*ResExport, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExportConfig not implemented")
}
func (*UnimplementedHiDotServer) ExportDot(ctx context.Context, req *ReqExport) (*ResExport, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExportDot not implemented")
}

func RegisterHiDotServer(s *grpc.Server, srv HiDotServer) {
	s.RegisterService(&_HiDot_serviceDesc, srv)
}

func _HiDot_FindDot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReqDirs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HiDotServer).FindDot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/go_out.HiDot/FindDot",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HiDotServer).FindDot(ctx, req.(*ReqDirs))
	}
	return interceptor(ctx, in, info, handler)
}

func _HiDot_ImportByConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReqImport)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HiDotServer).ImportByConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/go_out.HiDot/ImportByConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HiDotServer).ImportByConfig(ctx, req.(*ReqImport))
	}
	return interceptor(ctx, in, info, handler)
}

func _HiDot_ImportByDot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReqImport)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HiDotServer).ImportByDot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/go_out.HiDot/ImportByDot",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HiDotServer).ImportByDot(ctx, req.(*ReqImport))
	}
	return interceptor(ctx, in, info, handler)
}

func _HiDot_ExportConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReqExport)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HiDotServer).ExportConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/go_out.HiDot/ExportConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HiDotServer).ExportConfig(ctx, req.(*ReqExport))
	}
	return interceptor(ctx, in, info, handler)
}

func _HiDot_ExportDot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReqExport)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HiDotServer).ExportDot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/go_out.HiDot/ExportDot",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HiDotServer).ExportDot(ctx, req.(*ReqExport))
	}
	return interceptor(ctx, in, info, handler)
}

var _HiDot_serviceDesc = grpc.ServiceDesc{
	ServiceName: "go_out.HiDot",
	HandlerType: (*HiDotServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FindDot",
			Handler:    _HiDot_FindDot_Handler,
		},
		{
			MethodName: "ImportByConfig",
			Handler:    _HiDot_ImportByConfig_Handler,
		},
		{
			MethodName: "ImportByDot",
			Handler:    _HiDot_ImportByDot_Handler,
		},
		{
			MethodName: "ExportConfig",
			Handler:    _HiDot_ExportConfig_Handler,
		},
		{
			MethodName: "ExportDot",
			Handler:    _HiDot_ExportDot_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "hi.proto",
}
