//go:build wireinject

package main

import (
	"Webook/follow/events"
	grpc2 "Webook/follow/grpc"
	"Webook/follow/ioc"
	"Webook/follow/repository"
	"Webook/follow/repository/cache"
	"Webook/follow/repository/dao"
	"Webook/follow/service"
	"github.com/google/wire"
)

var serviceProviderSet = wire.NewSet(
	dao.NewGORMFollowRelationDAO,
	repository.NewFollowRelationRepository,
	cache.NewRedisFollowCache,
	service.NewFollowRelationService,
	grpc2.NewFollowRelationServiceServer,
)

var thirdProvider = wire.NewSet(
	ioc.InitRedisClient,
	ioc.InitEtcdClient,
	ioc.InitLogger,
	ioc.InitDB,
	ioc.InitKafka,
	ioc.InitProducer,
)

func Init() *App {
	wire.Build(
		thirdProvider, serviceProviderSet, ioc.InitGRPCxServer, events.NewFollowEventConsumer, events.NewProducer, ioc.NewConsumers, wire.Struct(new(App), "*"))
	return new(App)
}
