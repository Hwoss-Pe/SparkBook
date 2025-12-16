package events

import (
	"encoding/json"
	"github.com/IBM/sarama"
)

const (
	TopicCommentCreate = "comment_create_event"
	TopicCommentDelete = "comment_delete_event"
)

type CommentCreateEvent struct {
	Biz      string `json:"biz"`
	BizId    int64  `json:"biz_id"`
	Uid      int64  `json:"uid"`
	Content  string `json:"content"`
	ParentId int64  `json:"parent_id"`
	RootId   int64  `json:"root_id"`
}

type CommentDeleteEvent struct {
	Id int64 `json:"id"`
}

type Producer interface {
	ProduceCreate(evt CommentCreateEvent) error
	ProduceDelete(evt CommentDeleteEvent) error
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

func (s *saramaProducer) ProduceCreate(evt CommentCreateEvent) error {
	return s.send(TopicCommentCreate, evt)
}
func (s *saramaProducer) ProduceDelete(evt CommentDeleteEvent) error {
	return s.send(TopicCommentDelete, evt)
}
