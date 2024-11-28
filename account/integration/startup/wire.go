//go:build wireinject

package startup

import (
	grpc2 "Webook/account/grpc"
	"Webook/account/repository"
	"Webook/account/repository/dao"
	"Webook/account/service"
	"github.com/google/wire"
)

func InitAccountService() *grpc2.AccountServiceServer {
	wire.Build(InitTestDB,
		dao.NewCreditGORMDAO,
		repository.NewAccountRepository,
		service.NewAccountService,
		grpc2.NewAccountServiceServer)
	return new(grpc2.AccountServiceServer)
}
