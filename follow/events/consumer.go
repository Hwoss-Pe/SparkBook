package events

import (
	"Webook/follow/domain"
	"Webook/follow/repository"
	"Webook/pkg/logger"
	"Webook/pkg/saramax"
	"time"

	"github.com/IBM/sarama"
	"golang.org/x/net/context"
)

type FollowEventConsumer struct {
	client sarama.Client
	l      logger.Logger
	repo   repository.FollowRepository
}

func NewFollowEventConsumer(client sarama.Client, l logger.Logger, repo repository.FollowRepository) *FollowEventConsumer {
	return &FollowEventConsumer{client: client, l: l, repo: repo}
}

func (c *FollowEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("follow", c.client)
	if err != nil {
		return err
	}
	go func() {
		err2 := cg.Consume(context.Background(), []string{TopicFollow}, saramax.NewHandler[FollowEvent](c.l, c.consumeFollow))
		if err2 != nil {
			c.l.Error("退出了消费循环异常", logger.Error(err2))
		}
	}()
	go func() {
		err2 := cg.Consume(context.Background(), []string{TopicCancelFollow}, saramax.NewHandler[FollowEvent](c.l, c.consumeCancel))
		if err2 != nil {
			c.l.Error("退出了消费循环异常", logger.Error(err2))
		}
	}()
	return nil
}

func (c *FollowEventConsumer) consumeFollow(_ *sarama.ConsumerMessage, evt FollowEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if evt.Follower == evt.Followee {
		return nil
	}
	return c.repo.AddFollowRelation(ctx, domain.FollowRelation{Follower: evt.Follower, Followee: evt.Followee})
}

func (c *FollowEventConsumer) consumeCancel(_ *sarama.ConsumerMessage, evt FollowEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if evt.Follower == evt.Followee {
		return nil
	}
	return c.repo.InactiveFollowRelation(ctx, evt.Follower, evt.Followee)
}
