package events

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/goccy/go-json"
)

const topicSyncUser = "sync_user_event"

type Producer interface {
	ProduceUserEvent(ctx context.Context, evt UserEvent) error
}

type SaramaSyncProducer struct {
	p sarama.SyncProducer
}

func NewSaramaSyncProducer(p sarama.SyncProducer) Producer {
	return &SaramaSyncProducer{p: p}
}

type UserEvent struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

func (s *SaramaSyncProducer) ProduceUserEvent(ctx context.Context, evt UserEvent) error {
	data, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	_, _, err = s.p.SendMessage(&sarama.ProducerMessage{Topic: topicSyncUser, Value: sarama.ByteEncoder(data)})
	return err
}
