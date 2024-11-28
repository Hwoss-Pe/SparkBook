package repository

import (
	"Webook/account/domain"
	"context"
)

type AccountRepository interface {
	AddCredit(ctx context.Context, c domain.Credit) error
}
