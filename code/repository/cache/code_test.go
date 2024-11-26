package cache

import (
	"Webook/code/repository/cache/redismocks"
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
	"testing"
)

func TestRedisCodeCache_Set(t *testing.T) {
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) redis.Cmdable

		// 输入
		ctx   context.Context
		biz   string
		phone string
		code  string

		// 预期输出
		wantErr error
	}{
		{
			name: "设置成功",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				mockRes := redis.NewCmdResult(int64(0), nil)
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode,
					gomock.Any(), gomock.Any()).
					Return(mockRes)
				return cmd
			},
			ctx:   context.Background(),
			phone: "15212345678",
			code:  "123456",
		},
		{
			name: "发送太频繁",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				mockRes := redis.NewCmdResult(int64(-1), nil)
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode,
					gomock.Any(), gomock.Any()).
					Return(mockRes)
				return cmd
			},
			ctx:     context.Background(),
			phone:   "15212345678",
			code:    "123456",
			wantErr: ErrCodeSendTooMany,
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				mockRes := redis.NewCmdResult(int64(-2), nil)
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode,
					gomock.Any(), gomock.Any()).
					Return(mockRes)
				return cmd
			},
			ctx:     context.Background(),
			phone:   "15212345678",
			code:    "123456",
			wantErr: ErrUnknownForCode,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := NewRedisCodeCache(tc.mock(ctrl))
			err := c.Set(tc.ctx, tc.biz, tc.phone, tc.code)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
