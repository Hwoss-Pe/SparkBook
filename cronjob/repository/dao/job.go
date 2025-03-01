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
	db := G.db.WithContext(ctx)
	for {
		// 每一个循环都重新计算 time.Now，因为之前可能已经花了一些时间了
		now := time.Now().UnixMilli()
		var j Job
		//	扫描数据库里面下一执行时间小于当前的时间就进行抢占
		err := db.Where("next_time <= ? and status = ?", now, jobStatusWaiting).First(&j).Error
		if err != nil {
			// 数据库有问题
			return Job{}, err
		}
		// 然后要开始抢占更改它的utime
		// 这里利用 utime 来执行 CAS 操作
		res := db.Model(&Job{}).Where("id = ? and version = ? ", j.Id, j.Version).Updates(map[string]any{
			"utime":   now,
			"version": j.Version + 1,
			"status":  jobStatusRunning,
		})
		if res.Error != nil {
			return Job{}, err
		}
		//抢占成功
		if res.RowsAffected == 1 {
			return j, nil
		}
		//	没抢到就进行下一个循环
	}
}

func (G *GORMJobDAO) UpdateNextTime(ctx context.Context, id int64, t time.Time) error {
	return G.db.WithContext(ctx).Model(&Job{}).
		Where("id=?", id).Updates(map[string]any{
		"utime":     time.Now().UnixMilli(),
		"next_time": t.UnixMilli(),
	}).Error
}

func (G *GORMJobDAO) UpdateUtime(ctx context.Context, id int64) error {
	return G.db.WithContext(ctx).Model(&Job{}).
		Where("id=?", id).Updates(map[string]any{
		"utime": time.Now().UnixMilli(),
	}).Error
}

func (G *GORMJobDAO) Release(ctx context.Context, id int64) error {
	return G.db.WithContext(ctx).Model(&Job{}).
		Where("id = ?", id).Updates(map[string]any{
		"status": jobStatusEnd,
		"utime":  time.Now().UnixMilli(),
	}).Error
}

func (G *GORMJobDAO) Insert(ctx context.Context, j Job) error {
	now := time.Now().UnixMilli()
	j.Ctime = now
	j.Utime = now
	return G.db.WithContext(ctx).Create(&j).Error
}

func NewGORMJobDAO(db *gorm.DB) JobDAO {
	return &GORMJobDAO{db: db}
}
