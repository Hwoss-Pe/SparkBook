package repository

import (
	"Webook/article/domain"
	"context"
)

//go:generate mockgen -source=./reader.go -package=repomocks -destination=mocks/article_reader.mock.go ArticleReaderRepository
type ArticleReaderRepository interface {
	Save(ctx context.Context, art domain.Article) error
}
