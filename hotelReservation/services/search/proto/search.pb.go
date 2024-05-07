// Code generated by protoc-gen-go. DO NOT EDIT.
// source: search/proto/search.proto

/*
Package search is a generated protocol buffer package.

It is generated from these files:

	search/proto/search.proto

It has these top-level messages:

	NearbyRequest
	SearchResult
*/
package search

import (
	fmt "fmt"
	"github.com/AleckDarcy/ContextBus"
	math "math"

	proto "github.com/golang/protobuf/proto"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"

	cb "github.com/AleckDarcy/ContextBus/proto"
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

type NearbyRequest struct {
	Lat       float32     `protobuf:"fixed32,1,opt,name=lat" json:"lat,omitempty"`
	Lon       float32     `protobuf:"fixed32,2,opt,name=lon" json:"lon,omitempty"`
	InDate    string      `protobuf:"bytes,3,opt,name=inDate" json:"inDate,omitempty"`
	OutDate   string      `protobuf:"bytes,4,opt,name=outDate" json:"outDate,omitempty"`
	CBPayload *cb.Payload `protobuf:"bytes,10001,opt,name=CBPayload" json:"CBPayload,omitempty"`
}

func (m *NearbyRequest) Reset()                    { *m = NearbyRequest{} }
func (m *NearbyRequest) String() string            { return proto.CompactTextString(m) }
func (*NearbyRequest) ProtoMessage()               {}
func (*NearbyRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *NearbyRequest) GetLat() float32 {
	if m != nil {
		return m.Lat
	}
	return 0
}

func (m *NearbyRequest) GetLon() float32 {
	if m != nil {
		return m.Lon
	}
	return 0
}

func (m *NearbyRequest) GetInDate() string {
	if m != nil {
		return m.InDate
	}
	return ""
}

func (m *NearbyRequest) GetOutDate() string {
	if m != nil {
		return m.OutDate
	}
	return ""
}

func (m *NearbyRequest) GetCBPayload() *cb.Payload {
	if m != nil {
		return m.CBPayload
	}
	return nil
}

type SearchResult struct {
	HotelIds  []string    `protobuf:"bytes,1,rep,name=hotelIds" json:"hotelIds,omitempty"`
	CBPayload *cb.Payload `protobuf:"bytes,10001,opt,name=CBPayload" json:"CBPayload,omitempty"`
}

func (m *SearchResult) Reset()                    { *m = SearchResult{} }
func (m *SearchResult) String() string            { return proto.CompactTextString(m) }
func (*SearchResult) ProtoMessage()               {}
func (*SearchResult) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *SearchResult) GetHotelIds() []string {
	if m != nil {
		return m.HotelIds
	}
	return nil
}

func (m *SearchResult) GetCBPayload() *cb.Payload {
	if m != nil {
		return m.CBPayload
	}
	return nil
}

