package service

import (
	"Webook/account/domain"
	"Webook/account/repository"
	"context"
)

type accountService struct {
	repo repository.AccountRepository
}

func NewAccountService(repo repository.AccountRepository) AccountService {
	return &accountService{repo: repo}
}

func (a *accountService) Credit(ctx context.Context, cr domain.Credit) error {
	return a.repo.AddCredit(ctx, cr)
}
