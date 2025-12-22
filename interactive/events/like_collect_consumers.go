package events

import (
	"Webook/interactive/repository"
	"Webook/pkg/logger"
	"Webook/pkg/saramax"
	"time"

	"github.com/IBM/sarama"
	"golang.org/x/net/context"
)

type LikeEventConsumer struct {
	client       sarama.Client
	l            logger.Logger
	repo         repository.InteractiveRepository
	resolveOwner func(ctx context.Context, biz string, bizId int64) (int64, error)
}

func NewLikeEventConsumer(client sarama.Client, l logger.Logger, repo repository.InteractiveRepository, resolveOwner func(ctx context.Context, biz string, bizId int64) (int64, error)) *LikeEventConsumer {
	return &LikeEventConsumer{client: client, l: l, repo: repo, resolveOwner: resolveOwner}
}

func (c *LikeEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive_like", c.client)
	if err != nil {
		return err
	}
	go func() {
		for {
			err2 := cg.Consume(context.Background(), []string{TopicLike}, saramax.NewHandler[LikeEvent](c.l, c.consumeLike))
			if err2 != nil {
				c.l.Error("退出了消费循环异常", logger.Error(err2))
				time.Sleep(time.Second * 5)
			} else {
				c.l.Info("消费循环正常退出，可能是发生了Rebalance，正在重试", logger.String("topic", TopicLike))
			}
		}
	}()

	cgCancel, err := sarama.NewConsumerGroupFromClient("interactive_cancel_like", c.client)
	if err != nil {
		return err
	}
	go func() {
		for {
			err2 := cgCancel.Consume(context.Background(), []string{TopicCancelLike}, saramax.NewHandler[LikeEvent](c.l, c.consumeCancelLike))
			if err2 != nil {
				c.l.Error("退出了消费循环异常", logger.Error(err2))
				time.Sleep(time.Second * 5)
			} else {
				c.l.Info("消费循环正常退出，可能是发生了Rebalance，正在重试", logger.String("topic", TopicCancelLike))
			}
		}
	}()
	return nil
}

func (c *LikeEventConsumer) consumeLike(_ *sarama.ConsumerMessage, evt LikeEvent) error {
	c.l.Info("开始消费点赞事件", logger.Int64("uid", evt.Uid), logger.Int64("bizId", evt.BizId), logger.String("biz", evt.Biz))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return c.repo.IncrLike(ctx, evt.Biz, evt.BizId, evt.Uid)
}

func (c *LikeEventConsumer) consumeCancelLike(_ *sarama.ConsumerMessage, evt LikeEvent) error {
	c.l.Info("开始消费取消点赞事件", logger.Int64("uid", evt.Uid), logger.Int64("bizId", evt.BizId), logger.String("biz", evt.Biz))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return c.repo.DecrLike(ctx, evt.Biz, evt.BizId, evt.Uid)
}

type CollectEventConsumer struct {
	client       sarama.Client
	l            logger.Logger
	repo         repository.InteractiveRepository
	resolveOwner func(ctx context.Context, biz string, bizId int64) (int64, error)
}

func NewCollectEventConsumer(client sarama.Client, l logger.Logger, repo repository.InteractiveRepository, resolveOwner func(ctx context.Context, biz string, bizId int64) (int64, error)) *CollectEventConsumer {
	return &CollectEventConsumer{client: client, l: l, repo: repo, resolveOwner: resolveOwner}
}

func (c *CollectEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive_collect", c.client)
	if err != nil {
		return err
	}
	go func() {
		for {
			err2 := cg.Consume(context.Background(), []string{TopicCollect}, saramax.NewHandler[CollectEvent](c.l, c.consumeCollect))
			if err2 != nil {
				c.l.Error("退出了消费循环异常", logger.Error(err2))
				time.Sleep(time.Second * 5)
			} else {
				c.l.Info("消费循环正常退出，可能是发生了Rebalance，正在重试", logger.String("topic", TopicCollect))
			}
		}
	}()

	cgCancel, err := sarama.NewConsumerGroupFromClient("interactive_cancel_collect", c.client)
	if err != nil {
		return err
	}
	go func() {
		for {
			err2 := cgCancel.Consume(context.Background(), []string{TopicCancelCollect}, saramax.NewHandler[CollectEvent](c.l, c.consumeCancelCollect))
			if err2 != nil {
				c.l.Error("退出了消费循环异常", logger.Error(err2))
				time.Sleep(time.Second * 5)
			} else {
				c.l.Info("消费循环正常退出，可能是发生了Rebalance，正在重试", logger.String("topic", TopicCancelCollect))
			}
		}
	}()
	return nil
}

func (c *CollectEventConsumer) consumeCollect(_ *sarama.ConsumerMessage, evt CollectEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if c.resolveOwner != nil {
		owner, err := c.resolveOwner(ctx, evt.Biz, evt.BizId)
		if err == nil && owner == evt.Uid {
			return nil
		}
	}
	return c.repo.AddCollectionItem(ctx, evt.Biz, evt.BizId, evt.Cid, evt.Uid)
}

func (c *CollectEventConsumer) consumeCancelCollect(_ *sarama.ConsumerMessage, evt CollectEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if c.resolveOwner != nil {
		owner, err := c.resolveOwner(ctx, evt.Biz, evt.BizId)
		if err == nil && owner == evt.Uid {
			return nil
		}
	}
	return c.repo.RemoveCollectionItem(ctx, evt.Biz, evt.BizId, evt.Cid, evt.Uid)
}
