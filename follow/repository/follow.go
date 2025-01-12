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
	GetFollower(ctx context.Context, follower, offset, limit int64) ([]domain.FollowRelation, error)
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

func (c *CachedRelationRepository) GetFollower(ctx context.Context, followee, offset, limit int64) ([]domain.FollowRelation, error) {
	// 获取粉丝
	followerList, err := c.dao.FolloweeRelationList(ctx, followee, offset, limit)
	if err != nil {
		return nil, err
	}
	return c.genFollowRelationList(followerList), nil
}
func (c *CachedRelationRepository) GetFollowee(ctx context.Context, follower, offset, limit int64) ([]domain.FollowRelation, error) {
	followerList, err := c.dao.FollowRelationList(ctx, follower, offset, limit)
	if err != nil {
		return nil, err
	}
	return c.genFollowRelationList(followerList), nil
}

func (c *CachedRelationRepository) AddFollowRelation(ctx context.Context, f domain.FollowRelation) error {
	err := c.dao.CreateFollowRelation(ctx, c.toEntity(f))
	if err != nil {
		return err
	}
	return c.cache.Follow(ctx, f.Follower, f.Followee)
}

func (c *CachedRelationRepository) InactiveFollowRelation(ctx context.Context, follower int64, followee int64) error {
	err := c.dao.UpdateStatus(ctx, followee, follower, dao.FollowRelationStatusInactive)
	if err != nil {
		return err
	}
	return c.cache.CancelFollow(ctx, follower, followee)
}

func (c *CachedRelationRepository) GetFollowStatics(ctx context.Context, uid int64) (domain.FollowStatics, error) {
	info, err := c.cache.StaticsInfo(ctx, uid)
	//  快路径
	if err == nil {
		return info, err
	}
	// 慢路径
	info.Followers, err = c.dao.CntFollower(ctx, uid)
	if err != nil {
		return info, err
	}
	info.Followees, err = c.dao.CntFollowee(ctx, uid)
	if err != nil {
		return info, err
	}
	err = c.cache.SetStaticsInfo(ctx, uid, info)
	if err != nil {
		c.l.Error("缓存关注统计信息失败", logger.Error(err), logger.Int64("uid", uid))
	}
	return info, nil
}

func (c *CachedRelationRepository) FollowInfo(ctx context.Context, follower int64, followee int64) (domain.FollowRelation, error) {
	cd, err := c.dao.FollowerRelationDetail(ctx, follower, followee)
	if err != nil {
		return domain.FollowRelation{}, err
	}
	return c.toDomain(cd), nil
}

func NewFollowRelationRepository(dao dao.FollowRelationDao,
	cache cache.FollowCache, l logger.Logger) FollowRepository {
	return &CachedRelationRepository{
		dao:   dao,
		cache: cache,
		l:     l,
	}
}

func (c *CachedRelationRepository) toDomain(fr dao.FollowRelation) domain.FollowRelation {
	return domain.FollowRelation{
		Followee: fr.Followee,
		Follower: fr.Follower,
	}
}

func (c *CachedRelationRepository) toEntity(cd domain.FollowRelation) dao.FollowRelation {
	return dao.FollowRelation{
		Followee: cd.Followee,
		Follower: cd.Follower,
	}
}

func (c *CachedRelationRepository) genFollowRelationList(followList []dao.FollowRelation) []domain.FollowRelation {
	relations := make([]domain.FollowRelation, 0, len(followList))
	for _, cd := range followList {
		relations = append(relations, c.toDomain(cd))
	}
	return relations
}
