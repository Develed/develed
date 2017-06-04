// Code generated by protoc-gen-go.
// source: services.proto
// DO NOT EDIT!

/*
Package services is a generated protocol buffer package.

It is generated from these files:
	services.proto

It has these top-level messages:
	DrawRequest
	DrawResponse
	TextRequest
	TextResponse
*/
package services

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

type DrawRequest struct {
	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *DrawRequest) Reset()                    { *m = DrawRequest{} }
func (m *DrawRequest) String() string            { return proto.CompactTextString(m) }
func (*DrawRequest) ProtoMessage()               {}
func (*DrawRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *DrawRequest) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type DrawResponse struct {
	Code   int32  `protobuf:"varint,1,opt,name=code" json:"code,omitempty"`
	Status string `protobuf:"bytes,2,opt,name=status" json:"status,omitempty"`
}

func (m *DrawResponse) Reset()                    { *m = DrawResponse{} }
func (m *DrawResponse) String() string            { return proto.CompactTextString(m) }
func (*DrawResponse) ProtoMessage()               {}
func (*DrawResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *DrawResponse) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *DrawResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

type TextRequest struct {
	Text      string `protobuf:"bytes,1,opt,name=text" json:"text,omitempty"`
	Font      string `protobuf:"bytes,2,opt,name=font" json:"font,omitempty"`
	FontColor int64  `protobuf:"varint,3,opt,name=font_color,json=fontColor" json:"font_color,omitempty"`
	FontBg    int64  `protobuf:"varint,4,opt,name=font_bg,json=fontBg" json:"font_bg,omitempty"`
}

func (m *TextRequest) Reset()                    { *m = TextRequest{} }
func (m *TextRequest) String() string            { return proto.CompactTextString(m) }
func (*TextRequest) ProtoMessage()               {}
func (*TextRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *TextRequest) GetText() string {
	if m != nil {
		return m.Text
	}
	return ""
}

func (m *TextRequest) GetFont() string {
	if m != nil {
		return m.Font
	}
	return ""
}

func (m *TextRequest) GetFontColor() int64 {
	if m != nil {
		return m.FontColor
	}
	return 0
}

func (m *TextRequest) GetFontBg() int64 {
	if m != nil {
		return m.FontBg
	}
	return 0
}

type TextResponse struct {
	Code   int32  `protobuf:"varint,1,opt,name=code" json:"code,omitempty"`
	Status string `protobuf:"bytes,2,opt,name=status" json:"status,omitempty"`
}

