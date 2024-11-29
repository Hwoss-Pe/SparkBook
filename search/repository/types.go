package repository

import (
	"Webook/search/domain"
	"context"
)

type UserRepository interface {
	InputUser(ctx context.Context, msg domain.User) error
	SearchUser(ctx context.Context, keywords []string) ([]domain.User, error)
}

type ArticleRepository interface {
	InputArticle(ctx context.Context, msg domain.Article) error
	SearchArticle(ctx context.Context, uid int64, keywords []string) ([]domain.Article, error)
}
