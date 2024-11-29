package repository

import (
	"Webook/search/repository/dao"
	"context"
)

type AnyRepository interface {
	Input(ctx context.Context, index string, docID string, data string) error
}

type anyRepository struct {
	dao dao.AnyDAO
}

func NewAnyRepository(dao dao.AnyDAO) AnyRepository {
	return &anyRepository{dao: dao}
}

func (repo *anyRepository) Input(ctx context.Context, index string, docID string, data string) error {
	return repo.dao.Input(ctx, index, docID, data)
}
