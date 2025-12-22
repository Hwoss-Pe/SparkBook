package repository

import "C"
import (
	"Webook/interactive/domain"
	"Webook/interactive/repository/cache"
	"Webook/interactive/repository/dao"
	"Webook/pkg/logger"
	"errors"

	"github.com/ecodeclub/ekit/slice"
	"golang.org/x/net/context"
)

//go:generate mockgen -source=./interactive.go -package=repomocks -destination=mocks/interactive.mock.go InteractiveRepository
type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context,
		biz string, bizId int64) error
	// BatchIncrReadCnt 这里调用者要保证 bizs 和 bizIds 长度一样
	BatchIncrReadCnt(ctx context.Context, bizs []string, bizIds []int64) error
	IncrLike(ctx context.Context, biz string, bizId, uid int64) error
	DecrLike(ctx context.Context, biz string, bizId, uid int64) error
	AddCollectionItem(ctx context.Context, biz string, bizId, cid int64, uid int64) error
	RemoveCollectionItem(ctx context.Context, biz string, bizId, cid int64, uid int64) error
	Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error)
	Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error)
	Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error)
	GetByIds(ctx context.Context, biz string, ids []int64) ([]domain.Interactive, error)
	GetCollectedBizIds(ctx context.Context, biz string, uid int64, offset int, limit int) ([]int64, int64, error)
	GetLikedBizIds(ctx context.Context, biz string, uid int64, offset int, limit int) ([]int64, int64, error)
}

type CachedReadCntRepository struct {
	cache cache.InteractiveCache
	dao   dao.InteractiveDAO
	l     logger.Logger
}

func NewCachedInteractiveRepository(cache cache.InteractiveCache, dao dao.InteractiveDAO, l logger.Logger) InteractiveRepository {
	return &CachedReadCntRepository{cache: cache, dao: dao, l: l}
}

func (c *CachedReadCntRepository) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	err := c.dao.IncrReadCnt(ctx, biz, bizId)
	if err != nil {
		return err
	}
	return c.cache.IncrReadCntIfPresent(ctx, biz, bizId)
}

func (c *CachedReadCntRepository) BatchIncrReadCnt(ctx context.Context, bizs []string, bizIds []int64) error {
	return c.dao.BatchIncrReadCnt(ctx, bizs, bizIds)
}

func (c *CachedReadCntRepository) IncrLike(ctx context.Context, biz string, bizId, uid int64) error {
	err := c.dao.InsertLikeInfo(ctx, biz, bizId, uid)
	if err != nil {
		return err
	}
	return c.cache.IncrLikeCntIfPresent(ctx, biz, bizId)
}

func (c *CachedReadCntRepository) DecrLike(ctx context.Context, biz string, bizId, uid int64) error {
	err := c.dao.DeleteLikeInfo(ctx, biz, bizId, uid)
	if err != nil {
		return err
	}
	return c.cache.DecrLikeCntIfPresent(ctx, biz, bizId)
}

func (c *CachedReadCntRepository) AddCollectionItem(ctx context.Context, biz string, bizId, cid int64, uid int64) error {
	err := c.dao.InsertCollectionBiz(ctx, dao.UserCollectionBiz{
		Biz:   biz,
		Cid:   cid,
		BizId: bizId,
		Uid:   uid,
	})
	if err != nil {
		return err
	}
	return c.cache.IncrCollectCntIfPresent(ctx, biz, bizId)
}

func (c *CachedReadCntRepository) RemoveCollectionItem(ctx context.Context, biz string, bizId, cid int64, uid int64) error {
	err := c.dao.DeleteCollectionBiz(ctx, biz, bizId, cid, uid)
	if err != nil {
		return err
	}
	// 尝试递减缓存中的收藏数
	if dc, ok := c.cache.(interface {
		DecrCollectCntIfPresent(ctx context.Context, biz string, bizId int64) error
	}); ok {
		_ = dc.DecrCollectCntIfPresent(ctx, biz, bizId)
	}
	return nil
}

func (c *CachedReadCntRepository) Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error) {
	intr, err := c.cache.Get(ctx, biz, bizId)
	//缓存策略，只缓存这篇文章的点赞和收藏，对于用户干了什么的操作缓存的命中不大
	if err == nil {
		return intr, nil
	}
	ie, err := c.dao.Get(ctx, biz, bizId)
	//	这里多一个双写的策略，无论找到数据还是没找到都要进行回写
	//并且回写失败其实并不关心
	if err == nil || errors.Is(err, dao.ErrRecordNotFound) {
		res := c.toDomain(ie)
		if er := c.cache.Set(ctx, biz, bizId, res); er != nil {
			c.l.Error("回写缓存失败",
				logger.Int64("bizId", bizId),
				logger.String("biz", biz),
				logger.Error(er))
		}
		return res, nil
	}
	return domain.Interactive{}, err
}

func (c *CachedReadCntRepository) Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	_, err := c.dao.GetLikeInfo(ctx, biz, id, uid)
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, dao.ErrRecordNotFound):
		return false, nil
	default:
		return false, err
	}
}

func (c *CachedReadCntRepository) Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	_, err := c.dao.GetCollectionInfo(ctx, biz, id, uid)
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, dao.ErrRecordNotFound):
		return false, nil
	default:
		return false, err
	}
}

func (c *CachedReadCntRepository) GetByIds(ctx context.Context, biz string, ids []int64) ([]domain.Interactive, error) {
	vals, err := c.dao.GetByIds(ctx, biz, ids)
	if err != nil {
		return nil, err
	}
	return slice.Map[dao.Interactive, domain.Interactive](vals,
		func(idx int, src dao.Interactive) domain.Interactive {
			return c.toDomain(src)
		}), nil
}

func (c *CachedReadCntRepository) GetCollectedBizIds(ctx context.Context, biz string, uid int64, offset int, limit int) ([]int64, int64, error) {
	return c.dao.GetCollectedBizIds(ctx, biz, uid, offset, limit)
}

func (c *CachedReadCntRepository) GetLikedBizIds(ctx context.Context, biz string, uid int64, offset int, limit int) ([]int64, int64, error) {
	return c.dao.GetLikedBizIds(ctx, biz, uid, offset, limit)
}

func (c *CachedReadCntRepository) toDomain(intr dao.Interactive) domain.Interactive {
	return domain.Interactive{
		Biz:        intr.Biz,
		BizId:      intr.BizId,
		LikeCnt:    intr.LikeCnt,
		CollectCnt: intr.CollectCnt,
		ReadCnt:    intr.ReadCnt,
	}
}
