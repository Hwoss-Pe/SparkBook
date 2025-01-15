package events

import (
	"github.com/IBM/sarama"
	"github.com/goccy/go-json"
)

const topicReadEvent = "article_read_event"

type ReadEvent struct {
	Aid int64
	Uid int64
}
type Producer interface {
	ProduceReadEvent(evt ReadEvent) error
}

type SaramaSyncProducer struct {
	producer sarama.SyncProducer
}

func (s *SaramaSyncProducer) ProduceReadEvent(evt ReadEvent) error {
	data, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	_, _, err = s.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topicReadEvent,
		Value: sarama.ByteEncoder(data),
	})
	return err
}

func NewSaramaSyncProducer(producer sarama.SyncProducer) Producer {
	return &SaramaSyncProducer{producer: producer}
}
