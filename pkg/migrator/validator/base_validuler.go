package validator

import (
	"Webook/pkg/logger"
	"Webook/pkg/migrator/event"
	"context"
	"gorm.io/gorm"
	"time"
)

type baseValidator struct {
	base   *gorm.DB
	target *gorm.DB
	// 这边需要告知，是以 SRC 为准，还是以 DST 为准
	// 修复数据需要知道
	direction string
	l         logger.Logger
	producer  event.Producer
}

// 上报数据
func (b *baseValidator) notify(id int64, typ string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	evt := event.InconsistentEvent{
		Direction: b.direction,
		Id:        id,
		Type:      typ,
	}
	err := b.producer.ProduceInconsistentEvent(ctx, evt)
	if err != nil {
		b.l.Error("发送消息失败", logger.Error(err),
			logger.Field{Key: "event", Value: evt})
	}
}
