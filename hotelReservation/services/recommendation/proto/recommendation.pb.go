// Code generated by protoc-gen-go. DO NOT EDIT.
// source: services/recommendation/proto/recommendation.proto

/*
Package recommendation is a generated protocol buffer package.

It is generated from these files:

	services/recommendation/proto/recommendation.proto

It has these top-level messages:

	Request
	Result
*/
package recommendation

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

// The requirement of the recommendation.
type Request struct {
	Require string  `protobuf:"bytes,1,opt,name=require" json:"require,omitempty"`
	Lat     float64 `protobuf:"fixed64,2,opt,name=lat" json:"lat,omitempty"`
	Lon     float64 `protobuf:"fixed64,3,opt,name=lon" json:"lon,omitempty"`
}

func (m *Request) Reset()                    { *m = Request{} }
func (m *Request) String() string            { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()               {}
func (*Request) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Request) GetRequire() string {
	if m != nil {
		return m.Require
	}
	return ""
}

func (m *Request) GetLat() float64 {
	if m != nil {
		return m.Lat
	}
	return 0
}

func (m *Request) GetLon() float64 {
	if m != nil {
		return m.Lon
	}
	return 0
}

type Result struct {
	HotelIds []string `protobuf:"bytes,1,rep,name=HotelIds" json:"HotelIds,omitempty"`
}

func (m *Result) Reset()                    { *m = Result{} }
func (m *Result) String() string            { return proto.CompactTextString(m) }
func (*Result) ProtoMessage()               {}
func (*Result) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Result) GetHotelIds() []string {
	if m != nil {
		return m.HotelIds
	}
	return nil
}

func init() {
	proto.RegisterType((*Request)(nil), "recommendation.Request")
	proto.RegisterType((*Result)(nil), "recommendation.Result")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Recommendation service

type RecommendationClient interface {
	ResetDB(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Result, error)
	// GetRecommendations returns recommended hotels for a given requirement
	GetRecommendations(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Result, error)
}

type recommendationClient struct {
	cc *grpc.ClientConn
}

func NewRecommendationClient(cc *grpc.ClientConn) RecommendationClient {
	return &recommendationClient{cc}
}

func (c *recommendationClient) ResetDB(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := grpc.Invoke(ctx, "/recommendation.Recommendation/ResetDB", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *recommendationClient) GetRecommendations(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := grpc.Invoke(ctx, "/recommendation.Recommendation/GetRecommendations", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Recommendation service

type RecommendationServer interface {
	ResetDB(context.Context, *Request) (*Result, error)
	// GetRecommendations returns recommended hotels for a given requirement
	GetRecommendations(context.Context, *Request) (*Result, error)
}

func RegisterRecommendationServer(s *grpc.Server, srv RecommendationServer) {
	s.RegisterService(&_Recommendation_serviceDesc, srv)
}

func _Recommendation_ResetDB_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecommendationServer).ResetDB(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/recommendation.Recommendation/ResetDB",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecommendationServer).ResetDB(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Recommendation_GetRecommendations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecommendationServer).GetRecommendations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/recommendation.Recommendation/GetRecommendations",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecommendationServer).GetRecommendations(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

var _Recommendation_serviceDesc = grpc.ServiceDesc{
	ServiceName: "recommendation.Recommendation",
	HandlerType: (*RecommendationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ResetDB",
			Handler:    _Recommendation_ResetDB_Handler,
		},
		{
			MethodName: "GetRecommendations",
			Handler:    _Recommendation_GetRecommendations_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "services/recommendation/proto/recommendation.proto",
}

func init() {
	proto.RegisterFile("services/recommendation/proto/recommendation.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 205 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x32, 0x2a, 0x4e, 0x2d, 0x2a,
	0xcb, 0x4c, 0x4e, 0x2d, 0xd6, 0x2f, 0x4a, 0x4d, 0xce, 0xcf, 0xcd, 0x4d, 0xcd, 0x4b, 0x49, 0x2c,
	0xc9, 0xcc, 0xcf, 0xd3, 0x2f, 0x28, 0xca, 0x2f, 0xc9, 0x47, 0x13, 0xd4, 0x03, 0x0b, 0x0a, 0xf1,
	0xa1, 0x8a, 0x2a, 0xb9, 0x73, 0xb1, 0x07, 0xa5, 0x16, 0x96, 0xa6, 0x16, 0x97, 0x08, 0x49, 0x70,
	0xb1, 0x17, 0xa5, 0x16, 0x96, 0x66, 0x16, 0xa5, 0x4a, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0xc1,
	0xb8, 0x42, 0x02, 0x5c, 0xcc, 0x39, 0x89, 0x25, 0x12, 0x4c, 0x0a, 0x8c, 0x1a, 0x8c, 0x41, 0x20,
	0x26, 0x58, 0x24, 0x3f, 0x4f, 0x82, 0x19, 0x2a, 0x92, 0x9f, 0xa7, 0xa4, 0xc2, 0xc5, 0x16, 0x94,
	0x5a, 0x5c, 0x9a, 0x53, 0x22, 0x24, 0xc5, 0xc5, 0xe1, 0x91, 0x5f, 0x92, 0x9a, 0xe3, 0x99, 0x52,
	0x2c, 0xc1, 0xa8, 0xc0, 0xac, 0xc1, 0x19, 0x04, 0xe7, 0x1b, 0x4d, 0x66, 0xe4, 0xe2, 0x0b, 0x42,
	0x71, 0x81, 0x90, 0x15, 0xc8, 0x05, 0xc5, 0xa9, 0x25, 0x2e, 0x4e, 0x42, 0xe2, 0x7a, 0x68, 0x6e,
	0x86, 0x3a, 0x4d, 0x4a, 0x0c, 0x53, 0x02, 0x6c, 0x95, 0x2b, 0x97, 0x90, 0x7b, 0x6a, 0x09, 0xaa,
	0x81, 0xc5, 0x24, 0x1b, 0xe3, 0x24, 0x10, 0x85, 0x16, 0x2c, 0x49, 0x6c, 0xe0, 0xd0, 0x32, 0x06,
	0x04, 0x00, 0x00, 0xff, 0xff, 0xe0, 0x0b, 0xb1, 0x91, 0x63, 0x01, 0x00, 0x00,
}
