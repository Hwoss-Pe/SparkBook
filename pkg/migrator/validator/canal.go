package validator

import (
	"Webook/pkg/logger"
	"Webook/pkg/migrator"
	"Webook/pkg/migrator/event"
	"context"
	"errors"
	"gorm.io/gorm"
)

// CanalIncrValidator 利用canal而不是双写的方案
type CanalIncrValidator[T migrator.Entity] struct {
	baseValidator
}

func NewCanalIncrValidator[T migrator.Entity](base *gorm.DB, target *gorm.DB, direction string,
	l logger.Logger, producer event.Producer,
) *CanalIncrValidator[T] {
	return &CanalIncrValidator[T]{
		baseValidator: baseValidator{
			base:      base,
			target:    target,
			direction: direction,
			l:         l,
			producer:  producer,
		},
	}
}

// Validate 一次校验一条,校验都是根据id的
func (v *CanalIncrValidator[T]) Validate(ctx context.Context, id int64) error {
	var base T

	err := v.base.WithContext(ctx).Where("id = ?", id).First(&base).Error
	switch {
	case err == nil:
		var target T
		//如果源库有，找目的库，
		//目的库有就进行字段校验，没有就发送消息修复
		err1 := v.target.WithContext(ctx).Where("id =?", id).First(&target).Error
		switch {
		case err1 == nil:
			if !base.CompareTo(target) {
				v.notify(id, event.InconsistentEventTypeNotEqual)
			}
		case errors.Is(err1, gorm.ErrRecordNotFound):
			v.notify(id, event.InconsistentEventTypeTargetMissing)
		default:
			return err
		}
	case errors.Is(err, gorm.ErrRecordNotFound):
		var target T
		err1 := v.target.WithContext(ctx).Where("id = ?").First(&target).Error
		switch {
		case errors.Is(err1, gorm.ErrRecordNotFound):
			// 数据一致
		case err1 == nil:
			v.notify(id, event.InconsistentEventTypeBaseMissing)
		default:
			return err
		}
	default:
		return err
	}
	return nil
}
