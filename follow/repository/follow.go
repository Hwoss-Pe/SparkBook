package repository

import (
	"Webook/follow/domain"
	"Webook/follow/repository/cache"
	"Webook/follow/repository/dao"
	"Webook/pkg/logger"
	"golang.org/x/net/context"
)

type FollowRepository interface {
	GetFollowee(ctx context.Context, follower, offset, limit int64) ([]domain.FollowRelation, error)
	// AddFollowRelation 创建关注关系
	AddFollowRelation(ctx context.Context, f domain.FollowRelation) error
	// InactiveFollowRelation 取消关注
	InactiveFollowRelation(ctx context.Context, follower int64, followee int64) error
	GetFollowStatics(ctx context.Context, uid int64) (domain.FollowStatics, error)
	FollowInfo(ctx context.Context, follower int64, followee int64) (domain.FollowRelation, error)
}

type CachedRelationRepository struct {
	dao   dao.FollowRelationDao
	cache cache.FollowCache
	l     logger.Logger
}
