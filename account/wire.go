//go:build wireinject

package main

import (
	grpc2 "Webook/account/grpc"
	"Webook/account/ioc"
	"Webook/account/repository"
	"Webook/account/repository/dao"
	"Webook/account/service"
	"Webook/pkg/wego"
	"github.com/google/wire"
)

func Init() *wego.App {
	wire.Build(
		ioc.InitDB,
		ioc.InitLogger,
		ioc.InitEtcdClient,
		ioc.InitGRPCxServer,
		dao.NewCreditGORMDAO,
		repository.NewAccountRepository,
		service.NewAccountService,
		grpc2.NewAccountServiceServer,
		wire.Struct(new(wego.App), "GRPCServer"))
	return new(wego.App)
}
