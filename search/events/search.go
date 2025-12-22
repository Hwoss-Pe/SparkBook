package events

import (
	"Webook/pkg/logger"
	"Webook/pkg/saramax"
	"Webook/search/service"
	"context"
	"github.com/IBM/sarama"
	"time"
)

const DataSyncData = "search_sync_data"

type SyncDataEventConsumer struct {
	svc    service.SyncService
	client sarama.Client
	l      logger.Logger
}

func NewSyncDataEventConsumer(svc service.SyncService, client sarama.Client, l logger.Logger) *SyncDataEventConsumer {
	return &SyncDataEventConsumer{svc: svc, client: client, l: l}

}

type SyncDataEvent struct {
	IndexName string
	DocID     string
	Data      string
}

func (a *SyncDataEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("search_sync_data",
		a.client)
	if err != nil {
		return err
	}
	go func() {
		for {
			err := cg.Consume(context.Background(), []string{DataSyncData}, saramax.NewHandler[SyncDataEvent](a.l, a.Consume))
			if err != nil {
				a.l.Error("退出了消费循环异常", logger.Error(err))
			}
		}
	}()

	return err
}

func (a *SyncDataEventConsumer) Consume(sg *sarama.ConsumerMessage,
	evt SyncDataEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return a.svc.InputAny(ctx, evt.IndexName, evt.DocID, evt.Data)
}
