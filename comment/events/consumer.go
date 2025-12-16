package events

import (
	"Webook/comment/domain"
	"Webook/comment/repository"
	"Webook/pkg/logger"
	"Webook/pkg/saramax"
	"github.com/IBM/sarama"
	"golang.org/x/net/context"
	"time"
)

type CommentEventConsumer struct {
	client sarama.Client
	l      logger.Logger
	repo   repository.CommentRepository
}

func NewCommentEventConsumer(client sarama.Client, l logger.Logger, repo repository.CommentRepository) *CommentEventConsumer {
	return &CommentEventConsumer{client: client, l: l, repo: repo}
}

func (c *CommentEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("comment", c.client)
	if err != nil {
		return err
	}
	go func() {
		err2 := cg.Consume(context.Background(), []string{TopicCommentCreate}, saramax.NewHandler[CommentCreateEvent](c.l, c.consumeCreate))
		if err2 != nil {
			c.l.Error("退出了消费循环异常", logger.Error(err2))
		}
	}()
	go func() {
		err2 := cg.Consume(context.Background(), []string{TopicCommentDelete}, saramax.NewHandler[CommentDeleteEvent](c.l, c.consumeDelete))
		if err2 != nil {
			c.l.Error("退出了消费循环异常", logger.Error(err2))
		}
	}()
	return nil
}

func (c *CommentEventConsumer) consumeCreate(_ *sarama.ConsumerMessage, evt CommentCreateEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	cm := domain.Comment{
		Biz:     evt.Biz,
		BizId:   evt.BizId,
		Content: evt.Content,
	}
	if evt.ParentId > 0 {
		cm.ParentComment = &domain.Comment{Id: evt.ParentId}
	}
	if evt.RootId > 0 {
		cm.RootComment = &domain.Comment{Id: evt.RootId}
	}
	cm.Commentator = domain.User{ID: evt.Uid}
	return c.repo.CreateComment(ctx, cm)
}

func (c *CommentEventConsumer) consumeDelete(_ *sarama.ConsumerMessage, evt CommentDeleteEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return c.repo.DeleteComment(ctx, domain.Comment{Id: evt.Id})
}
