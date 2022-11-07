// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: protocol.proto

package protocol

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ExposeClient is the client API for Expose service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ExposeClient interface {
	// Sends a expose message
	ExposeCmd(ctx context.Context, in *ExposeRequest, opts ...grpc.CallOption) (*ExposeReply, error)
}

type exposeClient struct {
	cc grpc.ClientConnInterface
}

func NewExposeClient(cc grpc.ClientConnInterface) ExposeClient {
	return &exposeClient{cc}
}

func (c *exposeClient) ExposeCmd(ctx context.Context, in *ExposeRequest, opts ...grpc.CallOption) (*ExposeReply, error) {
	out := new(ExposeReply)
	err := c.cc.Invoke(ctx, "/protocol.Expose/ExposeCmd", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ExposeServer is the server API for Expose service.
// All implementations must embed UnimplementedExposeServer
// for forward compatibility
type ExposeServer interface {
	// Sends a expose message
	ExposeCmd(context.Context, *ExposeRequest) (*ExposeReply, error)
	mustEmbedUnimplementedExposeServer()
}

// UnimplementedExposeServer must be embedded to have forward compatible implementations.
type UnimplementedExposeServer struct {
}

func (UnimplementedExposeServer) ExposeCmd(context.Context, *ExposeRequest) (*ExposeReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExposeCmd not implemented")
}
func (UnimplementedExposeServer) mustEmbedUnimplementedExposeServer() {}

// UnsafeExposeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ExposeServer will
// result in compilation errors.
type UnsafeExposeServer interface {
	mustEmbedUnimplementedExposeServer()
}

func RegisterExposeServer(s grpc.ServiceRegistrar, srv ExposeServer) {
	s.RegisterService(&Expose_ServiceDesc, srv)
}

func _Expose_ExposeCmd_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExposeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExposeServer).ExposeCmd(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.Expose/ExposeCmd",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExposeServer).ExposeCmd(ctx, req.(*ExposeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Expose_ServiceDesc is the grpc.ServiceDesc for Expose service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Expose_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protocol.Expose",
	HandlerType: (*ExposeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ExposeCmd",
			Handler:    _Expose_ExposeCmd_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protocol.proto",
}

// ConnectClient is the client API for Connect service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ConnectClient interface {
	// Sends a connect message
	ConnectCmd(ctx context.Context, in *ConnectRequest, opts ...grpc.CallOption) (*ConnectReply, error)
}

type connectClient struct {
	cc grpc.ClientConnInterface
}

func NewConnectClient(cc grpc.ClientConnInterface) ConnectClient {
	return &connectClient{cc}
}

func (c *connectClient) ConnectCmd(ctx context.Context, in *ConnectRequest, opts ...grpc.CallOption) (*ConnectReply, error) {
	out := new(ConnectReply)
	err := c.cc.Invoke(ctx, "/protocol.Connect/ConnectCmd", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConnectServer is the server API for Connect service.
// All implementations must embed UnimplementedConnectServer
// for forward compatibility
type ConnectServer interface {
	// Sends a connect message
	ConnectCmd(context.Context, *ConnectRequest) (*ConnectReply, error)
	mustEmbedUnimplementedConnectServer()
}

// UnimplementedConnectServer must be embedded to have forward compatible implementations.
type UnimplementedConnectServer struct {
}

func (UnimplementedConnectServer) ConnectCmd(context.Context, *ConnectRequest) (*ConnectReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConnectCmd not implemented")
}
func (UnimplementedConnectServer) mustEmbedUnimplementedConnectServer() {}

// UnsafeConnectServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ConnectServer will
// result in compilation errors.
type UnsafeConnectServer interface {
	mustEmbedUnimplementedConnectServer()
}

func RegisterConnectServer(s grpc.ServiceRegistrar, srv ConnectServer) {
	s.RegisterService(&Connect_ServiceDesc, srv)
}

func _Connect_ConnectCmd_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConnectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConnectServer).ConnectCmd(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.Connect/ConnectCmd",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConnectServer).ConnectCmd(ctx, req.(*ConnectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Connect_ServiceDesc is the grpc.ServiceDesc for Connect service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Connect_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protocol.Connect",
	HandlerType: (*ConnectServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ConnectCmd",
			Handler:    _Connect_ConnectCmd_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protocol.proto",
}

// HelloClient is the client API for Hello service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HelloClient interface {
	// Sends a Hello message
	HelloCmd(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error)
}

type helloClient struct {
	cc grpc.ClientConnInterface
}

func NewHelloClient(cc grpc.ClientConnInterface) HelloClient {
	return &helloClient{cc}
}

func (c *helloClient) HelloCmd(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error) {
	out := new(HelloReply)
	err := c.cc.Invoke(ctx, "/protocol.Hello/HelloCmd", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HelloServer is the server API for Hello service.
// All implementations must embed UnimplementedHelloServer
// for forward compatibility
type HelloServer interface {
	// Sends a Hello message
	HelloCmd(context.Context, *HelloRequest) (*HelloReply, error)
	mustEmbedUnimplementedHelloServer()
}

// UnimplementedHelloServer must be embedded to have forward compatible implementations.
type UnimplementedHelloServer struct {
}

func (UnimplementedHelloServer) HelloCmd(context.Context, *HelloRequest) (*HelloReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HelloCmd not implemented")
}
func (UnimplementedHelloServer) mustEmbedUnimplementedHelloServer() {}

// UnsafeHelloServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HelloServer will
// result in compilation errors.
type UnsafeHelloServer interface {
	mustEmbedUnimplementedHelloServer()
}

func RegisterHelloServer(s grpc.ServiceRegistrar, srv HelloServer) {
	s.RegisterService(&Hello_ServiceDesc, srv)
}

func _Hello_HelloCmd_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HelloServer).HelloCmd(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.Hello/HelloCmd",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HelloServer).HelloCmd(ctx, req.(*HelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Hello_ServiceDesc is the grpc.ServiceDesc for Hello service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Hello_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protocol.Hello",
	HandlerType: (*HelloServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HelloCmd",
			Handler:    _Hello_HelloCmd_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protocol.proto",
}
