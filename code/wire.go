//go:build wireinject

package code

import (
	grpc2 "Webook/code/grpc"
	"Webook/code/ioc"
	"Webook/code/repository"
	"Webook/code/repository/cache"
	"Webook/code/service"
	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(
	ioc.InitRedis,
	ioc.InitEtcdClient,
	ioc.InitLogger,
)

func Init() *App {
	wire.Build(
		thirdProvider,
		ioc.InitSmsRpcClient,
		cache.NewRedisCodeCache,
		repository.NewCachedCodeRepository,
		service.NewSMSCodeService,
		grpc2.NewCodeServiceServer,
		ioc.InitGRPCxServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
