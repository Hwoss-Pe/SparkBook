//go:build wireinject

package main

import (
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
)

func Init() *App {
	wire.Build(
		thirdProvider,
		serviceProviderSet,
		ioc.InitGRPCxServer, wire.Struct(new(App), "*"))
	return new(App)
}
