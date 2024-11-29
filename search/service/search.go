package service

import (
	"Webook/search/domain"
	"Webook/search/repository"
	"context"
	"golang.org/x/sync/errgroup"
	"strings"
)

type SearchService interface {
	Search(ctx context.Context, uid int64, expression string) (domain.SearchResult, error)
}

type searchService struct {
	userRepo    repository.UserRepository
	articleRepo repository.ArticleRepository
}

func (s *searchService) Search(ctx context.Context, uid int64, expression string) (domain.SearchResult, error) {

	// 这边一般要对 expression 进行一些预处理,其实这里更好的是去对接算法的相关
	//没有使用 multi query 或者 multi match 之类的写法
	keywords := strings.Split(expression, " ")

	var eg errgroup.Group
	var res domain.SearchResult
	eg.Go(func() error {
		users, err := s.userRepo.SearchUser(ctx, keywords)
		res.Users = users
		return err
	})
	eg.Go(func() error {
		arts, err := s.articleRepo.SearchArticle(ctx, uid, keywords)
		res.Articles = arts
		return err
	})
	return res, eg.Wait()
}

func NewSearchService(userRepo repository.UserRepository, articleRepo repository.ArticleRepository) SearchService {
	return &searchService{userRepo: userRepo, articleRepo: articleRepo}
}
