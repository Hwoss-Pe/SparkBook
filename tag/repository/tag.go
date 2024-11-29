package repository

import (
	"Webook/pkg/logger"
	"Webook/tag/domain"
	"Webook/tag/repository/cache"
	"Webook/tag/repository/dao"
	"context"
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
	//TODO implement me
	panic("implement me")
}

func (c *CachedTagRepository) BindTagToBiz(ctx context.Context, uid int64, biz string, bizId int64, tags []int64) error {
	//TODO implement me
	panic("implement me")
}

func (c *CachedTagRepository) GetTags(ctx context.Context, uid int64) ([]domain.Tag, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CachedTagRepository) GetTagsById(ctx context.Context, ids []int64) ([]domain.Tag, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CachedTagRepository) GetBizTags(ctx context.Context, uid int64, biz string, bizId int64) ([]domain.Tag, error) {
	//TODO implement me
	panic("implement me")
}

func NewTagRepository(tagDAO dao.TagDAO, c cache.TagCache, l logger.Logger) TagRepository {
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
