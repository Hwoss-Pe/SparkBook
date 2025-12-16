//go:build wireinject

package main

import (
	"Webook/comment/events"
	grpc2 "Webook/comment/grpc"
	"Webook/comment/ioc"
	"Webook/comment/repository"
	"Webook/comment/repository/dao"
	"Webook/comment/service"
	"github.com/google/wire"
)

var serviceProviderSet = wire.NewSet(
	dao.NewGORMCommentDAO,
	repository.NewCommentRepo,
	service.NewCommentSvc,
	grpc2.NewGrpcServer,
)
var thirdProvider = wire.NewSet(
	ioc.InitLogger,
	ioc.InitDB,
	ioc.InitEtcdClient,
	ioc.InitKafka,
	ioc.InitProducer,
)

func Init() *App {
	wire.Build(
		thirdProvider,
		serviceProviderSet,
		ioc.InitGRPCxServer,
		events.NewProducer,
		events.NewCommentEventConsumer,
		ioc.NewConsumers,
		wire.Struct(new(App), "*"))
	return new(App)
}
