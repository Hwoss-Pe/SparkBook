//go:build wireinject

package main

import (
	"Webook/pkg/wego"
	"Webook/user/grpc"
	"Webook/user/ioc"
	"Webook/user/repository"
	"Webook/user/repository/cache"
	"Webook/user/repository/dao"
	"Webook/user/service"
	"github.com/google/wire"
)

//需求分析，用户模块，比较简单的就是注册，登录，详情，以及查找的功能

var thirdProvider = wire.NewSet(
	ioc.InitLogger,
	ioc.InitDB,
	ioc.InitRedis,
	ioc.InitEtcdClient,
)

func Init() *wego.App {
	wire.Build(
		thirdProvider,
		cache.NewRedisUserCache,
		dao.NewGORMUserDAO,
		repository.NewCachedUserRepository,
		service.NewUserService,
		grpc.NewUserServiceServer,
		ioc.InitGRPCxServer,
		wire.Struct(new(wego.App), "GRPCServer"),
	)
	return new(wego.App)
}
