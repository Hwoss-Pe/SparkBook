package service

import (
	"Webook/follow/domain"
	"Webook/follow/repository"
	"golang.org/x/net/context"
)

type FollowRelationService interface {
	// GetFollowee 获得某个人的关注列表
	GetFollowee(ctx context.Context, follower, offset, limit int64) ([]domain.FollowRelation, error)
	Follow(ctx context.Context, follower, followee int64) error
	CancelFollow(ctx context.Context, follower, followee int64) error
	// FollowInfo 获取某个人关注另一个人的详细信息
	FollowInfo(ctx context.Context, follower, followee int64) (domain.FollowRelation, error)
}

type followRelationService struct {
	repo repository.FollowRepository
}

func (f *followRelationService) GetFollowee(ctx context.Context, follower, offset, limit int64) ([]domain.FollowRelation, error) {
	//TODO implement me
	panic("implement me")
}

func (f *followRelationService) Follow(ctx context.Context, follower, followee int64) error {
	//TODO implement me
	panic("implement me")
}

func (f *followRelationService) CancelFollow(ctx context.Context, follower, followee int64) error {
	//TODO implement me
	panic("implement me")
}

func (f *followRelationService) FollowInfo(ctx context.Context, follower, followee int64) (domain.FollowRelation, error) {
	//TODO implement me
	panic("implement me")
}
