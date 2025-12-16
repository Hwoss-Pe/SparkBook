package service

import (
	"Webook/follow/domain"
	"Webook/follow/events"
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
	GetFollowStatics(ctx context.Context, uid int64) (domain.FollowStatics, error)
	GetFollower(ctx context.Context, follower, offset, limit int64) ([]domain.FollowRelation, error)
}

type followRelationService struct {
	repo     repository.FollowRepository
	producer events.Producer
}

func NewFollowRelationService(repo repository.FollowRepository, p events.Producer) FollowRelationService {
	return &followRelationService{
		repo:     repo,
		producer: p,
	}
}
func (f *followRelationService) GetFollower(ctx context.Context, follower, offset, limit int64) ([]domain.FollowRelation, error) {
	return f.repo.GetFollower(ctx, follower, offset, limit)
}
func (f *followRelationService) GetFollowee(ctx context.Context, follower, offset, limit int64) ([]domain.FollowRelation, error) {
	return f.repo.GetFollowee(ctx, follower, offset, limit)
}

func (f *followRelationService) Follow(ctx context.Context, follower, followee int64) error {
	if f.producer != nil {
		return f.producer.ProduceFollow(events.FollowEvent{Follower: follower, Followee: followee})
	}
	return f.repo.AddFollowRelation(ctx, domain.FollowRelation{Followee: followee, Follower: follower})
}

func (f *followRelationService) CancelFollow(ctx context.Context, follower, followee int64) error {
	if f.producer != nil {
		return f.producer.ProduceCancelFollow(events.FollowEvent{Follower: follower, Followee: followee})
	}
	return f.repo.InactiveFollowRelation(ctx, follower, followee)
}

func (f *followRelationService) FollowInfo(ctx context.Context, follower, followee int64) (domain.FollowRelation, error) {
	val, err := f.repo.FollowInfo(ctx, follower, followee)
	return val, err
}

func (f *followRelationService) GetFollowStatics(ctx context.Context, uid int64) (domain.FollowStatics, error) {
	return f.repo.GetFollowStatics(ctx, uid)
}
