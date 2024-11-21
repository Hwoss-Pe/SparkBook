package ratelimit

import (
	//"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	//"net/http"
	"time"
)

type Builder struct {
	prefix string
	cmd    redis.Cmdable
	//规定在单位时间内能通过的请求阈值
	interval time.Duration
	rate     int
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

//func (b *Builder) Build() *gin.HandlerFunc {
//	return func(ctx *gin.Context) {
//		//limitted, err := b.limit(ctx)
//		//if err != nil {
//		//	//如果限流插件出问题，这里直接返回
//		//	ctx.AbortWithStatus(http.StatusInternalServerError)
//		//	return
//		//}
//
//	}
//}
