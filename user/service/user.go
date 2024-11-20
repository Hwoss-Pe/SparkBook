package service

import (
	"Webook/user/domain"
	"Webook/user/repository"
	"context"
	"go.uber.org/zap"
)

// UserService 这里注册，登录，查找，更新的基本操作
type UserService interface {
	Signup(ctx context.Context, user domain.User) error
	Login(ctx context.Context, email, password string) (domain.User, error)
	Profile(ctx context.Context, id int64) (domain.User, error)
	UpdateNonSensitiveInfo(ctx context.Context, user domain.User) error

	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	// FindOrCreateByWechat 当用户第一次扫码登录的时候判断是要注册还是登录
	FindOrCreateByWechat(ctx context.Context, info domain.WechatInfo) (domain.User, error)
}

type userService struct {
	repo   repository.UserRepository
	logger *zap.Logger
}

func NewUserService(repo repository.UserRepository, logger *zap.Logger) UserService {
	return &userService{repo: repo, logger: logger}
}

// Signup 注册的逻辑对密码进行加密后直接写进数据库，在服务层进行透
func (u *userService) Signup(ctx context.Context, user domain.User) error {
	return u.repo.Create(ctx, user)
}

func (u *userService) Login(ctx context.Context, email, password string) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userService) Profile(ctx context.Context, id int64) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userService) UpdateNonSensitiveInfo(ctx context.Context, user domain.User) error {
	//TODO implement me
	panic("implement me")
}

func (u *userService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userService) FindOrCreateByWechat(ctx context.Context, info domain.WechatInfo) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}
