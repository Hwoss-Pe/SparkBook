package repository

import (
	"Webook/pkg/logger"
	"Webook/tag/domain"
	"Webook/tag/repository/cache"
	"Webook/tag/repository/dao"
	"context"
	"github.com/ecodeclub/ekit/slice"
	"time"
)

type TagRepository interface {
	CreateTag(ctx context.Context, tag domain.Tag) (int64, error)
	BindTagToBiz(ctx context.Context, uid int64, biz string, bizId int64, tags []int64) error
	GetTags(ctx context.Context, uid int64) ([]domain.Tag, error)
	GetTagsById(ctx context.Context, ids []int64) ([]domain.Tag, error)
	GetBizTags(ctx context.Context, uid int64, biz string, bizId int64) ([]domain.Tag, error)
}
type CachedTagRepository struct {
	dao   dao.TagDAO
	cache cache.TagCache
	l     logger.Logger
}

func (c *CachedTagRepository) CreateTag(ctx context.Context, tag domain.Tag) (int64, error) {
	id, err := c.dao.CreateTag(ctx, c.toEntity(tag))
	if err != nil {
		return 0, err
	}
	err = c.cache.Append(ctx, tag.Uid, tag)
	if err != nil {
		logger.Error(err)
	}
	return id, nil
}

func (c *CachedTagRepository) BindTagToBiz(ctx context.Context, uid int64, biz string, bizId int64, tags []int64) error {
	return c.dao.CreateTagBiz(ctx, slice.Map(tags, func(idx int, src int64) dao.TagBiz {
		return dao.TagBiz{
			Tid:   src,
			BizId: bizId,
			Biz:   biz,
			Uid:   uid,
		}
	}))
}

func (c *CachedTagRepository) GetTags(ctx context.Context, uid int64) ([]domain.Tag, error) {
	res, err := c.cache.GetTags(ctx, uid)
	if err == nil {
		return res, nil
	}
	tags, err := c.dao.GetTagsByUid(ctx, uid)
	if err != nil {
		return nil, err
	}

	res = slice.Map(tags, func(idx int, src dao.Tag) domain.Tag {
		return c.toDomain(src)
	})
	err = c.cache.Append(ctx, uid, res...)
	if err != nil {
		logger.Error(err)
	}
	return res, nil
}

func (c *CachedTagRepository) GetTagsById(ctx context.Context, ids []int64) ([]domain.Tag, error) {
	tags, err := c.dao.GetTagsById(ctx, ids)
	if err != nil {
		return nil, err
	}
	return slice.Map(tags, func(idx int, src dao.Tag) domain.Tag {
		return c.toDomain(src)
	}), nil
}

func (c *CachedTagRepository) GetBizTags(ctx context.Context, uid int64, biz string, bizId int64) ([]domain.Tag, error) {
	tags, err := c.dao.GetTagsByBiz(ctx, uid, biz, bizId)
	if err != nil {
		return nil, err
	}
	return slice.Map(tags, func(idx int, src dao.Tag) domain.Tag {
		return c.toDomain(src)
	}), nil
}

func NewTagRepository(tagDAO dao.TagDAO, c cache.TagCache, l logger.Logger) TagRepository {
	return &CachedTagRepository{
		dao:   tagDAO,
		l:     l,
		cache: c,
	}
}

func NewPreTagRepository(tagDAO dao.TagDAO, c cache.TagCache, l logger.Logger) *CachedTagRepository {
	return &CachedTagRepository{
		dao:   tagDAO,
		l:     l,
		cache: c,
	}
}
func (c *CachedTagRepository) toDomain(tag dao.Tag) domain.Tag {
	return domain.Tag{
		Id:   tag.Id,
		Name: tag.Name,
		Uid:  tag.Uid,
	}
}

func (c *CachedTagRepository) toEntity(tag domain.Tag) dao.Tag {
	return dao.Tag{
		Id:   tag.Id,
		Name: tag.Name,
		Uid:  tag.Uid,
	}
}

// PreloadUserTags 预加载机制，从数据库里面取出来然后挨个放进redis
func (c *CachedTagRepository) PreloadUserTags(ctx context.Context) error {
	offset := 0
	const batch = 100
	for {
		dbCtx, cancel := context.WithTimeout(ctx, time.Second)
		tags, err := c.dao.GetTags(dbCtx, offset, batch)
		cancel()
		if err != nil {
			return err
		}
		//  goroutine 来并发
		for _, tag := range tags {
			cCtx, cancel := context.WithTimeout(ctx, time.Second)
			err = c.cache.Append(cCtx, tag.Uid, c.toDomain(tag))
			cancel()
			if err != nil {
				return err
			}
		}
		if len(tags) < batch {
			break
		}
		offset = offset + batch
	}
	return nil
}
