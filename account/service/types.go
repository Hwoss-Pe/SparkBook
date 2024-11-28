package service

import (
	"Webook/account/domain"
	"context"
)

type AccountService interface {
	Credit(ctx context.Context, cr domain.Credit) error
}
