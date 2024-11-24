package fixer

import "C"
import (
	"Webook/pkg/logger"
	"Webook/pkg/migrator"
	"Webook/pkg/migrator/event"
	"Webook/pkg/migrator/fixer"
	"Webook/pkg/saramax"
	"context"
	"errors"
	"github.com/IBM/sarama"
	"gorm.io/gorm"
	"time"
)

type Consumer[T migrator.Entity] struct {
	client sarama.Client
	l      logger.Logger
	//这里需要指针
	srcFirst *fixer.OverrideFixer[T]
	dstFirst *fixer.OverrideFixer[T]
	topic    string
}

func NewConsumer[T migrator.Entity](client sarama.Client, l logger.Logger, src *gorm.DB, dest *gorm.DB, topic string) (*Consumer[T], error) {
	srcFirst, err := fixer.NewOverrideFixer[T](src, dest)
	if err != nil {
		return nil, err
	}
	destFirst, err := fixer.NewOverrideFixer[T](dest, src)
	if err != nil {
		return nil, err
	}
	return &Consumer[T]{
		client:   client,
		l:        l,
		srcFirst: srcFirst,
		dstFirst: destFirst,
		topic:    topic,
	}, nil
}

func (c *Consumer[T]) Start() error {
	client, err := sarama.NewConsumerGroupFromClient("migrator-fix", c.client)
	if err != nil {
		return err
	}
	//异步消费
	go func() {

		err = client.Consume(context.Background(), []string{c.topic},
			//这里用的是自己包装好的，传入对应的consume实现就行，并且参数T和这里的T不是一个东西
			saramax.NewHandler(c.l, c.Consume))
		if err != nil {
			c.l.Error("退出了消费循环异常", logger.Error(err))
		}

	}()
	return err
}

func (c *Consumer[T]) Consume(msgs *sarama.ConsumerMessage, t event.InconsistentEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	switch t.Direction {
	case "SRC":
		return c.srcFirst.Fix(ctx, t.Id)
	case "DST":
		return c.dstFirst.Fix(ctx, t.Id)
	}
	return errors.New("未知的校验方向")
}
