//go:build wireinject

package main

import (
	"Webook/search/events"
	grpc2 "Webook/search/grpc"
	"Webook/search/ioc"
	"Webook/search/repository"
	"Webook/search/repository/dao"
	"Webook/search/service"
	"github.com/google/wire"
)

var serviceProviderSet = wire.NewSet(
	dao.NewUserElasticDAO,
	dao.NewArticleElasticDAO,
	dao.NewAnyESDAO,
	dao.NewTagESDAO,
	repository.NewUserRepository,
	repository.NewArticleRepository,
	repository.NewAnyRepository,
	service.NewSyncService,
	service.NewSearchService,
)

var thirdProvider = wire.NewSet(
	ioc.InitESClient,
	ioc.InitEtcdClient,
	ioc.InitLogger,
	ioc.InitKafka)

func Init() *App {
	wire.Build(
		thirdProvider,
		serviceProviderSet,
		grpc2.NewSyncServiceServer,
		grpc2.NewSearchServiceServer,
		events.NewUserConsumer,
		events.NewArticleConsumer,
		events.NewSyncDataEventConsumer,
		ioc.InitGRPCxServer,
		ioc.NewConsumers,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
