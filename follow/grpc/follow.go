package grpc2

import (
	followv1 "Webook/api/proto/gen/api/proto/follow/v1"
	"Webook/follow/service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type FollowServiceServer struct {
	svc service.FollowRelationService
	followv1.UnimplementedFollowServiceServer
}

func NewFollowRelationServiceServer(svc service.FollowRelationService) *FollowServiceServer {
	return &FollowServiceServer{
		svc: svc,
	}
}
func (f *FollowServiceServer) Follow(ctx context.Context, request *followv1.FollowRequest) (*followv1.FollowResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (f *FollowServiceServer) Register(server grpc.ServiceRegistrar) {
	followv1.RegisterFollowServiceServer(server, f)
}

func (f *FollowServiceServer) CancelFollow(ctx context.Context, request *followv1.CancelFollowRequest) (*followv1.CancelFollowResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (f *FollowServiceServer) GetFollowee(ctx context.Context, request *followv1.GetFolloweeRequest) (*followv1.GetFolloweeResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (f *FollowServiceServer) FollowInfo(ctx context.Context, request *followv1.FollowInfoRequest) (*followv1.FollowInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (f *FollowServiceServer) GetFollower(ctx context.Context, request *followv1.GetFollowerRequest) (*followv1.GetFollowerResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (f *FollowServiceServer) GetFollowStatics(ctx context.Context, request *followv1.GetFollowStaticRequest) (*followv1.GetFollowStaticResponse, error) {
	//TODO implement me
	panic("implement me")
}
