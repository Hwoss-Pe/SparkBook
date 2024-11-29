package events

import (
	"Webook/pkg/logger"
	"Webook/pkg/saramax"
	"Webook/reward/domain"
	"Webook/reward/service"
	"context"
	"github.com/IBM/sarama"
	"strings"
	"time"
)

type PaymentEvent struct {
	//譬如打赏服务的就是reward-%d,在调用支付服务的时候会把bizNO传过去
	BizTradeNO string
	Status     uint8
}
type PaymentEventConsumer struct {
	client sarama.Client
	l      logger.Logger
	svc    service.RewardService
}

func (p *PaymentEventConsumer) Start() error {

	cg, err := sarama.NewConsumerGroupFromClient("reward",
		p.client)
	if err != nil {
		return err
	}
	//	异步消费
	go func() {
		err := cg.Consume(context.Background(), []string{"payment_events"}, saramax.NewHandler[PaymentEvent](p.l, p.Consume))
		if err != nil {
			p.l.Error("退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}
func (p *PaymentEventConsumer) Consume(
	msg *sarama.ConsumerMessage,
	evt PaymentEvent) error {
	if !strings.HasPrefix(evt.BizTradeNO, "reward") {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	return p.svc.UpdateReward(ctx, evt.BizTradeNO, evt.ToDomainStatus())
}
func (p PaymentEvent) ToDomainStatus() domain.RewardStatus {
	// 	PaymentStatusInit
	//	PaymentStatusSuccess
	//	PaymentStatusFailed
	//	PaymentStatusRefund
	switch p.Status {
	// 这里不能引用 payment 里面的定义，只能手写
	case 1:
		return domain.RewardStatusInit
	case 2:
		return domain.RewardStatusPayed
	case 3, 4:
		return domain.RewardStatusFailed
	default:
		return domain.RewardStatusUnknown
	}
}