func init() {
	proto.RegisterType((*NearbyRequest)(nil), "search.NearbyRequest")
	proto.RegisterType((*SearchResult)(nil), "search.SearchResult")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Search service

type SearchClient interface {
	ResetDB(ctx context.Context, in *NearbyRequest, opts ...grpc.CallOption) (*SearchResult, error)
	Nearby(ctx context.Context, in *NearbyRequest, opts ...grpc.CallOption) (*SearchResult, error)
}

type searchClient struct {
	cc *grpc.ClientConn
}

func NewSearchClient(cc *grpc.ClientConn) SearchClient {
	return &searchClient{cc}
}

func (c *searchClient) ResetDB(ctx context.Context, in *NearbyRequest, opts ...grpc.CallOption) (*SearchResult, error) {
	out := new(SearchResult)
	err := grpc.Invoke(ctx, "/search.Search/ResetDB", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchClient) Nearby(ctx context.Context, in *NearbyRequest, opts ...grpc.CallOption) (*SearchResult, error) {
	out := new(SearchResult)

	// ContextBus
	// fmt.Println("Send /search.Search/Nearby request")

	err := grpc.Invoke(ctx, "/search.Search/Nearby", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}

	// ContextBus
	// fmt.Println("Receive /search.Search/Nearby response")

	return out, nil
}

// Server API for Search service

type SearchServer interface {
	ResetDB(context.Context, *NearbyRequest) (*SearchResult, error)
	Nearby(context.Context, *NearbyRequest) (*SearchResult, error)
}

func RegisterSearchServer(s *grpc.Server, srv SearchServer) {
	s.RegisterService(&_Search_serviceDesc, srv)
}

func _Search_ResetDB_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NearbyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServer).ResetDB(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/search.Search/ResetDB",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServer).ResetDB(ctx, req.(*NearbyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Search_Nearby_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NearbyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}

	// ContextBus
	cbCtx, cbOK := ContextBus.FromPayload(ctx, in.GetCBPayload())
	_, _ = cbCtx, cbOK

	// ContextBus
	if cbOK {
		ContextBus.OnSubmission(cbCtx, &cb.EventWhere{}, &cb.EventRecorder{
			Type: cb.EventRecorderType_EventRecorderThirdParty,
			Name: "Search.Nearby.Handler.1",
		}, &cb.EventMessage{
			Attrs:   nil,
			Message: "NearbyHandler starts",
			Paths:   nil,
		})
	}

	if interceptor == nil {
		return srv.(SearchServer).Nearby(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/search.Search/Nearby",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServer).Nearby(ctx, req.(*NearbyRequest))
	}
	val, err := interceptor(ctx, in, info, handler)
	// ContextBus
	if cbOK { // todo: do payload
		res := val.(*SearchResult)

		_ = res
		ContextBus.OnSubmission(cbCtx, &cb.EventWhere{}, &cb.EventRecorder{
			Type: cb.EventRecorderType_EventRecorderThirdParty,
			Name: "Search.Nearby.Handler.2",
		}, &cb.EventMessage{
			Attrs:   nil,
			Message: "NearbyHandler ends",
			Paths:   nil,
		})
	}

	return val, err
}

var _Search_serviceDesc = grpc.ServiceDesc{
	ServiceName: "search.Search",
	HandlerType: (*SearchServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ResetDB",
			Handler:    _Search_ResetDB_Handler,
		},
		{
			MethodName: "Nearby",
			Handler:    _Search_Nearby_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "search/proto/search.proto",
}

func init() { proto.RegisterFile("search/proto/search.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 257 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x91, 0xc1, 0x4a, 0x03, 0x31,
	0x10, 0x86, 0x49, 0x57, 0xd2, 0xee, 0xa8, 0x20, 0xa1, 0x4a, 0x5c, 0x10, 0x96, 0x9e, 0xf6, 0xb4,
	0x0b, 0x2d, 0xfa, 0x00, 0x6b, 0x2f, 0x5e, 0x8a, 0xc4, 0x9b, 0x97, 0x92, 0x6d, 0x07, 0x2a, 0x2c,
	0x1b, 0xdd, 0x4c, 0xd0, 0x3e, 0x86, 0x57, 0x9f, 0x56, 0x4c, 0x52, 0x6d, 0x8f, 0xf6, 0x36, 0xdf,
	0xff, 0xff, 0x61, 0x26, 0x33, 0x70, 0x6d, 0x51, 0xf7, 0xab, 0x4d, 0xf5, 0xda, 0x1b, 0x32, 0x55,
	0x80, 0xd2, 0x83, 0xe0, 0x81, 0xb2, 0x9b, 0x95, 0xe9, 0x08, 0x3f, 0x68, 0xd9, 0x38, 0x5b, 0xed,
	0xd5, 0x21, 0x36, 0xf9, 0x62, 0x70, 0xbe, 0x40, 0xdd, 0x37, 0x5b, 0x85, 0x6f, 0x0e, 0x2d, 0x89,
	0x0b, 0x48, 0x5a, 0x4d, 0x92, 0xe5, 0xac, 0x18, 0xa8, 0x9f, 0xd2, 0x2b, 0xa6, 0x93, 0x83, 0xa8,
	0x98, 0x4e, 0x5c, 0x01, 0x7f, 0xe9, 0xe6, 0x9a, 0x50, 0x26, 0x39, 0x2b, 0x52, 0x15, 0x49, 0x48,
	0x18, 0x1a, 0x47, 0xde, 0x38, 0xf1, 0xc6, 0x0e, 0xc5, 0x0c, 0xd2, 0xfb, 0xfa, 0x51, 0x6f, 0x5b,
	0xa3, 0xd7, 0xf2, 0x73, 0x91, 0xb3, 0xe2, 0x74, 0x3a, 0x2e, 0xf7, 0xe7, 0x89, 0xa6, 0xfa, 0xcb,
	0x4d, 0x96, 0x70, 0xf6, 0xe4, 0x7f, 0xa1, 0xd0, 0xba, 0x96, 0x44, 0x06, 0xa3, 0x8d, 0x21, 0x6c,
	0x1f, 0xd6, 0x56, 0xb2, 0x3c, 0x29, 0x52, 0xf5, 0xcb, 0x47, 0x35, 0x98, 0xbe, 0x03, 0x0f, 0x0d,
	0xc4, 0x1d, 0x0c, 0x15, 0x5a, 0xa4, 0x79, 0x2d, 0x2e, 0xcb, 0xb8, 0xc8, 0x83, 0xbd, 0x64, 0xe3,
	0x9d, 0x7c, 0x30, 0xd2, 0x2d, 0xf0, 0x10, 0xfb, 0xd7, 0xb3, 0x7a, 0xf4, 0x1c, 0xef, 0xd3, 0x70,
	0x7f, 0x87, 0xd9, 0x77, 0x00, 0x00, 0x00, 0xff, 0xff, 0xe6, 0x96, 0xc8, 0x4f, 0xcb, 0x01, 0x00,
	0x00,
}
