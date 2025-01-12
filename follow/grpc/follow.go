package grpc2

import (
	followv1 "Webook/api/proto/gen/api/proto/follow/v1"
	"Webook/follow/domain"
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
	err := f.svc.Follow(ctx, request.Follower, request.Followee)
	return &followv1.FollowResponse{}, err
}

func (f *FollowServiceServer) Register(server grpc.ServiceRegistrar) {
	followv1.RegisterFollowServiceServer(server, f)
}

func (f *FollowServiceServer) CancelFollow(ctx context.Context, request *followv1.CancelFollowRequest) (*followv1.CancelFollowResponse, error) {
	err := f.svc.CancelFollow(ctx, request.Follower, request.Followee)
	return &followv1.CancelFollowResponse{}, err
}

// GetFollowee 查看自己的关注列表
func (f *FollowServiceServer) GetFollowee(ctx context.Context, request *followv1.GetFolloweeRequest) (*followv1.GetFolloweeResponse, error) {
	list, err := f.svc.GetFollowee(ctx, request.Follower, request.Offset, request.Limit)
	if err != nil {
		return nil, err
	}
	res := make([]*followv1.FollowRelation, 0, len(list))
	for _, relation := range list {
		res = append(res, f.convertToView(relation))
	}
	return &followv1.GetFolloweeResponse{
		FollowRelations: res,
	}, nil
}

func (f *FollowServiceServer) FollowInfo(ctx context.Context, request *followv1.FollowInfoRequest) (*followv1.FollowInfoResponse, error) {
	info, err := f.svc.FollowInfo(ctx, request.Follower, request.Followee)
	if err != nil {
		return nil, err
	}
	return &followv1.FollowInfoResponse{
		FollowRelation: f.convertToView(info),
	}, nil
}

// GetFollower 查看自己的粉丝列表
func (f *FollowServiceServer) GetFollower(ctx context.Context, request *followv1.GetFollowerRequest) (*followv1.GetFollowerResponse, error) {
	list, err := f.svc.GetFollower(ctx, request.Followee, request.Offset, request.Limit)
	if err != nil {
		return nil, err
	}

	res := make([]*followv1.FollowRelation, 0, len(list))
	for _, relation := range list {
		res = append(res, f.convertToView(relation))
	}
	return &followv1.GetFollowerResponse{
		FollowRelations: res,
	}, nil
}

func (f *FollowServiceServer) GetFollowStatics(ctx context.Context, request *followv1.GetFollowStaticRequest) (*followv1.GetFollowStaticResponse, error) {
	followee := request.Followee
	fs, err := f.svc.GetFollowStatics(ctx, followee)
	if err != nil {
		return nil, err
	}
	return &followv1.GetFollowStaticResponse{
		FollowStatic: f.convertToViewStatic(fs),
	}, nil
}

func (f *FollowServiceServer) convertToView(relation domain.FollowRelation) *followv1.FollowRelation {
	return &followv1.FollowRelation{
		Followee: relation.Followee,
		Follower: relation.Follower,
	}
}
func (f *FollowServiceServer) convertToViewStatic(relation domain.FollowStatics) *followv1.FollowStatic {
	return &followv1.FollowStatic{
		Followers: relation.Followers,
		Followees: relation.Followers,
	}
}
