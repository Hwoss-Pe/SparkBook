package grpc2

import (
	"Webook/account/domain"
	"Webook/account/service"
	"Webook/api/proto/gen/api/proto/account/v1"
	"context"
	"github.com/ecodeclub/ekit/slice"
	"google.golang.org/grpc"
)

type AccountServiceServer struct {
	service service.AccountService
	accountv1.UnimplementedAccountServiceServer
}

func NewAccountServiceServer(svc service.AccountService) *AccountServiceServer {
	return &AccountServiceServer{service: svc}
}
func (a *AccountServiceServer) Register(server *grpc.Server) {
	accountv1.RegisterAccountServiceServer(server, a)
}

func (a *AccountServiceServer) Credit(ctx context.Context,
	req *accountv1.CreditRequest) (*accountv1.CreditResponse, error) {
	err := a.service.Credit(ctx, a.toDomain(req))
	return &accountv1.CreditResponse{}, err
}

func (a *AccountServiceServer) toDomain(c *accountv1.CreditRequest) domain.Credit {
	return domain.Credit{
		Biz:   c.Biz,
		BizId: c.BizId,
		Items: slice.Map(c.Items, func(idx int, src *accountv1.CreditItem) domain.CreditItem {
			return a.itemToDomain(src)
		}),
	}
}

func (a *AccountServiceServer) itemToDomain(c *accountv1.CreditItem) domain.CreditItem {
	return domain.CreditItem{
		Account: c.Account,
		Amt:     c.Amt,
		Uid:     c.Uid,
		// 两者取值都是一样的， 偷个懒，直接转
		AccountType: domain.AccountType(c.AccountType),
		Currency:    c.Currency,
	}
}
