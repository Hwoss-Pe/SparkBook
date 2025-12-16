package events

import (
	"encoding/json"
	"github.com/IBM/sarama"
)

const (
	TopicLike          = "intr_like_event"
	TopicCancelLike    = "intr_cancel_like_event"
	TopicCollect       = "intr_collect_event"
	TopicCancelCollect = "intr_cancel_collect_event"
)

type LikeEvent struct {
	Biz   string `json:"biz"`
	BizId int64  `json:"biz_id"`
	Uid   int64  `json:"uid"`
}

type CollectEvent struct {
	Biz   string `json:"biz"`
	BizId int64  `json:"biz_id"`
	Cid   int64  `json:"cid"`
	Uid   int64  `json:"uid"`
}

type Producer interface {
	ProduceLike(evt LikeEvent) error
	ProduceCancelLike(evt LikeEvent) error
	ProduceCollect(evt CollectEvent) error
	ProduceCancelCollect(evt CollectEvent) error
}

type saramaProducer struct {
	p sarama.SyncProducer
}

func NewInteractiveProducer(p sarama.SyncProducer) Producer {
	return &saramaProducer{p: p}
}

func (s *saramaProducer) send(topic string, v any) error {
	bs, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, _, err = s.p.SendMessage(&sarama.ProducerMessage{Topic: topic, Value: sarama.ByteEncoder(bs)})
	return err
}

func (s *saramaProducer) ProduceLike(evt LikeEvent) error       { return s.send(TopicLike, evt) }
func (s *saramaProducer) ProduceCancelLike(evt LikeEvent) error { return s.send(TopicCancelLike, evt) }
func (s *saramaProducer) ProduceCollect(evt CollectEvent) error { return s.send(TopicCollect, evt) }
func (s *saramaProducer) ProduceCancelCollect(evt CollectEvent) error {
	return s.send(TopicCancelCollect, evt)
}
