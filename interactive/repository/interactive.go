package repository

import (
	"Webook/interactive/domain"
	"golang.org/x/net/context"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context,
		biz string, bizId int64) error
	// BatchIncrReadCnt 这里调用者要保证 bizs 和 bizIds 长度一样
	BatchIncrReadCnt(ctx context.Context, bizs []string, bizIds []int64) error
	IncrLike(ctx context.Context, biz string, bizId, uid int64) error
	DecrLike(ctx context.Context, biz string, bizId, uid int64) error
	AddCollectionItem(ctx context.Context, biz string, bizId, cid int64, uid int64) error
	Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error)
	Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error)
	Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error)
	GetByIds(ctx context.Context, biz string, ids []int64) ([]domain.Interactive, error)
}
