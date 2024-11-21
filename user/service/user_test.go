package service

import (
	"Webook/user/domain"
	"Webook/user/repository"
	repomocks "Webook/user/repository/mocks"
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"testing"
	"time"
)

// 测试这边只进行一个例子
func TestUserService_Login(t *testing.T) {
	//	固定使用一个事件
	ctime := time.Now()
	testcase := []struct {
		name string
		mock func(ctrl *gomock.Controller) repository.UserRepository
		//	输入
		ctx      context.Context
		email    string
		password string

		//输出
		wantErr  error
		wantUser domain.User
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				mockUserRepository := repomocks.NewMockUserRepository(ctrl)
				mockUserRepository.EXPECT().FindByEmail(context.Background(), "123@qq.com").Return(
					domain.User{
						Id:       123,
						Email:    "123@qq.com",
						Password: "$2a$10$s51GBcU20dkNUVTpUAQqpe6febjXkRYvhEwa5OkN5rU6rw2KTbNUi",
						Phone:    "15261890000",
						Ctime:    ctime,
					}, nil)
				return mockUserRepository
			},
			ctx:      context.Background(),
			email:    "123@qq.com",
			password: "hello#world123",
			wantUser: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Password: "$2a$10$s51GBcU20dkNUVTpUAQqpe6febjXkRYvhEwa5OkN5rU6rw2KTbNUi",
				Phone:    "15261890000",
				Ctime:    ctime,
			},
			wantErr: nil,
		},
		{
			name: "用户未找到",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				mockUserRepository := repomocks.NewMockUserRepository(ctrl)
				mockUserRepository.EXPECT().FindByEmail(context.Background(), "123@qq.com").Return(
					domain.User{}, repository.ErrUserNotFound)
				return mockUserRepository
			},
			ctx:      context.Background(),
			email:    "123@qq.com",
			password: "hello#world123",
			wantErr:  ErrInvalidUserOrPassword,
		},
		{
			name: "密码错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				mockUserRepository := repomocks.NewMockUserRepository(ctrl)
				mockUserRepository.EXPECT().FindByEmail(context.Background(), "123@qq.com").Return(
					domain.User{
						Id:       123,
						Email:    "123@qq.com",
						Password: "$2a$10$s51GBcU20dkNUVTpUAQqpe6febjXkRYvhEwa5OkN5rU6rw2KTbNUi",
						Phone:    "15261890000",
						Ctime:    ctime,
					}, nil)
				return mockUserRepository
			},
			ctx:   context.Background(),
			email: "123@qq.com",
			//这里的密码是错误的
			password: "hello#world",
			wantErr:  ErrInvalidUserOrPassword,
		},
	}

	for _, tc := range testcase {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			repo := tc.mock(controller)
			svc := NewUserService(repo)
			user, err := svc.Login(tc.ctx, tc.email, tc.password)
			assert.Equal(t, err, tc.wantErr)
			assert.Equal(t, user, tc.wantUser)
		})
	}
}

func TestPasswordEncrypt(t *testing.T) {
	pwd := []byte("hello#world123")
	// 加密
	encrypted, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	// 比较
	println(string(encrypted))
	err = bcrypt.CompareHashAndPassword(encrypted, pwd)
	require.NoError(t, err)
}

func Test(t *testing.T) {

	code := "500"
	i, _ := strconv.ParseInt(code, 10, 64)
	fmt.Println(i)
}
