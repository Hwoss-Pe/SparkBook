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
	cg, err := sarama.NewConsumerGroupFromClient("interactive", c.client)
	if err != nil {
		return err
	}
	go func() {
		err2 := cg.Consume(context.Background(), []string{TopicLike}, saramax.NewHandler[LikeEvent](c.l, c.consumeLike))
		if err2 != nil {
			c.l.Error("退出了消费循环异常", logger.Error(err2))
		}
	}()
	go func() {
		err2 := cg.Consume(context.Background(), []string{TopicCancelLike}, saramax.NewHandler[LikeEvent](c.l, c.consumeCancelLike))
		if err2 != nil {
			c.l.Error("退出了消费循环异常", logger.Error(err2))
		}
	}()
	return nil
}

func (c *LikeEventConsumer) consumeLike(_ *sarama.ConsumerMessage, evt LikeEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return c.repo.IncrLike(ctx, evt.Biz, evt.BizId, evt.Uid)
}

func (c *LikeEventConsumer) consumeCancelLike(_ *sarama.ConsumerMessage, evt LikeEvent) error {
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
	cg, err := sarama.NewConsumerGroupFromClient("interactive", c.client)
	if err != nil {
		return err
	}
	go func() {
		err2 := cg.Consume(context.Background(), []string{TopicCollect}, saramax.NewHandler[CollectEvent](c.l, c.consumeCollect))
		if err2 != nil {
			c.l.Error("退出了消费循环异常", logger.Error(err2))
		}
	}()
	go func() {
		err2 := cg.Consume(context.Background(), []string{TopicCancelCollect}, saramax.NewHandler[CollectEvent](c.l, c.consumeCancelCollect))
		if err2 != nil {
			c.l.Error("退出了消费循环异常", logger.Error(err2))
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
