//go:build wireinject

package main

import (
	"Webook/bff/ioc"
	"Webook/bff/web"
	"Webook/bff/web/jwt"
	"Webook/pkg/wego"
	"github.com/google/wire"
)

func InitApp() *wego.App {
	wire.Build(
		ioc.InitLogger,
		ioc.InitRedis,
		ioc.InitEtcdClient,

		ioc.NewArticleHandler,
		web.NewCommentHandler,
		web.NewUserHandler,
		web.NewRewardHandler,
		jwt.NewRedisHandler,

		ioc.InitUserClient,
		ioc.InitIntrClient,
		ioc.InitFollowClient,
		ioc.InitRewardClient,
		ioc.InitRankingClient,
		ioc.InitCodeClient,
		ioc.InitArticleClient,
		ioc.InitCommentClient,
		ioc.InitGinServer,
		wire.Struct(new(wego.App), "WebServer"),
	)
	return new(wego.App)
}
