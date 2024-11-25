package dao

import (
	"context"
	"github.com/ecodeclub/ekit/sqlx"
	"gorm.io/gorm"
)

type AsyncSmsDAO interface {
	Insert(ctx context.Context, s AsyncSms) error
	GetWaitingSMS(ctx context.Context) (AsyncSms, error)
	MarkSuccess(ctx context.Context, id int64) error
	MarkFailed(ctx context.Context, id int64) error
}
type AsyncSms struct {
	Id       int64
	Config   sqlx.JsonColumn[SmsConfig]
	RetryCnt int
	RetryMax int

	Status uint8
	Ctime  int64
	Utime  int64 `gorm:"index"`
}

const (
	// 因为本身状态没有暴露出去，所以不需要在 domain 里面定义
	asyncStatusWaiting = iota
	// 失败了，并且超过了重试次数
	asyncStatusFailed
	asyncStatusSuccess
)

type SmsConfig struct {
	TplId   string
	Args    []string
	Numbers []string
}

func NewGORMAsyncSmsDAO(db *gorm.DB) AsyncSmsDAO {
	return &GORMAsyncSmsDAO{
		db: db,
	}
}

type GORMAsyncSmsDAO struct {
	db *gorm.DB
}

func (G *GORMAsyncSmsDAO) Insert(ctx context.Context, s AsyncSms) error {
	return G.db.Create(&s).Error
}

func (G *GORMAsyncSmsDAO) GetWaitingSMS(ctx context.Context) (AsyncSms, error) {
	//TODO implement me
	panic("implement me")
}

func (G *GORMAsyncSmsDAO) MarkSuccess(ctx context.Context, id int64) error {
	//TODO implement me
	panic("implement me")
}

func (G *GORMAsyncSmsDAO) MarkFailed(ctx context.Context, id int64) error {
	//TODO implement me
	panic("implement me")
}
