package wechat

import (
	"Webook/payment/domain"
	"context"
)

//对于大部分的支付方式，都会有预支付的操作，其实就是创建code_url的操作

type PaymentService interface {
	PrePay(ctx context.Context, pmt domain.Payment) (string, error)
}
