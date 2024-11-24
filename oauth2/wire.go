//go:build wireinject

package main

import (
	"Webook/oauth2/grpc"
	"Webook/oauth2/ioc"
	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(
	ioc.InitLogger,
)

func Init() *App {
	wire.Build(
		thirdProvider,
		ioc.InitPrometheus,
		grpc.NewOauth2ServiceServer,
		ioc.InitEtcdClient,
		ioc.InitGRPCxServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
