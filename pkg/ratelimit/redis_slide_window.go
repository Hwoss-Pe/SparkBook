package ratelimit

import (
	"context"
	_ "embed"
	"github.com/redis/go-redis/v9"
	"time"
)

//go:embed lua/slide_window.lua
var luaScript string

type RedisSlidingWindowLimiter struct {
	cmd redis.Cmdable
	//也是单位时间内允许多少通过
	interval time.Duration
	rate     int
}

func (r *RedisSlidingWindowLimiter) Limit(ctx context.Context, key string) (bool, error) {

	return r.cmd.Eval(ctx, luaScript, []string{key}, r.interval.Milliseconds(), r.rate, time.Now().UnixMilli()).Bool()
}

func NewRedisSlidingWindowLimiter(cmd redis.Cmdable, interval time.Duration, rate int) *RedisSlidingWindowLimiter {
	return &RedisSlidingWindowLimiter{cmd: cmd, interval: interval, rate: rate}
}
