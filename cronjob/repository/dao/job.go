package dao

import (
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"time"
)

var ErrNoMoreJob = gorm.ErrRecordNotFound

const (
	// 等待被调度，意思就是没有人正在调度
	jobStatusWaiting = iota
	// 已经被 goroutine 抢占了
	jobStatusRunning
	// 不再需要调度了，比如说被终止了，或者被删除了。
	// 这里没有严格区分这两种情况的必要性
	jobStatusEnd
)

type JobDAO interface {
	Preempt(ctx context.Context) (Job, error)
	UpdateNextTime(ctx context.Context, id int64, t time.Time) error
	UpdateUtime(ctx context.Context, id int64) error
	Release(ctx context.Context, id int64) error
	Insert(ctx context.Context, j Job) error
}

type Job struct {
	Id         int64  `gorm:"primaryKey,autoIncrement"`
	Name       string `gorm:"type=varchar(256);unique"`
	Executor   string
	Cfg        string
	Expression string
	Version    int64
	NextTime   int64 `gorm:"index"`
	Status     int   //标识当前使用状态是否抢占
	Ctime      int64
	Utime      int64
}

type GORMJobDAO struct {
	db *gorm.DB
}

func (G *GORMJobDAO) Preempt(ctx context.Context) (Job, error) {
	//TODO implement me
	panic("implement me")
}

func (G *GORMJobDAO) UpdateNextTime(ctx context.Context, id int64, t time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (G *GORMJobDAO) UpdateUtime(ctx context.Context, id int64) error {
	//TODO implement me
	panic("implement me")
}

func (G *GORMJobDAO) Release(ctx context.Context, id int64) error {
	//TODO implement me
	panic("implement me")
}

func (G *GORMJobDAO) Insert(ctx context.Context, j Job) error {
	//TODO implement me
	panic("implement me")
}

func NewGORMJobDAO(db *gorm.DB) JobDAO {
	return &GORMJobDAO{db: db}
}
