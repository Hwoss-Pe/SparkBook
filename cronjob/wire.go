//go:build wireinject

package main

import (
	grpc2 "Webook/cronjob/grpc"
	"Webook/cronjob/ioc"
	"Webook/cronjob/repository"
	"Webook/cronjob/repository/dao"
	"Webook/cronjob/service"
	"Webook/pkg/grpcx"
	"github.com/google/wire"
)

var serviceProviderSet = wire.NewSet(
	dao.NewGORMJobDAO,
	repository.NewPreemptCronJobRepository,
	service.NewCronJobService)

var thirdProvider = wire.NewSet(
	ioc.InitDB,
	ioc.InitLogger,
)

func Init() *App {
	wire.Build(
		thirdProvider,
		serviceProviderSet,
		grpc2.NewCronJobServiceServer,
		ioc.InitGRPCxServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}

type App struct {
	server *grpcx.Server
}
