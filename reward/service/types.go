package service

import (
	"Webook/reward/domain"
	"context"
)

//go:generate mockgen -source=./types.go -destination=mocks/reward.mock.go -package=svcmocks RewardService
type RewardService interface {
	// PreReward 准备打赏，
	//    也可以直接理解为对标到创建一个打赏的订单
	// 因为目前   只支持微信扫码支付，所以实际上直接把接口定义成这个样子就可以了
	PreReward(ctx context.Context, r domain.Reward) (domain.CodeURL, error)
	GetReward(ctx context.Context, rid, uid int64) (domain.Reward, error)
	UpdateReward(ctx context.Context, bizTradeNO string, status domain.RewardStatus) error
}
