package saramax

import (
	"github.com/IBM/sarama"
)

// HandlerFunc 可以直接做成函数，直接传入对应实现
type HandlerFunc func(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error

func (h HandlerFunc) Setup(session sarama.ConsumerGroupSession) error {
	//消费前的逻辑
	return nil
}

func (h HandlerFunc) Cleanup(session sarama.ConsumerGroupSession) error {
	//消费后的逻辑
	return nil
}

func (h HandlerFunc) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	return h(session, claim)
}
