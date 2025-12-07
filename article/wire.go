//go:build wireinject

package main

import (
	"Webook/article/events"
	grpc2 "Webook/article/grpc"
	"Webook/article/ioc"
	"Webook/article/repository"
	"Webook/article/repository/cache"
	"Webook/article/repository/dao"
	"Webook/article/service"
	"Webook/pkg/wego"
	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(
	ioc.InitRedis,
	ioc.InitLogger,
	ioc.InitUserRpcClient,
	ioc.InitProducer,
	ioc.InitEtcdClient,
	ioc.InitDB,
)

func Init() *wego.App {
	wire.Build(
		thirdProvider,
		events.NewSaramaSyncProducer,
		cache.NewRedisArticleCache,
		dao.NewGORMArticleDAO,
		dao.NewGORMArticleReaderDAO,
		repository.NewArticleRepository,
		repository.NewArticleAuthorRepository,
		repository.NewCachedArticleReaderRepository,
		repository.NewGrpcAuthorRepository,
		service.NewArticleService,
		grpc2.NewArticleServiceServer,
		ioc.InitGRPCxServer,
		wire.Struct(new(wego.App), "GRPCServer"),
	)
	return new(wego.App)
}
