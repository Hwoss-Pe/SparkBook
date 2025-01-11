package dao

import (
	"context"
	"errors"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"gorm.io/gorm"
)

var errUnknownPattern = errors.New("未知的双写 pattern")

type DoubleWriteDAO struct {
	src InteractiveDAO
	dst InteractiveDAO
	// 双写的模式
	pattern *atomicx.Value[string]
}

const (
	patternSrcOnly  = "src_only"
	patternSrcFirst = "src_first"
	patternDstFirst = "dst_first"
	patternDstOnly  = "dst_only"
)

func NewDoubleWriteDAO(
	src *gorm.DB,
	dst *gorm.DB) *DoubleWriteDAO {
	return &DoubleWriteDAO{
		src: NewGORMInteractiveDAO(src),
		dst: NewGORMInteractiveDAO(dst),
		// 默认就是 src 优先
		pattern: atomicx.NewValueOf(patternSrcOnly),
	}
}

func (d *DoubleWriteDAO) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	switch d.pattern.Load() {
	case patternSrcOnly:
		return d.src.IncrReadCnt(ctx, biz, bizId)
	case patternSrcFirst:
		err := d.src.IncrReadCnt(ctx, biz, bizId)
		if err == nil {
			// 如果 src 出错了，不需要考虑 dst 了
			err1 := d.dst.IncrReadCnt(ctx, biz, bizId)
			if err1 != nil {
				// 这里只能记录一下日志，做好监控
			}
		}
		return err
	case patternDstFirst:
		err := d.dst.IncrReadCnt(ctx, biz, bizId)
		if err == nil {
			// 如果 dst 出错了，不需要考虑 src 了
			err1 := d.src.IncrReadCnt(ctx, biz, bizId)
			if err1 != nil {
				// 这里只能记录一下日志，做好监控
			}
		}
		return err
	case patternDstOnly:
		return d.dst.IncrReadCnt(ctx, biz, bizId)
	default:
		return errUnknownPattern
	}
}

func (d *DoubleWriteDAO) InsertLikeInfo(ctx context.Context, biz string, bizId, uid int64) error {
	//TODO implement me
	panic("implement me")
}

func (d *DoubleWriteDAO) GetLikeInfo(ctx context.Context, biz string, bizId, uid int64) (UserLikeBiz, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DoubleWriteDAO) DeleteLikeInfo(ctx context.Context, biz string, bizId, uid int64) error {
	//TODO implement me
	panic("implement me")
}

func (d *DoubleWriteDAO) Get(ctx context.Context, biz string, bizId int64) (Interactive, error) {
	switch d.pattern.Load() {
	case patternSrcOnly, patternSrcFirst:
		return d.src.Get(ctx, biz, bizId)
	case patternDstFirst, patternDstOnly:
		return d.dst.Get(ctx, biz, bizId)
	default:
		return Interactive{}, errUnknownPattern
	}
}

func (d *DoubleWriteDAO) InsertCollectionBiz(ctx context.Context, cb UserCollectionBiz) error {
	//TODO implement me
	panic("implement me")
}

func (d *DoubleWriteDAO) GetCollectionInfo(ctx context.Context, biz string, bizId, uid int64) (UserCollectionBiz, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DoubleWriteDAO) BatchIncrReadCnt(ctx context.Context, bizs []string, ids []int64) error {
	//TODO implement me
	panic("implement me")
}

func (d *DoubleWriteDAO) GetByIds(ctx context.Context, biz string, ids []int64) ([]Interactive, error) {
	//TODO implement me
	panic("implement me")
}
