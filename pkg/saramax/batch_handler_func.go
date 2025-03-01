package saramax

import (
	"Webook/pkg/logger"
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"time"
)

type BatchHandler[T any] struct {
	l  logger.Logger
	fn func(msgs []*sarama.ConsumerMessage, t []T) error
}

func NewBatchHandler[T any](l logger.Logger, fn func(msgs []*sarama.ConsumerMessage, t []T) error) *BatchHandler[T] {
	return &BatchHandler[T]{l: l, fn: fn}
}

func (b *BatchHandler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (b *BatchHandler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim 这里走的批量提交异步消费，可以解决消息积压的问题
func (b *BatchHandler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	const batchSize = 10
	messages := claim.Messages()
	for {
		//一次处理十条数据
		msgs := make([]*sarama.ConsumerMessage, 0, batchSize)
		ts := make([]T, 0, batchSize)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		//	判断是否超时
		done := false
		for i := 0; i < batchSize && !done; i++ {
			select {
			case <-ctx.Done():
				//如果是超时就直接退出这个循环，或者没有消息可以消费了
				done = true
			case msg, ok := <-messages:
				if !ok {
					//如果取出信息出现错误可能是kafka问题，这里直接断开
					cancel()
					return nil
				}
				msgs = append(msgs, msg)
				var t T
				err := json.Unmarshal(msg.Value, &t)
				if err != nil {
					// 消息格式都不对，没啥好处理的
					b.l.Error("反序列化消息体失败",
						logger.String("topic", msg.Topic),
						logger.Int32("partition", msg.Partition),
						logger.Int64("offset", msg.Offset),
						logger.Error(err))
					// 不中断，继续下一个
					session.MarkMessage(msg, "")
					continue
				}
				ts = append(ts, t)
			}
		}
		err := b.fn(msgs, ts)
		//	批量提交
		if err == nil {
			for _, msg := range msgs {
				session.MarkMessage(msg, "")
			}
		} else {
			//	重试或者发送到死信队列也可以的
			const maxRetries = 3
			var retryCount int64
			var err error
			for retryCount < maxRetries {
				err = b.fn(msgs, ts)
				if err == nil {
					// 成功，标记消息
					for _, msg := range msgs {
						session.MarkMessage(msg, "")
					}
					break
				}
				retryCount++
				b.l.Error("批量提交失败，正在重试", logger.Int64("retryCount", retryCount), logger.Error(err))
				time.Sleep(time.Second * time.Duration(retryCount))
			}

			if err != nil {
				b.l.Error("批量提交失败，重试已超出最大次数", logger.Error(err))
			}

		}
		cancel()
	}
}
