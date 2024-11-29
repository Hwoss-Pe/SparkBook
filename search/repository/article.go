package repository

import (
	"Webook/search/domain"
	"Webook/search/repository/dao"
	"context"
	"github.com/ecodeclub/ekit/slice"
)

type articleRepository struct {
	dao  dao.ArticleDAO
	tags dao.TagDAO
}

func (a *articleRepository) SearchArticle(ctx context.Context,
	uid int64,
	keywords []string) ([]domain.Article, error) {
	ids, err := a.tags.Search(ctx, uid, "article", keywords)
	if err != nil {
		return nil, err
	}
	arts, err := a.dao.Search(ctx, ids, keywords)
	if err != nil {
		return nil, err
	}
	return slice.Map(arts, func(idx int, src dao.Article) domain.Article {
		return domain.Article{
			Id:      src.Id,
			Title:   src.Title,
			Status:  src.Status,
			Content: src.Content,
			Tags:    src.Tags,
		}
	}), nil
}

func (a *articleRepository) InputArticle(ctx context.Context, msg domain.Article) error {
	return a.dao.InputArticle(ctx, dao.Article{
		Id:      msg.Id,
		Title:   msg.Title,
		Status:  msg.Status,
		Content: msg.Content,
	})
}

func NewArticleRepository(d dao.ArticleDAO, td dao.TagDAO) ArticleRepository {
	return &articleRepository{
		dao:  d,
		tags: td,
	}
}
