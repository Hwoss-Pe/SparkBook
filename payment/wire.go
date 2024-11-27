//go:build wireinject

package main

import (
	grpc2 "Webook/payment/grpc"
	"Webook/payment/ioc"
	"Webook/payment/repository"
	"Webook/payment/repository/dao"
	"Webook/payment/web"
	"Webook/pkg/wego"
	"github.com/google/wire"
)

func InitApp() *wego.App {
	wire.Build(
		ioc.InitLogger,
		ioc.InitKafka,
		ioc.InitProducer,
		ioc.InitWechatClient,
		dao.NewPaymentGORMDAO,
		ioc.InitDB,
		repository.NewPaymentRepository,
		grpc2.NewWechatServiceServer,
		ioc.InitWechatNativeService,
		ioc.InitWechatConfig,
		ioc.InitEtcdClient,
		ioc.InitWechatNotifyHandler,
		ioc.InitGRPCServer,
		web.NewWechatHandler,
		ioc.InitGinServer,
		wire.Struct(new(wego.App), "WebServer", "GRPCServer"))
	return new(wego.App)
}
