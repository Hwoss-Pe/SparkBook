//go:build wireinject

package main

import (
	"Webook/pkg/wego"
	grpc2 "Webook/tag/grpc"
	"Webook/tag/ioc"
	"Webook/tag/repository"
	"Webook/tag/repository/cache"
	"Webook/tag/repository/dao"
	"Webook/tag/service"
	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(
	ioc.InitRedis,
	ioc.InitLogger,
	ioc.InitDB,
	ioc.InitEtcdClient,
)

func Init() *wego.App {
	wire.Build(
		thirdProvider,
		cache.NewRedisTagCache,
		dao.NewGORMTagDAO,
		repository.NewTagRepository,
		ioc.InitKafka,
		ioc.InitProducer,
		service.NewTagService,
		grpc2.NewTagServiceServer,
		ioc.InitGRPCxServer,
		wire.Struct(new(wego.App), "GRPCServer"),
	)
	return new(wego.App)
}
