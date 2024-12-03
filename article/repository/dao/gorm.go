package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type GORMArticleDAO struct {
	db *gorm.DB
}

func (G *GORMArticleDAO) Insert(ctx context.Context, art Article) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (G *GORMArticleDAO) UpdateById(ctx context.Context, art Article) error {
	//TODO implement me
	panic("implement me")
}

func (G *GORMArticleDAO) GetByAuthor(ctx context.Context, author int64, offset, limit int) ([]Article, error) {
	//TODO implement me
	panic("implement me")
}

func (G *GORMArticleDAO) GetById(ctx context.Context, id int64) (Article, error) {
	//TODO implement me
	panic("implement me")
}

func (G *GORMArticleDAO) GetPubById(ctx context.Context, id int64) (PublishedArticle, error) {
	//TODO implement me
	panic("implement me")
}

func (G *GORMArticleDAO) Sync(ctx context.Context, art Article) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (G *GORMArticleDAO) SyncStatus(ctx context.Context, author, id int64, status uint8) error {
	//TODO implement me
	panic("implement me")
}

func (G *GORMArticleDAO) ListPubByUtime(ctx context.Context, utime time.Time, offset int, limit int) ([]PublishedArticle, error) {
	var res []PublishedArticle
	err := G.db.WithContext(ctx).
		Order("utime DESC").
		Where("utime < ? ", utime.UnixMilli()).
		Limit(limit).Offset(offset).Find(&res).Error
	return res, err
}
