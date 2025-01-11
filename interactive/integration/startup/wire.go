//go:build wireinject

package startup

import (
	grpc2 "Webook/interactive/grpc"
	"Webook/interactive/repository"
	"Webook/interactive/repository/cache"
	"Webook/interactive/repository/dao"
	"Webook/interactive/service"
	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(
	InitRedis, InitTestDB,
	InitLog,
	InitKafka,
)

func InitGRPCServer() *grpc2.InteractiveServiceServer {
	wire.Build(
		grpc2.NewInteractiveServiceServer,
		thirdProvider,
		dao.NewGORMInteractiveDAO,
		cache.NewRedisInteractiveCache,
		repository.NewCachedInteractiveRepository,
		service.NewInteractiveService,
	)
	return new(grpc2.InteractiveServiceServer)
}
