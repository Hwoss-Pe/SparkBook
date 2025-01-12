package cache

import (
	"Webook/follow/domain"
	"golang.org/x/net/context"
)

type FollowCache interface {
	StaticsInfo(ctx context.Context, uid int64) (domain.FollowRelation, error)
	SetStaticsInfo(ctx context.Context, uid int64, statics domain.FollowStatics) error
	Follow(ctx context.Context, follow, followee int64) error
	CancelFollow(ctx context.Context, follower, followee int64) error
}
