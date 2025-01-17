// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: api/proto/ranking/v1/ranking.proto

package rankingv1

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

// RankingServiceClient is the client API for RankingService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RankingServiceClient interface {
	RankTopN(ctx context.Context, in *RankTopNRequest, opts ...grpc.CallOption) (*RankTopNResponse, error)
	TopN(ctx context.Context, in *TopNRequest, opts ...grpc.CallOption) (*TopNResponse, error)
}

type rankingServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRankingServiceClient(cc grpc.ClientConnInterface) RankingServiceClient {
	return &rankingServiceClient{cc}
}

func (c *rankingServiceClient) RankTopN(ctx context.Context, in *RankTopNRequest, opts ...grpc.CallOption) (*RankTopNResponse, error) {
	out := new(RankTopNResponse)
	err := c.cc.Invoke(ctx, "/ranking.v1.RankingService/RankTopN", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rankingServiceClient) TopN(ctx context.Context, in *TopNRequest, opts ...grpc.CallOption) (*TopNResponse, error) {
	out := new(TopNResponse)
	err := c.cc.Invoke(ctx, "/ranking.v1.RankingService/TopN", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RankingServiceServer is the server API for RankingService service.
// All implementations must embed UnimplementedRankingServiceServer
// for forward compatibility
type RankingServiceServer interface {
	RankTopN(context.Context, *RankTopNRequest) (*RankTopNResponse, error)
	TopN(context.Context, *TopNRequest) (*TopNResponse, error)
	mustEmbedUnimplementedRankingServiceServer()
}

// UnimplementedRankingServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRankingServiceServer struct {
}

func (UnimplementedRankingServiceServer) RankTopN(context.Context, *RankTopNRequest) (*RankTopNResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RankTopN not implemented")
}
func (UnimplementedRankingServiceServer) TopN(context.Context, *TopNRequest) (*TopNResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TopN not implemented")
}
func (UnimplementedRankingServiceServer) mustEmbedUnimplementedRankingServiceServer() {}

// UnsafeRankingServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RankingServiceServer will
// result in compilation errors.
type UnsafeRankingServiceServer interface {
	mustEmbedUnimplementedRankingServiceServer()
}

func RegisterRankingServiceServer(s grpc.ServiceRegistrar, srv RankingServiceServer) {
	s.RegisterService(&RankingService_ServiceDesc, srv)
}

func _RankingService_RankTopN_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RankTopNRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RankingServiceServer).RankTopN(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ranking.v1.RankingService/RankTopN",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RankingServiceServer).RankTopN(ctx, req.(*RankTopNRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RankingService_TopN_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TopNRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RankingServiceServer).TopN(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ranking.v1.RankingService/TopN",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RankingServiceServer).TopN(ctx, req.(*TopNRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RankingService_ServiceDesc is the grpc.ServiceDesc for RankingService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RankingService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ranking.v1.RankingService",
	HandlerType: (*RankingServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RankTopN",
			Handler:    _RankingService_RankTopN_Handler,
		},
		{
			MethodName: "TopN",
			Handler:    _RankingService_TopN_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/proto/ranking/v1/ranking.proto",
}
