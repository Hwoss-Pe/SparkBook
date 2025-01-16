package repository

import (
	"Webook/article/domain"
	"Webook/article/repository/dao"
	"context"
)

// ArticleAuthorRepository 演示在 service 层面上分流
//
//go:generate mockgen -source=./author.go -package=repomocks -destination=mocks/article_author.mock.go ArticleAuthorRepository
type ArticleAuthorRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
}

// CachedArticleAuthorRepository 按照道理，这里也是可以搞缓存的
type CachedArticleAuthorRepository struct {
	dao dao.ArticleDAO
}

func (c *CachedArticleAuthorRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	return c.dao.Insert(ctx, c.toEntity(art))
}

func (c *CachedArticleAuthorRepository) Update(ctx context.Context, art domain.Article) error {
	return c.dao.UpdateById(ctx, c.toEntity(art))
}

func NewArticleAuthorRepository(dao dao.ArticleDAO) ArticleAuthorRepository {
	return &CachedArticleAuthorRepository{
		dao: dao,
	}
}

func (c *CachedArticleAuthorRepository) toEntity(art domain.Article) dao.Article {
	return dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	}
}
