package service

import (
	"Webook/search/domain"
	"context"
)

type SearchService interface {
	Search(ctx context.Context, uid int64, expression string) (domain.SearchResult, error)
}
