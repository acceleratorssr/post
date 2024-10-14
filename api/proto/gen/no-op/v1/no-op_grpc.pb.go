// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: no-op/v1/no-op.proto

package noopv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	NoOp_NoOp_FullMethodName = "/noop.v1.NoOp/NoOp"
)

// NoOpClient is the client API for NoOp service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// example:
type NoOpClient interface {
	NoOp(ctx context.Context, in *NoOpRequest, opts ...grpc.CallOption) (*NoOpResponse, error)
}

type noOpClient struct {
	cc grpc.ClientConnInterface
}

func NewNoOpClient(cc grpc.ClientConnInterface) NoOpClient {
	return &noOpClient{cc}
}

func (c *noOpClient) NoOp(ctx context.Context, in *NoOpRequest, opts ...grpc.CallOption) (*NoOpResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(NoOpResponse)
	err := c.cc.Invoke(ctx, NoOp_NoOp_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NoOpServer is the server API for NoOp service.
// All implementations must embed UnimplementedNoOpServer
// for forward compatibility.
//
// example:
type NoOpServer interface {
	NoOp(context.Context, *NoOpRequest) (*NoOpResponse, error)
	mustEmbedUnimplementedNoOpServer()
}

// UnimplementedNoOpServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedNoOpServer struct{}

func (UnimplementedNoOpServer) NoOp(context.Context, *NoOpRequest) (*NoOpResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NoOp not implemented")
}
func (UnimplementedNoOpServer) mustEmbedUnimplementedNoOpServer() {}
func (UnimplementedNoOpServer) testEmbeddedByValue()              {}

// UnsafeNoOpServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NoOpServer will
// result in compilation errors.
type UnsafeNoOpServer interface {
	mustEmbedUnimplementedNoOpServer()
}

func RegisterNoOpServer(s grpc.ServiceRegistrar, srv NoOpServer) {
	// If the following call pancis, it indicates UnimplementedNoOpServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&NoOp_ServiceDesc, srv)
}

func _NoOp_NoOp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NoOpRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NoOpServer).NoOp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NoOp_NoOp_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NoOpServer).NoOp(ctx, req.(*NoOpRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// NoOp_ServiceDesc is the grpc.ServiceDesc for NoOp service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NoOp_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "noop.v1.NoOp",
	HandlerType: (*NoOpServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NoOp",
			Handler:    _NoOp_NoOp_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "no-op/v1/no-op.proto",
}
