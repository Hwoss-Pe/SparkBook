package repository

import (
	"Webook/user/domain"
	"Webook/user/repository/cache"
	"Webook/user/repository/dao"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	// Update 更新数据，只有非 0 值才会更新
	Update(ctx context.Context, u domain.User) error
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindById(ctx context.Context, id int64) (domain.User, error)
	// FindByWechat 暂时可以认为按照 openId来查询
	// 将来可能需要按照 unionId 来查询
	FindByWechat(ctx context.Context, openId string) (domain.User, error)
}

type CachedUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}
