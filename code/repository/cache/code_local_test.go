package cache

import (
	"context"
	"errors"
	lru "github.com/hashicorp/golang-lru"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestLocalCodeCache_Set(t *testing.T) {
	testCases := []struct {
		name string
		// 虽然叫做 mock，但是实际用不到 mock 生成的代码
		mock    func() *lru.Cache
		biz     string
		code    string
		phone   string
		wantErr error
	}{
		{
			name: "设置成功",
			mock: func() *lru.Cache {
				c, err := lru.New(10)
				require.NoError(t, err)
				// 什么也不需要做
				return c
			},
			biz:   "login",
			code:  "123456",
			phone: "152",
		},
		{
			name: "发送太频繁",
			mock: func() *lru.Cache {
				c, err := lru.New(10)
				require.NoError(t, err)
				c.Add("phone_code:login:152", codeItem{
					code: "123456",
					cnt:  3,
					// 还有九分钟多过期
					expire: time.Now().Add(time.Minute*9 + time.Second*30),
				})
				// 什么也不需要做
				return c
			},
			biz:     "login",
			code:    "123456",
			phone:   "152",
			wantErr: ErrCodeSendTooMany,
		},
		{
			name: "系统错误",
			mock: func() *lru.Cache {
				c, err := lru.New(10)
				require.NoError(t, err)
				// 随便塞了一个类型
				c.Add("phone_code:login:152", "abc")
				// 什么也不需要做
				return c
			},
			biz:     "login",
			code:    "123456",
			phone:   "152",
			wantErr: errors.New("系统错误"),
		},
		{
			name: "重发",
			mock: func() *lru.Cache {
				c, err := lru.New(10)
				require.NoError(t, err)
				c.Add("phone_code:login:152", codeItem{
					code: "123456",
					cnt:  3,
					// 还有八分钟
					expire: time.Now().Add(time.Minute * 8),
				})
				// 什么也不需要做
				return c
			},
			biz:   "login",
			code:  "123456",
			phone: "152",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := NewLocalCodeCache(tc.mock(), time.Minute*10)
			err := c.Set(context.Background(), tc.biz, tc.phone, tc.code)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestLocalCodeCache_Verify(t *testing.T) {
	testCases := []struct {
		name string
		// 虽然叫做 mock，但是实际用不到 mock 生成的代码
		mock  func() *lru.Cache
		biz   string
		code  string
		phone string

		wantErr error
		wantOk  bool
	}{
		{
			name: "验证正确",
			mock: func() *lru.Cache {
				c, err := lru.New(10)
				require.NoError(t, err)
				c.Add("phone_code:login:152", codeItem{
					code:   "123456",
					cnt:    3,
					expire: time.Now().Add(time.Minute * 8),
				})
				return c
			},
			biz:    "login",
			code:   "123456",
			phone:  "152",
			wantOk: true,
		},
		{
			name: "验证错误",
			mock: func() *lru.Cache {
				c, err := lru.New(10)
				require.NoError(t, err)
				c.Add("phone_code:login:152", codeItem{
					code:   "123456",
					cnt:    3,
					expire: time.Now().Add(time.Minute * 8),
				})
				return c
			},
			biz:   "login",
			code:  "abcde",
			phone: "152",
		},
		{
			name: "系统错误",
			mock: func() *lru.Cache {
				c, err := lru.New(10)
				require.NoError(t, err)
				// 随便塞了一个类型
				c.Add("phone_code:login:152", "abc")
				// 什么也不需要做
				return c
			},
			biz:     "login",
			code:    "123456",
			phone:   "152",
			wantErr: errors.New("系统错误"),
		},
		{
			name: "没有发验证码",
			mock: func() *lru.Cache {
				c, err := lru.New(10)
				require.NoError(t, err)
				return c
			},
			biz:     "login",
			code:    "123456",
			phone:   "152",
			wantErr: ErrKeyNotExist,
		},
		{
			name: "验证太多次数",
			mock: func() *lru.Cache {
				c, err := lru.New(10)
				require.NoError(t, err)
				c.Add("phone_code:login:152", codeItem{
					code:   "123456",
					cnt:    0,
					expire: time.Now().Add(time.Minute * 8),
				})
				return c
			},
			biz:     "login",
			code:    "abcde",
			phone:   "152",
			wantErr: ErrCodeVerifyTooManyTimes,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := NewLocalCodeCache(tc.mock(), time.Minute*10)
			ok, err := c.Verify(context.Background(), tc.biz, tc.phone, tc.code)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantOk, ok)
		})
	}
}
