package dao

import (
	"context"
	"time"
)

//go:generate mockgen -source=./types.go -package=artdaomocks -destination=mocks/article.mock.go ArticleDAO
type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, art Article) error
	GetByAuthor(ctx context.Context, author int64, offset, limit int) ([]Article, error)
	GetById(ctx context.Context, id int64) (Article, error)
	GetPubById(ctx context.Context, id int64) (PublishedArticle, error)
	Sync(ctx context.Context, art Article) (int64, error)
	SyncStatus(ctx context.Context, author, id int64, status uint8) error
	ListPubByUtime(ctx context.Context, utime time.Time, offset int, limit int) ([]PublishedArticle, error)
}
