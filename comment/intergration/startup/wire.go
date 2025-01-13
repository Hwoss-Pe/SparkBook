//go:build wireinject

package startup

import (
	grpc2 "Webook/comment/grpc"
	"Webook/comment/repository"
	"Webook/comment/repository/dao"
	"Webook/comment/service"
	"Webook/pkg/logger"
	"github.com/google/wire"
)

var serviceProviderSet = wire.NewSet(
	dao.NewGORMCommentDAO,
	repository.NewCommentRepo,
	service.NewCommentSvc,
	grpc2.NewGrpcServer,
)

var thirdProvider = wire.NewSet(
	logger.NewNoOpLogger,
	InitTestDB,
)

func InitGRPCServer() *grpc2.CommentServiceServer {
	wire.Build(thirdProvider, serviceProviderSet)
	return new(grpc2.CommentServiceServer)
}
