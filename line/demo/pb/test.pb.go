// Code generated by protoc-gen-go. DO NOT EDIT.
// source: test.proto

package proto

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

// The request message containing the user's name.
type TestRequest struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TestRequest) Reset()         { *m = TestRequest{} }
func (m *TestRequest) String() string { return proto.CompactTextString(m) }
func (*TestRequest) ProtoMessage()    {}
func (*TestRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_test_bda11bd0903c5340, []int{0}
}
func (m *TestRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TestRequest.Unmarshal(m, b)
}
func (m *TestRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TestRequest.Marshal(b, m, deterministic)
}
func (dst *TestRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TestRequest.Merge(dst, src)
}
func (m *TestRequest) XXX_Size() int {
	return xxx_messageInfo_TestRequest.Size(m)
}
func (m *TestRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_TestRequest.DiscardUnknown(m)
}

var xxx_messageInfo_TestRequest proto.InternalMessageInfo

func (m *TestRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

// The response message containing the greetings
type TestReply struct {
	Message              string   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TestReply) Reset()         { *m = TestReply{} }
func (m *TestReply) String() string { return proto.CompactTextString(m) }
func (*TestReply) ProtoMessage()    {}
func (*TestReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_test_bda11bd0903c5340, []int{1}
}
func (m *TestReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TestReply.Unmarshal(m, b)
}
func (m *TestReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TestReply.Marshal(b, m, deterministic)
}
func (dst *TestReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TestReply.Merge(dst, src)
}
func (m *TestReply) XXX_Size() int {
	return xxx_messageInfo_TestReply.Size(m)
}
func (m *TestReply) XXX_DiscardUnknown() {
	xxx_messageInfo_TestReply.DiscardUnknown(m)
}

var xxx_messageInfo_TestReply proto.InternalMessageInfo

func (m *TestReply) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*TestRequest)(nil), "proto.TestRequest")
	proto.RegisterType((*TestReply)(nil), "proto.TestReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// TestClient is the client API for Test service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TestClient interface {
	// Sends a greeting
	SayHello(ctx context.Context, in *TestRequest, opts ...grpc.CallOption) (*TestReply, error)
}

type testClient struct {
	cc *grpc.ClientConn
}

func NewTestClient(cc *grpc.ClientConn) TestClient {
	return &testClient{cc}
}

func (c *testClient) SayHello(ctx context.Context, in *TestRequest, opts ...grpc.CallOption) (*TestReply, error) {
	out := new(TestReply)
	err := c.cc.Invoke(ctx, "/proto.Test/SayHello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TestServer is the server API for Test service.
type TestServer interface {
	// Sends a greeting
	SayHello(context.Context, *TestRequest) (*TestReply, error)
}

func RegisterTestServer(s *grpc.Server, srv TestServer) {
	s.RegisterService(&_Test_serviceDesc, srv)
}

func _Test_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Test/SayHello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestServer).SayHello(ctx, req.(*TestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Test_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Test",
	HandlerType: (*TestServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _Test_SayHello_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "test.proto",
}

func init() { proto.RegisterFile("test.proto", fileDescriptor_test_bda11bd0903c5340) }

var fileDescriptor_test_bda11bd0903c5340 = []byte{
	// 133 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2a, 0x49, 0x2d, 0x2e,
	0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x53, 0x4a, 0x8a, 0x5c, 0xdc, 0x21, 0xa9,
	0xc5, 0x25, 0x41, 0xa9, 0x85, 0xa5, 0xa9, 0xc5, 0x25, 0x42, 0x42, 0x5c, 0x2c, 0x79, 0x89, 0xb9,
	0xa9, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x60, 0xb6, 0x92, 0x2a, 0x17, 0x27, 0x44, 0x49,
	0x41, 0x4e, 0xa5, 0x90, 0x04, 0x17, 0x7b, 0x6e, 0x6a, 0x71, 0x71, 0x62, 0x3a, 0x4c, 0x0d, 0x8c,
	0x6b, 0x64, 0xc5, 0xc5, 0x02, 0x52, 0x26, 0x64, 0xc4, 0xc5, 0x11, 0x9c, 0x58, 0xe9, 0x91, 0x9a,
	0x93, 0x93, 0x2f, 0x24, 0x04, 0xb1, 0x4c, 0x0f, 0xc9, 0x0a, 0x29, 0x01, 0x14, 0xb1, 0x82, 0x9c,
	0x4a, 0x25, 0x86, 0x24, 0x36, 0xb0, 0x90, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0xb5, 0x95, 0xe4,
	0x47, 0xa1, 0x00, 0x00, 0x00,
}
