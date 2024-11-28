package repository

import (
	"Webook/account/domain"
	"Webook/account/repository/dao"
	"context"
	"time"
)

type accountRepository struct {
	dao dao.AccountDAO
}

func NewAccountRepository(dao dao.AccountDAO) AccountRepository {
	return &accountRepository{dao: dao}
}
func (a *accountRepository) AddCredit(ctx context.Context, c domain.Credit) error {
	activities := make([]dao.AccountActivity, 0, len(c.Items))
	now := time.Now().UnixMilli()
	for _, itm := range c.Items {
		activities = append(activities, dao.AccountActivity{
			Uid:         itm.Uid,
			Biz:         c.Biz,
			BizId:       c.BizId,
			Account:     itm.Account,
			AccountType: itm.AccountType.AsUint8(),
			Amount:      itm.Amt,
			Currency:    itm.Currency,
			Ctime:       now,
			Utime:       now,
		})
	}
	return a.dao.AddActivities(ctx, activities...)
}
