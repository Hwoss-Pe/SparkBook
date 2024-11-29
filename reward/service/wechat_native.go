package service

import (
	accountv1 "Webook/api/proto/gen/api/proto/account/v1"
	pmtv1 "Webook/api/proto/gen/api/proto/payment/v1"
	"Webook/pkg/logger"
	"Webook/reward/domain"
	"Webook/reward/repository"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type WechatNativeRewardService struct {
	client pmtv1.WechatPaymentServiceClient
	repo   repository.RewardRepository
	l      logger.Logger
	acli   accountv1.AccountServiceClient
}

// PreReward 先查缓存，没有创建，调用支付服务，然后写缓存
func (w *WechatNativeRewardService) PreReward(ctx context.Context, r domain.Reward) (domain.CodeURL, error) {
	//决定引入缓存来缓存codeURL，因此先查缓存
	codeUrl, err := w.repo.GetCachedCodeURL(ctx, r)
	if err == nil {
		return codeUrl, nil
	}
	r.Status = domain.RewardStatusInit
	//	创建预支付
	rid, err := w.repo.CreateReward(ctx, r)
	if err != nil {
		return domain.CodeURL{}, err
	}
	//	调用支付服务获取url
	response, err := w.client.NativePrePay(ctx, &pmtv1.PrePayRequest{
		// 想办法拼接出来一个 biz_trade_id
		BizTradeNo:  fmt.Sprintf("reward-%d", rid),
		Description: fmt.Sprintf("打赏-%s", r.Target.BizName),
		Amt: &pmtv1.Amount{
			Total:    r.Amt,
			Currency: "CNY",
		},
	})
	if err != nil {
		return domain.CodeURL{}, err
	}
	cu := domain.CodeURL{
		Rid: rid,
		URL: response.CodeUrl,
	}
	//	然后缓存
	err1 := w.repo.CachedCodeURL(ctx, cu, r)
	if err1 != nil {
		w.l.Error("缓存二维码失败", logger.Error(err1))
	}
	return cu, err
}

// GetReward 走一个快慢路径，快路径就是查自己服务维护状态
func (w *WechatNativeRewardService) GetReward(ctx context.Context, rid, uid int64) (domain.Reward, error) {
	// 快路径
	r, err := w.repo.GetReward(ctx, rid)
	if err != nil {
		return domain.Reward{}, err
	}
	if r.Uid != uid {
		// 说明是非法查询
		return domain.Reward{}, errors.New("查询的打赏记录和打赏人对不上")
	}
	//必须判断这个支付操作是否已经是完结才能返回
	if r.Completed() {
		return r, nil
	}
	//	当然慢路径就是你还没支付我就亲自调用支付服务查对应的状态
	response, err := w.client.GetPayment(ctx, &pmtv1.GetPaymentRequest{
		BizTradeNo: w.bizTradeNO(r.Id),
	})
	if err != nil {
		// 这边我们直接返回从数据库查询的数据
		w.l.Error("慢路径查询支付结果失败",
			logger.Int64("rid", r.Id), logger.Error(err))
		return r, nil
	}
	//	查完后自己更新本地数据库
	switch response.Status {
	case pmtv1.PaymentStatus_PaymentStatusFailed:
		r.Status = domain.RewardStatusFailed
	case pmtv1.PaymentStatus_PaymentStatusInit:
		r.Status = domain.RewardStatusInit
	case pmtv1.PaymentStatus_PaymentStatusSuccess:
		r.Status = domain.RewardStatusPayed
	case pmtv1.PaymentStatus_PaymentStatusRefund:
		// 理论上来说不可能出现这个，直接设置为失败
		r.Status = domain.RewardStatusFailed
	}
	err = w.repo.UpdateStatus(ctx, rid, r.Status)
	if err != nil {
		w.l.Error("更新本地打赏状态失败",
			logger.Int64("rid", r.Id), logger.Error(err))
		return r, nil
	}
	return r, nil
}

// UpdateReward 更新本地表然后进行对已经支付的进行入账,对于抽成在这里计算
func (w *WechatNativeRewardService) UpdateReward(ctx context.Context, bizTradeNO string, status domain.RewardStatus) error {
	rid := w.toRid(bizTradeNO)
	err := w.repo.UpdateStatus(ctx, rid, status)
	if err != nil {
		return err
	}
	//	更新
	if status == domain.RewardStatusPayed {
		r, err := w.repo.GetReward(ctx, rid)
		if err != nil {
			return err
		}
		//	这里模拟一个金额的抽成
		weAmt := int64(float64(r.Amt) * 0.1)
		_, err = w.acli.Credit(ctx, &accountv1.CreditRequest{
			Biz:   "reward",
			BizId: r.Uid,
			Items: []*accountv1.CreditItem{
				{
					AccountType: accountv1.AccountType_AccountTypeSystem,
					Amt:         weAmt,
					Currency:    "CNY"},
				{
					Account:     r.Uid,
					Uid:         r.Uid,
					AccountType: accountv1.AccountType_AccountTypeReward,
					Amt:         r.Amt - weAmt,
					Currency:    "CNY",
				},
			},
		})
		if err != nil {
			w.l.Error("入账失败了，快来修数据啊！！！",
				logger.String("biz_trade_no", bizTradeNO),
				logger.Error(err))
			// 做好监控和告警，这里
			return err
		}
	}
	return nil
}
func (w *WechatNativeRewardService) bizTradeNO(rid int64) string {
	return fmt.Sprintf("reward-%d", rid)
}
func (w *WechatNativeRewardService) toRid(tradeNO string) int64 {
	ridStr := strings.Split(tradeNO, "-")
	val, _ := strconv.ParseInt(ridStr[1], 10, 64)
	return val
}

func NewWechatNativeRewardService(
	client pmtv1.WechatPaymentServiceClient,
	repo repository.RewardRepository,
	l logger.Logger,
	acli accountv1.AccountServiceClient,
) RewardService {
	return &WechatNativeRewardService{client: client, repo: repo, l: l, acli: acli}
}
