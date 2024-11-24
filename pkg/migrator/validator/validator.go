package validator

import (
	"Webook/pkg/logger"
	"Webook/pkg/migrator"
	"Webook/pkg/migrator/event"
	"context"
	"errors"
	"github.com/ecodeclub/ekit/slice"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"time"
)

type Validator[T migrator.Entity] struct {
	baseValidator
	baseSize int
	utime    int64
	// 如果没有数据了，就睡眠
	// 如果不是正数，那么就说明直接返回，结束这一次的循环
	sleepInterval time.Duration
}

func NewValidator[T migrator.Entity](
	base *gorm.DB,
	target *gorm.DB,
	direction string,
	l logger.Logger,
	producer event.Producer) *Validator[T] {
	return &Validator[T]{baseValidator: baseValidator{
		base:      base,
		target:    target,
		direction: direction,
		l:         l,
		producer:  producer,
		// 默认是全量校验，并且数据没了就结束
	}, baseSize: 100, sleepInterval: 0}
}
func (v *Validator[T]) Utime(utime int64) *Validator[T] {
	v.utime = utime
	return v
}
func (v *Validator[T]) SleepInterval(i time.Duration) *Validator[T] {
	v.sleepInterval = i
	return v
}

// Validate 执行校验，目的库和源库的调换
// 第一次校验校验后可能漏掉的就是target里面有但是base没有的，这里才需要第二次校验
func (v *Validator[T]) Validate(ctx context.Context) error {
	var eg errgroup.Group
	eg.Go(func() error {
		return v.baseToTarget(ctx)
	})
	eg.Go(func() error {
		return v.TargetToBase(ctx)
	})
	return eg.Wait()
}
func (v *Validator[T]) baseToTarget(ctx context.Context) error {
	offset := 0
	for {
		var src T
		dbCtx, cancel := context.WithTimeout(ctx, time.Second)
		err := v.base.WithContext(dbCtx).
			Order("id").
			//这个字段来判断是不是进行增量校验的，如果是会传进对应的时间
			Where("utime >= ?", v.utime).
			Offset(offset).
			First(&src).Error
		cancel()
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			//	判断是不是已经没数据了，只有增量才会sleep
			if v.sleepInterval <= 0 {
				return nil
			}
			time.Sleep(v.sleepInterval)
			continue
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			//退出循环
			return nil
		case err == nil:
			v.dstDiff(ctx, src)
		default:
			v.l.Error("src => dst 查询源表失败", logger.Error(err))
		}
		offset++
	}
}
func (v *Validator[T]) dstDiff(ctx context.Context, src T) {
	var target T
	dbCtx, cancel := context.WithTimeout(ctx, time.Second)
	//在目的库里面找对应的id
	err := v.target.WithContext(dbCtx).
		Where("id=?", src.Id()).First(&target).Error
	cancel()
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		//目的库没有通知修复
		v.notify(src.Id(), event.InconsistentEventTypeTargetMissing)
	case err == nil:
		// 查询到了数据
		equal := src.CompareTo(target)
		if !equal {
			v.notify(src.Id(), event.InconsistentEventTypeNotEqual)
		}
	default:
		v.l.Error("src => dst 查询目标表失败", logger.Error(err))
	}
}

// TargetToBase target里面有但是base没有的，只需要从base里面取出来进行判断
func (v *Validator[T]) TargetToBase(ctx context.Context) error {
	offset := 0
	for {
		var ts []T
		dbCtx, cancel := context.WithTimeout(ctx, time.Second)
		err := v.target.WithContext(dbCtx).
			//model这里主要是指定表名，new一下更加保健
			Model(new(T)).Select("id").
			Limit(v.baseSize).Offset(offset).
			Find(&ts).Error
		cancel()
		switch {
		case err == nil:
			v.srcMissingRecords(ctx, ts)
		case errors.Is(err, gorm.ErrRecordNotFound):
			if v.sleepInterval > 0 {
				time.Sleep(v.sleepInterval)
				// 在 sleep 的时候。不需要调整偏移量，增量才需要
				continue
			}
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			return nil //
		default:
			v.l.Error("dst => src 查询目标表失败", logger.Error(err))
		}
		if len(ts) < v.baseSize {
			//没有数据的时候就退出，这里是补修复不需要在这里进行阻塞
			return nil
		}
		offset += v.baseSize
	}
}

func (v *Validator[T]) srcMissingRecords(ctx context.Context, ts []T) {
	//	获取target的id,然后根据id去src里面查一下做一个对比，找出差集进行删除
	var ids []int64
	for _, t := range ts {
		ids = append(ids, t.Id())
	}
	var srcs []T
	dbCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	//这里
	err := v.base.WithContext(dbCtx).Select("id").
		//这里另一种思路是not in ，得到ids后需要判断这个长度，到底是全部一样还是存在部分一样
		Where("id  IN?", ids).
		Find(&srcs).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		// 说明 ids 全部没有
		v.notifySrcMissing(ts)
	case err == nil:
		missing := slice.DiffSetFunc(ts, srcs, func(src, dst T) bool {
			return src.Id() == dst.Id()
		})
		v.notifySrcMissing(missing)
	default:
		v.l.Error("dst => src 查询源表失败", logger.Error(err))
	}
}
func (v *Validator[T]) notifySrcMissing(ts []T) {
	for _, t := range ts {
		v.notify(t.Id(), event.InconsistentEventTypeBaseMissing)
	}
}
