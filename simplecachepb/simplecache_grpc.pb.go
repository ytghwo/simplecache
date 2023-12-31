// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.19.4
// source: simplecache.proto

package simplecachepb

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

const (
	SimpleCache_Get_FullMethodName = "/simplecachepb.SimpleCache/Get"
)

// SimpleCacheClient is the client API for SimpleCache service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SimpleCacheClient interface {
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
}

type simpleCacheClient struct {
	cc grpc.ClientConnInterface
}

func NewSimpleCacheClient(cc grpc.ClientConnInterface) SimpleCacheClient {
	return &simpleCacheClient{cc}
}

func (c *simpleCacheClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, SimpleCache_Get_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SimpleCacheServer is the server API for SimpleCache service.
// All implementations must embed UnimplementedSimpleCacheServer
// for forward compatibility
type SimpleCacheServer interface {
	Get(context.Context, *GetRequest) (*GetResponse, error)
	mustEmbedUnimplementedSimpleCacheServer()
}

// UnimplementedSimpleCacheServer must be embedded to have forward compatible implementations.
type UnimplementedSimpleCacheServer struct {
}

func (UnimplementedSimpleCacheServer) Get(context.Context, *GetRequest) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedSimpleCacheServer) mustEmbedUnimplementedSimpleCacheServer() {}

// UnsafeSimpleCacheServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SimpleCacheServer will
// result in compilation errors.
type UnsafeSimpleCacheServer interface {
	mustEmbedUnimplementedSimpleCacheServer()
}

func RegisterSimpleCacheServer(s grpc.ServiceRegistrar, srv SimpleCacheServer) {
	s.RegisterService(&SimpleCache_ServiceDesc, srv)
}

func _SimpleCache_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SimpleCacheServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SimpleCache_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SimpleCacheServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SimpleCache_ServiceDesc is the grpc.ServiceDesc for SimpleCache service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SimpleCache_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "simplecachepb.SimpleCache",
	HandlerType: (*SimpleCacheServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _SimpleCache_Get_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "simplecache.proto",
}
