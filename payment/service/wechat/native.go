package wechat

import (
	"Webook/payment/domain"
	"Webook/payment/events"
	"Webook/payment/repository"
	"Webook/pkg/logger"
	"context"
	"errors"
	"fmt"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
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
	producer  events.Producer

	// 在微信 native 里面，分别是
	// SUCCESS：支付成功
	// REFUND：转入退款
	// NOTPAY：未支付
	// CLOSED：已关闭
	// REVOKED：已撤销（付款码支付）
	// USERPAYING：用户支付中（付款码支付）
	// PAYERROR：支付失败(其他原因，如银行返回失败)
	nativeCBTypeToStatus map[string]domain.PaymentStatus
}

func NewNativePaymentService(svc *native.NativeApiService,
	repo repository.PaymentRepository,
	producer events.Producer,
	l logger.Logger,
	appid, mchid string) *NativePaymentService {
	return &NativePaymentService{
		l:     l,
		repo:  repo,
		svc:   svc,
		appId: appid,
		mchId: mchid,
		//这个是微信回调的域名
		notifyURL: "http://wechat.hwoss.com/pay/callback",
		nativeCBTypeToStatus: map[string]domain.PaymentStatus{
			"SUCCESS":  domain.PaymentStatusSuccess,
			"PAYERROR": domain.PaymentStatusFailed,
			"NOTPAY":   domain.PaymentStatusInit,
			"CLOSED":   domain.PaymentStatusFailed,
			"REVOKED":  domain.PaymentStatusFailed,
			"REFUND":   domain.PaymentStatusRefund,
		},
	}
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
func (n *NativePaymentService) FindExpiredPayment(ctx context.Context, offset, limit int, t time.Time) ([]domain.Payment, error) {
	return n.repo.FindExpiredPayment(ctx, offset, limit, t)
}
func (n *NativePaymentService) HandleCallback(ctx context.Context, txn *payments.Transaction) error {
	return n.updateByTxn(ctx, txn)
}

// SyncWechatInfo 主动去微信那边根据bizNO和商户Id查对应订单结果，这时候回返回的是新的事务id，查完更新对应的数据库
func (n *NativePaymentService) SyncWechatInfo(ctx context.Context, bizTradeNO string) error {
	txn, _, err := n.svc.QueryOrderByOutTradeNo(ctx, native.QueryOrderByOutTradeNoRequest{
		OutTradeNo: core.String(bizTradeNO),
		Mchid:      core.String(n.mchId),
	})
	if err != nil {
		return err
	}
	return n.updateByTxn(ctx, txn)
}

// 在数据库更新，并且发送到kafka如果需要对订单数据进行采集
func (n *NativePaymentService) updateByTxn(ctx context.Context, txn *payments.Transaction) error {
	status, ok := n.nativeCBTypeToStatus[*txn.TradeState]
	if !ok {
		return fmt.Errorf("%w, %s", errUnknownTransactionState, *txn.TradeState)
	}
	pmt := domain.Payment{
		BizTradeNO: *txn.OutTradeNo,
		TxnID:      *txn.TransactionId,
		Status:     status,
	}
	err := n.repo.UpdatePayment(ctx, pmt)
	if err != nil {
		return err
	}
	err1 := n.producer.ProducePaymentEvent(ctx, events.PaymentEvent{
		BizTradeNO: pmt.BizTradeNO,
		Status:     pmt.Status.AsUint8(),
	})
	if err1 != nil {
		n.l.Error("发送支付事件失败", logger.Error(err),
			logger.String("biz_trade_no", pmt.BizTradeNO))
	}
	//哪怕发送失败，但是数据库已经有数据可以接受
	return nil
}
