package ioc

import (
	"Webook/bff/web"
	jwt2 "Webook/bff/web/jwt"
	"Webook/bff/web/middleware"
	"Webook/pkg/ginx"
	"Webook/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"strings"
	"time"
)

func InitGinServer(l logger.Logger,
	jwtHdl jwt2.Handler,
	user *web.UserHandler,
	article *web.ArticleHandler,
	reward *web.RewardHandler) *ginx.Server {
	engine := gin.Default()
	engine.Use(corsHdl(), timeout(), middleware.NewJWTLoginMiddlewareBuilder(jwtHdl).Build())
	user.RegisterRoute(engine)
	article.RegisterRoute(engine)
	reward.RegisterRoute(engine)
	addr := viper.GetString("http.addr")
	ginx.InitCounter(prometheus.CounterOpts{
		Namespace: "daming_geektime",
		Subsystem: "webbok_bff",
		Name:      "http",
	})
	ginx.SetLogger(l)
	return &ginx.Server{
		Engine: engine,
		Addr:   addr,
	}
}
func timeout() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, ok := ctx.Request.Context().Deadline()
		if !ok {
			withTimeout, cancelFunc := context.WithTimeout(ctx.Request.Context(), time.Second*10)
			defer cancelFunc()
			ctx.Request = ctx.Request.Clone(withTimeout)
		}
		ctx.Next()
	}
}
func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowHeaders: []string{"Content-Type", "Authorization"},

		ExposeHeaders: []string{"x-jwt-token", "x-refresh-token"},

		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "yourcompany.com")
		},
		MaxAge: 12 * time.Hour,
	})
}
