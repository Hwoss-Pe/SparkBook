package dao

import (
	"context"
	"github.com/ecodeclub/ekit/sqlx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

var ErrWaitingSMSNotFound = gorm.ErrRecordNotFound

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
	//SELECT for UPDATE 数据库层面加锁查询
	var s AsyncSms
	err := G.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 为了避开一些偶发性的失败，只找 1 分钟前的异步短信发送
		now := time.Now().UnixMilli()
		endTime := now - time.Minute.Milliseconds()
		err := tx.Clauses(clause.Locking{
			Strength: "UPDATE",
		}).Where("utime< ? and status = ?", endTime, asyncStatusWaiting).
			First(&s).Error
		if err != nil {
			return err
		}
		// 只要更新了更新时间，根据前面的规则，就不可能被别的节点抢占了
		// 更新时间和计数
		err = tx.Model(&AsyncSms{}).Where("id = ?", s.Id).Updates(
			map[string]any{
				"retry_cnt": gorm.Expr("retry_cnt + 1"),
				"utime":     now,
			}).Error

		return err
	})
	return s, err
}

func (G *GORMAsyncSmsDAO) MarkSuccess(ctx context.Context, id int64) error {
	now := time.Now()
	return G.db.WithContext(ctx).Model(&AsyncSms{}).Where("id = ?", id).Updates(map[string]any{
		"utime":  now,
		"status": asyncStatusSuccess,
	}).Error
}

func (G *GORMAsyncSmsDAO) MarkFailed(ctx context.Context, id int64) error {
	now := time.Now()
	//	只有达到最大重试次数才会标记失败
	return G.db.WithContext(ctx).Model(&AsyncSms{}).Where("id = ? retry_cnt >= retry_max", id).Updates(map[string]any{
		"utime":  now,
		"status": asyncStatusFailed,
	}).Error
}
