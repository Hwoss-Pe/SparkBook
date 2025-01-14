//go:build wireinject

package main

import (
	grpc2 "Webook/ranking/grpc"
	"Webook/ranking/ioc"
	"Webook/ranking/repository"
	"Webook/ranking/repository/cache"
	"Webook/ranking/service"
	"github.com/google/wire"
)

var serviceProviderSet = wire.NewSet(
	cache.NewRankingLocalCache,
	cache.NewRedisRankingCache,
	repository.NewCachedRankingRepository,
	service.NewBatchRankingService,
)

var thirdProvider = wire.NewSet(
	ioc.InitRedis,
	ioc.InitInterActiveRpcClient,
	ioc.InitArticleRpcClient,
	ioc.InitEtcdClient,
	ioc.InitLogger,
)

func Init() *App {
	wire.Build(
		thirdProvider,
		serviceProviderSet,
		grpc2.NewRankingServiceServer,
		ioc.InitGRPCxServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
