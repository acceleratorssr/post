// interactive.proto

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: intr/v1/intr.proto

package intrv1

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
	LikeService_IncrReadCount_FullMethodName       = "/intr.v1.LikeService/IncrReadCount"
	LikeService_Like_FullMethodName                = "/intr.v1.LikeService/Like"
	LikeService_UnLike_FullMethodName              = "/intr.v1.LikeService/UnLike"
	LikeService_Collect_FullMethodName             = "/intr.v1.LikeService/Collect"
	LikeService_GetListBatchOfLikes_FullMethodName = "/intr.v1.LikeService/GetListBatchOfLikes"
)

// LikeServiceClient is the client API for LikeService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LikeServiceClient interface {
	IncrReadCount(ctx context.Context, in *IncrReadCountRequest, opts ...grpc.CallOption) (*IncrReadCountResponse, error)
	Like(ctx context.Context, in *LikeRequest, opts ...grpc.CallOption) (*LikeResponse, error)
	UnLike(ctx context.Context, in *UnLikeRequest, opts ...grpc.CallOption) (*UnLikeResponse, error)
	Collect(ctx context.Context, in *CollectRequest, opts ...grpc.CallOption) (*CollectResponse, error)
	GetListBatchOfLikes(ctx context.Context, in *GetListBatchOfLikesRequest, opts ...grpc.CallOption) (*GetListBatchOfLikesResponse, error)
}

type likeServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLikeServiceClient(cc grpc.ClientConnInterface) LikeServiceClient {
	return &likeServiceClient{cc}
}

func (c *likeServiceClient) IncrReadCount(ctx context.Context, in *IncrReadCountRequest, opts ...grpc.CallOption) (*IncrReadCountResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(IncrReadCountResponse)
	err := c.cc.Invoke(ctx, LikeService_IncrReadCount_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *likeServiceClient) Like(ctx context.Context, in *LikeRequest, opts ...grpc.CallOption) (*LikeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LikeResponse)
	err := c.cc.Invoke(ctx, LikeService_Like_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *likeServiceClient) UnLike(ctx context.Context, in *UnLikeRequest, opts ...grpc.CallOption) (*UnLikeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UnLikeResponse)
	err := c.cc.Invoke(ctx, LikeService_UnLike_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *likeServiceClient) Collect(ctx context.Context, in *CollectRequest, opts ...grpc.CallOption) (*CollectResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CollectResponse)
	err := c.cc.Invoke(ctx, LikeService_Collect_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *likeServiceClient) GetListBatchOfLikes(ctx context.Context, in *GetListBatchOfLikesRequest, opts ...grpc.CallOption) (*GetListBatchOfLikesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetListBatchOfLikesResponse)
	err := c.cc.Invoke(ctx, LikeService_GetListBatchOfLikes_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LikeServiceServer is the server API for LikeService service.
// All implementations must embed UnimplementedLikeServiceServer
// for forward compatibility.
type LikeServiceServer interface {
	IncrReadCount(context.Context, *IncrReadCountRequest) (*IncrReadCountResponse, error)
	Like(context.Context, *LikeRequest) (*LikeResponse, error)
	UnLike(context.Context, *UnLikeRequest) (*UnLikeResponse, error)
	Collect(context.Context, *CollectRequest) (*CollectResponse, error)
	GetListBatchOfLikes(context.Context, *GetListBatchOfLikesRequest) (*GetListBatchOfLikesResponse, error)
	mustEmbedUnimplementedLikeServiceServer()
}

// UnimplementedLikeServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedLikeServiceServer struct{}

func (UnimplementedLikeServiceServer) IncrReadCount(context.Context, *IncrReadCountRequest) (*IncrReadCountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IncrReadCount not implemented")
}
func (UnimplementedLikeServiceServer) Like(context.Context, *LikeRequest) (*LikeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Like not implemented")
}
func (UnimplementedLikeServiceServer) UnLike(context.Context, *UnLikeRequest) (*UnLikeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnLike not implemented")
}
func (UnimplementedLikeServiceServer) Collect(context.Context, *CollectRequest) (*CollectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Collect not implemented")
}
func (UnimplementedLikeServiceServer) GetListBatchOfLikes(context.Context, *GetListBatchOfLikesRequest) (*GetListBatchOfLikesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetListBatchOfLikes not implemented")
}
func (UnimplementedLikeServiceServer) mustEmbedUnimplementedLikeServiceServer() {}
func (UnimplementedLikeServiceServer) testEmbeddedByValue()                     {}

// UnsafeLikeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LikeServiceServer will
// result in compilation errors.
type UnsafeLikeServiceServer interface {
	mustEmbedUnimplementedLikeServiceServer()
}

func RegisterLikeServiceServer(s grpc.ServiceRegistrar, srv LikeServiceServer) {
	// If the following call pancis, it indicates UnimplementedLikeServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&LikeService_ServiceDesc, srv)
}

func _LikeService_IncrReadCount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IncrReadCountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LikeServiceServer).IncrReadCount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LikeService_IncrReadCount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LikeServiceServer).IncrReadCount(ctx, req.(*IncrReadCountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LikeService_Like_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LikeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LikeServiceServer).Like(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LikeService_Like_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LikeServiceServer).Like(ctx, req.(*LikeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LikeService_UnLike_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnLikeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LikeServiceServer).UnLike(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LikeService_UnLike_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LikeServiceServer).UnLike(ctx, req.(*UnLikeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LikeService_Collect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CollectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LikeServiceServer).Collect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LikeService_Collect_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LikeServiceServer).Collect(ctx, req.(*CollectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LikeService_GetListBatchOfLikes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetListBatchOfLikesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LikeServiceServer).GetListBatchOfLikes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LikeService_GetListBatchOfLikes_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LikeServiceServer).GetListBatchOfLikes(ctx, req.(*GetListBatchOfLikesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// LikeService_ServiceDesc is the grpc.ServiceDesc for LikeService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LikeService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "intr.v1.LikeService",
	HandlerType: (*LikeServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IncrReadCount",
			Handler:    _LikeService_IncrReadCount_Handler,
		},
		{
			MethodName: "Like",
			Handler:    _LikeService_Like_Handler,
		},
		{
			MethodName: "UnLike",
			Handler:    _LikeService_UnLike_Handler,
		},
		{
			MethodName: "Collect",
			Handler:    _LikeService_Collect_Handler,
		},
		{
			MethodName: "GetListBatchOfLikes",
			Handler:    _LikeService_GetListBatchOfLikes_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "intr/v1/intr.proto",
}
