package events

import (
	"Webook/interactive/repository"
	"Webook/pkg/logger"
	"Webook/pkg/saramax"
	"time"

	"github.com/IBM/sarama"
	"golang.org/x/net/context"
)

const topicReadEvent = "article_read_event"

type ReadEvent struct {
	Aid int64
	Uid int64
}

type InteractiveReadEventConsumer struct {
	client sarama.Client
	l      logger.Logger
	repo   repository.InteractiveRepository
}

func (i *InteractiveReadEventConsumer) Start() error {
	client, err := sarama.NewConsumerGroupFromClient("interactive_read", i.client)
	if err != nil {
		return err
	}
	//	开始消费
	go func() {
		for {
			err2 := client.Consume(context.Background(), []string{topicReadEvent}, saramax.NewHandler[ReadEvent](i.l, i.Consume))
			if err2 != nil {
				i.l.Error("退出了消费循环异常", logger.Error(err))
				time.Sleep(time.Second * 5)
			} else {
				i.l.Info("消费循环正常退出，可能是发生了Rebalance，正在重试", logger.String("topic", topicReadEvent))
			}
		}
	}()
	return err
}

func (i *InteractiveReadEventConsumer) StartBatch() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive",
		i.client)
	if err != nil {
		return err
	}
	go func() {
		for {
			err := cg.Consume(context.Background(),
				[]string{topicReadEvent},
				saramax.NewBatchHandler[ReadEvent](i.l, i.BatchConsume))
			if err != nil {
				i.l.Error("退出了消费循环异常", logger.Error(err))
				time.Sleep(time.Second * 5)
			} else {
				i.l.Info("消费循环正常退出，可能是发生了Rebalance，正在重试", logger.String("topic", topicReadEvent))
			}
		}
	}()
	return err
}

func (i *InteractiveReadEventConsumer) BatchConsume(msgs []*sarama.ConsumerMessage,
	evts []ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	bizs := make([]string, 0, len(msgs))
	ids := make([]int64, 0, len(msgs))
	for _, evt := range evts {
		//这里的批量消费那边是bizId,在这传了uid是否合适？，Aid实际对应的才是bizId
		bizs = append(bizs, "article")
		ids = append(ids, evt.Aid)
	}
	return i.repo.BatchIncrReadCnt(ctx, bizs, ids)
}

func (i *InteractiveReadEventConsumer) Consume(msg *sarama.ConsumerMessage, t ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := i.repo.IncrReadCnt(ctx, "article", t.Aid)
	return err
}

func NewInteractiveReadEventConsumer(client sarama.Client, l logger.Logger, repo repository.InteractiveRepository) *InteractiveReadEventConsumer {
	return &InteractiveReadEventConsumer{client: client, l: l, repo: repo}
}
