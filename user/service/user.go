package service

import (
	"Webook/user/domain"
	"Webook/user/repository"
	"context"
	"errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserDuplicateEmail = repository.ErrUserDuplicate
var ErrInvalidUserOrPassword = errors.New("邮箱或者密码不正确")

// UserService 这里注册，登录，查找，更新的基本操作
//
//go:generate mockgen -source=./user.go -package=svcmocks -destination=mocks/user.mock.go UserService
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

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo, logger: zap.L()}
}

// Signup 注册的逻辑对密码进行加密后直接写进数据库，在服务层进行透
func (u *userService) Signup(ctx context.Context,
	user domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	return u.repo.Create(ctx, user)
}

func (u *userService) Login(ctx context.Context, email,
	password string) (domain.User, error) {
	user, err := u.repo.FindByEmail(ctx, email)
	if errors.Is(err, repository.ErrUserNotFound) {
		//	就是数据库找不到
		return domain.User{}, ErrInvalidUserOrPassword
	}
	//	如果找到的对应数据再来判断加密后密码是否是一致的

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}

	return user, nil
}

func (u *userService) Profile(ctx context.Context,
	id int64) (domain.User, error) {
	return u.repo.FindById(ctx, id)
}

func (u *userService) UpdateNonSensitiveInfo(ctx context.Context,
	user domain.User) error {
	///因为我们要修改非敏感信息的话这里对于那些信息要么屏蔽，要么就依赖下一层修改非0值
	user.Email = ""
	user.Phone = ""
	user.Password = ""
	return u.repo.Update(ctx, user)
}

func (u *userService) FindOrCreate(ctx context.Context,
	phone string) (domain.User, error) {
	user, err := u.repo.FindByPhone(ctx, phone)
	//查下是否存在，存在就直接返回，不存在就执行注册
	if !errors.Is(repository.ErrUserNotFound, err) {
		return user, err
	}
	err = u.repo.Create(ctx, domain.User{
		Phone: phone,
	})
	if err != nil && !errors.Is(repository.ErrUserNotFound, err) {
		return domain.User{}, err
	}
	return u.repo.FindByPhone(ctx, phone)
}

func (u *userService) FindOrCreateByWechat(ctx context.Context,
	info domain.WechatInfo) (domain.User, error) {

	user, err := u.repo.FindByWechat(ctx, info.OpenId)
	if !errors.Is(err, repository.ErrUserNotFound) {
		return user, err
	}
	err = u.repo.Create(ctx, domain.User{
		WechatInfo: info,
	})
	zap.L().Info("该微信用户为注册，自动注册新用户")
	return u.repo.FindByWechat(ctx, info.OpenId)
}