func (m *TextResponse) Reset()                    { *m = TextResponse{} }
func (m *TextResponse) String() string            { return proto.CompactTextString(m) }
func (*TextResponse) ProtoMessage()               {}
func (*TextResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *TextResponse) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *TextResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func init() {
	proto.RegisterType((*DrawRequest)(nil), "services.DrawRequest")
	proto.RegisterType((*DrawResponse)(nil), "services.DrawResponse")
	proto.RegisterType((*TextRequest)(nil), "services.TextRequest")
	proto.RegisterType((*TextResponse)(nil), "services.TextResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for ImageSink service

type ImageSinkClient interface {
	Draw(ctx context.Context, in *DrawRequest, opts ...grpc.CallOption) (*DrawResponse, error)
}

type imageSinkClient struct {
	cc *grpc.ClientConn
}

func NewImageSinkClient(cc *grpc.ClientConn) ImageSinkClient {
	return &imageSinkClient{cc}
}

func (c *imageSinkClient) Draw(ctx context.Context, in *DrawRequest, opts ...grpc.CallOption) (*DrawResponse, error) {
	out := new(DrawResponse)
	err := grpc.Invoke(ctx, "/services.ImageSink/Draw", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ImageSink service

type ImageSinkServer interface {
	Draw(context.Context, *DrawRequest) (*DrawResponse, error)
}

func RegisterImageSinkServer(s *grpc.Server, srv ImageSinkServer) {
	s.RegisterService(&_ImageSink_serviceDesc, srv)
}

func _ImageSink_Draw_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DrawRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ImageSinkServer).Draw(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.ImageSink/Draw",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ImageSinkServer).Draw(ctx, req.(*DrawRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ImageSink_serviceDesc = grpc.ServiceDesc{
	ServiceName: "services.ImageSink",
	HandlerType: (*ImageSinkServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Draw",
			Handler:    _ImageSink_Draw_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "services.proto",
}

// Client API for Textd service

type TextdClient interface {
	Write(ctx context.Context, in *TextRequest, opts ...grpc.CallOption) (*TextResponse, error)
}

type textdClient struct {
	cc *grpc.ClientConn
}

func NewTextdClient(cc *grpc.ClientConn) TextdClient {
	return &textdClient{cc}
}

func (c *textdClient) Write(ctx context.Context, in *TextRequest, opts ...grpc.CallOption) (*TextResponse, error) {
	out := new(TextResponse)
	err := grpc.Invoke(ctx, "/services.Textd/Write", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Textd service

type TextdServer interface {
	Write(context.Context, *TextRequest) (*TextResponse, error)
}

func RegisterTextdServer(s *grpc.Server, srv TextdServer) {
	s.RegisterService(&_Textd_serviceDesc, srv)
}

func _Textd_Write_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TextRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TextdServer).Write(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.Textd/Write",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TextdServer).Write(ctx, req.(*TextRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Textd_serviceDesc = grpc.ServiceDesc{
	ServiceName: "services.Textd",
	HandlerType: (*TextdServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Write",
			Handler:    _Textd_Write_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "services.proto",
}

func init() { proto.RegisterFile("services.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 258 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x51, 0xb1, 0x4e, 0xc3, 0x40,
	0x0c, 0x6d, 0x68, 0x12, 0x88, 0x5b, 0x31, 0x9c, 0x44, 0x89, 0x90, 0x90, 0x42, 0xa6, 0x4c, 0x1d,
	0xca, 0x00, 0x62, 0x03, 0xba, 0xb0, 0x1e, 0x48, 0x8c, 0xe8, 0x9a, 0x98, 0x28, 0x82, 0xe6, 0xca,
	0x9d, 0x0b, 0xfd, 0x7c, 0x64, 0xa7, 0xa0, 0x13, 0x1b, 0x53, 0x9e, 0x9f, 0xed, 0xe7, 0x97, 0x77,
	0x70, 0xec, 0xd1, 0x7d, 0x76, 0x35, 0xfa, 0xf9, 0xc6, 0x59, 0xb2, 0xea, 0xe8, 0xa7, 0x2e, 0x2f,
	0x60, 0xb2, 0x74, 0xe6, 0x4b, 0xe3, 0xc7, 0x16, 0x3d, 0x29, 0x05, 0x71, 0x63, 0xc8, 0xe4, 0x51,
	0x11, 0x55, 0x53, 0x2d, 0xb8, 0xbc, 0x81, 0xe9, 0x30, 0xe2, 0x37, 0xb6, 0xf7, 0xc8, 0x33, 0xb5,
	0x6d, 0x50, 0x66, 0x12, 0x2d, 0x58, 0xcd, 0x20, 0xf5, 0x64, 0x68, 0xeb, 0xf3, 0x83, 0x22, 0xaa,
	0x32, 0xbd, 0xaf, 0xca, 0x35, 0x4c, 0x9e, 0x70, 0x47, 0x81, 0x3c, 0xe1, 0x8e, 0x64, 0x35, 0xd3,
	0x82, 0x99, 0x7b, 0xb5, 0x3d, 0xed, 0x17, 0x05, 0xab, 0x73, 0x00, 0xfe, 0xbe, 0xd4, 0xf6, 0xdd,
	0xba, 0x7c, 0x5c, 0x44, 0xd5, 0x58, 0x67, 0xcc, 0xdc, 0x33, 0xa1, 0x4e, 0xe1, 0x50, 0xda, 0xab,
	0x36, 0x8f, 0xa5, 0x97, 0x72, 0x79, 0xd7, 0xb2, 0xd5, 0xe1, 0xdc, 0xff, 0xad, 0x2e, 0x96, 0x90,
	0x3d, 0xac, 0x4d, 0x8b, 0x8f, 0x5d, 0xff, 0xa6, 0xae, 0x20, 0xe6, 0x7f, 0x56, 0x27, 0xf3, 0xdf,
	0xe4, 0x82, 0x98, 0xce, 0x66, 0x7f, 0xe9, 0xe1, 0x5e, 0x39, 0x5a, 0xdc, 0x42, 0xc2, 0x0e, 0x1a,
	0x75, 0x0d, 0xc9, 0xb3, 0xeb, 0x08, 0x43, 0x89, 0x20, 0x8a, 0x50, 0x22, 0xb4, 0x5c, 0x8e, 0x56,
	0xa9, 0xbc, 0xd1, 0xe5, 0x77, 0x00, 0x00, 0x00, 0xff, 0xff, 0x6e, 0x73, 0xf6, 0xbf, 0xb5, 0x01,
	0x00, 0x00,
}