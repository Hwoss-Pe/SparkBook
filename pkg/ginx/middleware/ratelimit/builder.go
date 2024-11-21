package ratelimit

import (
	_ "embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"time"
)

type Builder struct {
	prefix string
	cmd    redis.Cmdable
	//规定在单位时间内能通过的请求阈值
	interval time.Duration
	rate     int
	//l logger.Logger 防止耦合
}

func NewBuilder(cmd redis.Cmdable, interval time.Duration, rate int) *Builder {
	return &Builder{
		prefix:   "ip-limiter",
		cmd:      cmd,
		interval: interval,
		rate:     rate}
}
func (b *Builder) Prefix(prefix string) *Builder {
	b.prefix = prefix
	return b
}

func (b *Builder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limited, err := b.Limit(ctx)
		if err != nil {
			log.Println(err)
			//如果限流插件出问题，这里直接返回
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if limited {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		//走下一个插件
		ctx.Next()
	}
}

//go:embed slide_window.lua
var luaScript string

// Limit 参数顺序，时间段，阈值，当前时间做score
func (b *Builder) Limit(ctx *gin.Context) (bool, error) {
	key := fmt.Sprintf("%s:%s", b.prefix, ctx.ClientIP())
	return b.cmd.Eval(ctx, luaScript, []string{key}, b.interval.Milliseconds(),
		b.rate, time.Now().UnixMilli()).Bool()
}
