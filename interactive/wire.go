package main

import (
	"Webook/interactive/events"
	grpc2 "Webook/interactive/grpc"
	"Webook/interactive/ioc"
	"Webook/interactive/repository"
	"Webook/interactive/repository/cache"
	"Webook/interactive/service"
	"Webook/sms/repository/dao"
	"github.com/google/wire"
)

var serviceProviderSet = wire.NewSet(
	dao.NewGORMAsyncSmsDAO,
	cache.NewRedisInteractiveCache,
	repository.NewCachedInteractiveRepository,
	service.NewInteractiveService,
)

var thirdProvider = wire.NewSet(
	ioc.InitSRC,
	ioc.InitDST,
	ioc.InitDoubleWritePool,
	ioc.InitBizDB,
	ioc.InitRedisClient,
	ioc.InitLogger,
	ioc.InitKafka,
	ioc.InitEtcdClient,
	ioc.InitProducer,
)

var migratorSet = wire.NewSet(
	ioc.InitMigratorWeb,
	ioc.InitFixDataConsumer,
	ioc.InitMigratorProducer,
)

func Init() *App {
	wire.Build(thirdProvider,
		serviceProviderSet,
		migratorSet,
		grpc2.NewInteractiveServiceServer,
		events.NewInteractiveReadEventConsumer,
		ioc.InitGRPCxServer,
		ioc.NewConsumers,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
