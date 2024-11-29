//go:build wireinject

package startup

import (
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
	dao.NewTagESDAO,
	dao.NewAnyESDAO,
	repository.NewUserRepository,
	repository.NewAnyRepository,
	repository.NewArticleRepository,
	service.NewSyncService,
	service.NewSearchService,
)

var thirdProvider = wire.NewSet(
	InitESClient,
	ioc.InitLogger)

func InitSearchServer() *grpc2.SearchServiceServer {
	wire.Build(
		thirdProvider,
		serviceProviderSet,
		grpc2.NewSearchServiceServer,
	)
	return new(grpc2.SearchServiceServer)
}

func InitSyncServer() *grpc2.SyncServiceServer {
	wire.Build(
		thirdProvider,
		serviceProviderSet,
		grpc2.NewSyncServiceServer,
	)
	return new(grpc2.SyncServiceServer)
}
