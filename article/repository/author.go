package repository

import (
	userv1 "Webook/api/proto/gen/api/proto/user/v1"
	"Webook/article/domain"
	"Webook/article/repository/dao"
	"context"
)

// AuthorRepository 封装user的client用于获取用户信息
type AuthorRepository interface {
	// FindAuthor id为文章id
	FindAuthor(ctx context.Context, id int64) (domain.Author, error)
}

type GrpcAuthorRepository struct {
	client userv1.UsersServiceClient
	dao    dao.ArticleDAO
}

func (g *GrpcAuthorRepository) FindAuthor(ctx context.Context, id int64) (domain.Author, error) {
	art, err := g.dao.GetPubById(ctx, id)
	if err != nil {
		return domain.Author{}, nil
	}
	u, err := g.client.Profile(ctx, &userv1.ProfileRequest{
		Id: art.AuthorId,
	})
	if err != nil {
		return domain.Author{}, err
	}
	return domain.Author{
		Id:   u.User.Id,
		Name: u.User.Nickname,
	}, nil
}

func NewGrpcAuthorRepository(articleDao dao.ArticleDAO, client userv1.UsersServiceClient) AuthorRepository {
	return &GrpcAuthorRepository{
		client: client,
		dao:    articleDao,
	}
}
