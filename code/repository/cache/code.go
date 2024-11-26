package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	//go:embed lua/set_code.lua
	luaSetCode string
	//go:embed lua/verify_code.lua
	luaVerifyCode             string
	ErrCodeSendTooMany        = errors.New("发送验证码太频繁")
	ErrUnknownForCode         = errors.New("发送验证码遇到未知错误")
	ErrCodeVerifyTooManyTimes = errors.New("验证次数太多")
)

//go:generate mockgen -source=./code.go -package=cachemocks -destination=mocks/code.mock.go CodeCache
type CodeCache interface {
	Set(ctx context.Context, biz string,
		phone string, code string) error

	Verify(ctx context.Context, biz string,
		phone string, inputCode string) (bool, error)
}

type RedisCodeCache struct {
	redis redis.Cmdable
}

func NewRedisCodeCache(cmd redis.Cmdable) CodeCache {
	return &RedisCodeCache{
		redis: cmd,
	}
}

// Set 如果该手机在该业务场景下，验证码不存在（都已经过期），那么发送
// 如果已经有一个验证码，但是发出去已经一分钟了，允许重发
// 如果已经有一个验证码，但是没有过期时间，说明有不知名错误
// 如果已经有一个验证码，但是发出去不到一分钟，不允许重发
// 验证码有效期 10 分钟
func (c *RedisCodeCache) Set(ctx context.Context,
	biz string,
	phone string,
	code string) error {
	res, err := c.redis.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		return nil
	case -1:
		//	最近发过
		return ErrCodeSendTooMany
	default:
		// 系统错误，比如说 -2，是 key 冲突
		// 其它响应码，不知道是啥鬼东西
		return ErrUnknownForCode
	}
}

// Verify 验证验证码
// 如果验证码是一致的，那么删除
// 如果验证码不一致，那么保留的
func (c *RedisCodeCache) Verify(ctx context.Context,
	biz string, phone string, inputCode string) (bool, error) {
	res, err := c.redis.Eval(ctx, luaVerifyCode, []string{c.key(biz, phone)}, inputCode).Int()
	if err != nil {
		return false, err
	}
	switch res {
	case 0:
		return true, nil
	case -1:
		//	验证次数耗尽，一般都是意味着有人在捣乱
		return false, ErrCodeVerifyTooManyTimes
	default:
		// 验证码不对
		return false, nil
	}
}

func (c *RedisCodeCache) key(biz string, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
