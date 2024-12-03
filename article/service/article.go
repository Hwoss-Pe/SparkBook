package service

import (
	"Webook/article/domain"
	"Webook/article/events"
	"Webook/article/repository"
	"Webook/pkg/logger"
	"context"
	"time"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Publish(ctx context.Context, art domain.Article) (int64, error)
	Withdraw(ctx context.Context, uid, id int64) error
	PublishV1(ctx context.Context, art domain.Article) (int64, error)
	List(ctx context.Context, author int64, offset, limit int) ([]domain.Article, error)
	GetById(ctx context.Context, id int64) (domain.Article, error)
	// GetPublishedById 查找已经发表的
	GetPublishedById(ctx context.Context, id, uid int64) (domain.Article, error)
	// ListPub 根据更新时间来分页，更新时间必须小于 startTime
	ListPub(ctx context.Context, startTime time.Time, offset, limit int) ([]domain.Article, error)
}

type articleService struct {
	// 1. 在 service 这一层使用两个 repository
	//authorRepo repository.ArticleAuthorRepository
	//readerRepo repository.ArticleReaderRepository

	logger logger.Logger

	repo repository.ArticleRepository

	// 搞个异步的
	producer events.Producer
}

func (a *articleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (a *articleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (a *articleService) Withdraw(ctx context.Context, uid, id int64) error {
	//TODO implement me
	panic("implement me")
}

func (a *articleService) PublishV1(ctx context.Context, art domain.Article) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (a *articleService) List(ctx context.Context, author int64,
	offset, limit int) ([]domain.Article, error) {
	return a.repo.List(ctx, author, offset, limit)
}

func (a *articleService) GetById(ctx context.Context, id int64) (domain.Article, error) {
	return a.repo.GetById(ctx, id)
}

func (a *articleService) GetPublishedById(ctx context.Context, id, uid int64) (domain.Article, error) {
	//TODO implement me
	panic("implement me")
}

func (a *articleService) ListPub(ctx context.Context, startTime time.Time, offset, limit int) ([]domain.Article, error) {
	return a.repo.ListPub(ctx, startTime, offset, limit)
}
