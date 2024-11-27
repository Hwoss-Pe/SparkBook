package wechat

import (
	"Webook/payment/domain"
	"Webook/payment/repository"
	"Webook/pkg/logger"
	"context"
	"errors"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"time"
)

var errUnknownTransactionState = errors.New("未知的微信事务状态")

type NativePaymentService struct {
	svc       *native.NativeApiService
	appId     string
	mchId     string
	notifyURL string
	l         logger.Logger
	repo      repository.PaymentRepository
}

func (n *NativePaymentService) GetPayment(ctx context.Context, bizTradeId string) (domain.Payment, error) {
	return n.repo.GetPayment(ctx, bizTradeId)
}

func (n *NativePaymentService) Prepay(ctx context.Context, pmt domain.Payment) (string, error) {
	//先记录订单初始状态，然后发送prepay获取二维码链接
	err := n.repo.AddPayment(ctx, pmt)
	if err != nil {
		return "", err
	}
	now := time.Now().Add(time.Minute * 30)
	resp, _, err := n.svc.Prepay(ctx, native.PrepayRequest{
		Appid:       &n.appId,
		Mchid:       &n.mchId,
		Description: &pmt.Description,
		OutTradeNo:  &pmt.BizTradeNO,
		TimeExpire:  &now,
		NotifyUrl:   &n.notifyURL,
		Amount: &native.Amount{
			Total:    &pmt.Amt.Total,
			Currency: &pmt.Amt.Currency,
		},
	})
	if err != nil {
		return "", err
	}
	return *resp.CodeUrl, nil
}
