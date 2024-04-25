// Code generated by protoc-gen-go. DO NOT EDIT.
// source: services/user/proto/user.proto

/*
Package user is a generated protocol buffer package.

It is generated from these files:

	services/user/proto/user.proto

It has these top-level messages:

	Request
	Result
*/
package user

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

type Request struct {
	Username string `protobuf:"bytes,1,opt,name=username" json:"username,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password" json:"password,omitempty"`
}

func (m *Request) Reset()                    { *m = Request{} }
func (m *Request) String() string            { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()               {}
func (*Request) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Request) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *Request) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type Result struct {
	Correct bool `protobuf:"varint,1,opt,name=correct" json:"correct,omitempty"`
}

func (m *Result) Reset()                    { *m = Result{} }
func (m *Result) String() string            { return proto.CompactTextString(m) }
func (*Result) ProtoMessage()               {}
func (*Result) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Result) GetCorrect() bool {
	if m != nil {
		return m.Correct
	}
	return false
}

func init() {
	proto.RegisterType((*Request)(nil), "user.Request")
	proto.RegisterType((*Result)(nil), "user.Result")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for User service

type UserClient interface {
	ResetDB(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Result, error)
	// CheckUser returns whether the username and password are correct
	CheckUser(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Result, error)
}

type userClient struct {
	cc *grpc.ClientConn
}

func NewUserClient(cc *grpc.ClientConn) UserClient {
	return &userClient{cc}
}

func (c *userClient) ResetDB(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := grpc.Invoke(ctx, "/user.User/ResetDB", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) CheckUser(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := grpc.Invoke(ctx, "/user.User/CheckUser", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for User service

type UserServer interface {
	ResetDB(context.Context, *Request) (*Result, error)
	// CheckUser returns whether the username and password are correct
	CheckUser(context.Context, *Request) (*Result, error)
}

func RegisterUserServer(s *grpc.Server, srv UserServer) {
	s.RegisterService(&_User_serviceDesc, srv)
}

func _User_ResetDB_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).ResetDB(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.User/ResetDB",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).ResetDB(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_CheckUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).CheckUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.User/CheckUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).CheckUser(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

var _User_serviceDesc = grpc.ServiceDesc{
	ServiceName: "user.User",
	HandlerType: (*UserServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ResetDB",
			Handler:    _User_ResetDB_Handler,
		},
		{
			MethodName: "CheckUser",
			Handler:    _User_CheckUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "services/user/proto/user.proto",
}

func init() { proto.RegisterFile("services/user/proto/user.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 178 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x2b, 0x4e, 0x2d, 0x2a,
	0xcb, 0x4c, 0x4e, 0x2d, 0xd6, 0x2f, 0x2d, 0x4e, 0x2d, 0xd2, 0x2f, 0x28, 0xca, 0x2f, 0xc9, 0x07,
	0x33, 0xf5, 0xc0, 0x4c, 0x21, 0x16, 0x10, 0x5b, 0xc9, 0x91, 0x8b, 0x3d, 0x28, 0xb5, 0xb0, 0x34,
	0xb5, 0xb8, 0x44, 0x48, 0x8a, 0x8b, 0x03, 0x24, 0x94, 0x97, 0x98, 0x9b, 0x2a, 0xc1, 0xa8, 0xc0,
	0xa8, 0xc1, 0x19, 0x04, 0xe7, 0x83, 0xe4, 0x0a, 0x12, 0x8b, 0x8b, 0xcb, 0xf3, 0x8b, 0x52, 0x24,
	0x98, 0x20, 0x72, 0x30, 0xbe, 0x92, 0x12, 0x17, 0x5b, 0x50, 0x6a, 0x71, 0x69, 0x4e, 0x89, 0x90,
	0x04, 0x17, 0x7b, 0x72, 0x7e, 0x51, 0x51, 0x6a, 0x72, 0x09, 0xd8, 0x00, 0x8e, 0x20, 0x18, 0xd7,
	0x28, 0x82, 0x8b, 0x25, 0xb4, 0x38, 0xb5, 0x48, 0x48, 0x0d, 0x64, 0x5d, 0x71, 0x6a, 0x89, 0x8b,
	0x93, 0x10, 0xaf, 0x1e, 0xd8, 0x31, 0x50, 0xdb, 0xa5, 0x78, 0x60, 0x5c, 0xb0, 0x49, 0x1a, 0x5c,
	0x9c, 0xce, 0x19, 0xa9, 0xc9, 0xd9, 0x60, 0x4d, 0xf8, 0x54, 0x3a, 0xb1, 0x45, 0x81, 0x3d, 0x92,
	0xc4, 0x06, 0xf6, 0x95, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0x9d, 0x25, 0xc8, 0xa5, 0xf7, 0x00,
	0x00, 0x00,
}
