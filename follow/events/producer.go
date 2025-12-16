package events

import (
	"encoding/json"
	"github.com/IBM/sarama"
)

const (
	TopicFollow       = "follow_event"
	TopicCancelFollow = "follow_cancel_event"
)

type FollowEvent struct {
	Follower int64 `json:"follower"`
	Followee int64 `json:"followee"`
}

type Producer interface {
	ProduceFollow(evt FollowEvent) error
	ProduceCancelFollow(evt FollowEvent) error
}

type saramaProducer struct{ p sarama.SyncProducer }

func NewProducer(p sarama.SyncProducer) Producer { return &saramaProducer{p: p} }

func (s *saramaProducer) send(topic string, v any) error {
	bs, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, _, err = s.p.SendMessage(&sarama.ProducerMessage{Topic: topic, Value: sarama.ByteEncoder(bs)})
	return err
}

func (s *saramaProducer) ProduceFollow(evt FollowEvent) error { return s.send(TopicFollow, evt) }
func (s *saramaProducer) ProduceCancelFollow(evt FollowEvent) error {
	return s.send(TopicCancelFollow, evt)
}
