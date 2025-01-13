package service

import (
	"Webook/cronjob/domain"
	"Webook/cronjob/repository"
	"Webook/pkg/logger"
	"golang.org/x/net/context"
	"time"
)

//go:generate mockgen -source=./cron_job.go -package=svcmocks -destination=mocks/cron_job.mock.go CronJobService
type CronJobService interface {
	Preempt(ctx context.Context) (domain.CronJob, error)
	ResetNextTime(ctx context.Context, job domain.CronJob) error
	AddJob(ctx context.Context, j domain.CronJob) error
}

type cronJobService struct {
	repo            repository.CronJobRepository
	l               logger.Logger
	refreshInterval time.Duration
}

func (c *cronJobService) Preempt(ctx context.Context) (domain.CronJob, error) {
	j, err := c.repo.Preempt(ctx)
	if err != nil {
		return domain.CronJob{}, err
	}
	ch := make(chan struct{})
	go func() {
		// 这边要启动一个 goroutine 开始续约，也就是在持续占有期间
		ticker := time.NewTicker(c.refreshInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				//一直更新utime
				c.refresh(j.Id)
			//	下面的close也会让他收到信号然后退出循环
			case <-ch:
				return
			}
		}
	}()
	j.CancelFunc = func() {
		close(ch)
		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
		defer cancelFunc()
		err = c.repo.Release(ctx, j.Id)
		if err != nil {
			c.l.Error("释放任务失败",
				logger.Error(err),
				logger.Int64("id", j.Id))
		}
	}
	return j, nil
}

func (c *cronJobService) ResetNextTime(ctx context.Context, job domain.CronJob) error {
	t := job.Next(time.Now())
	//可以判断出是否下次还能进行
	if !t.IsZero() {
		return c.repo.UpdateNextTime(ctx, job.Id, t)
	}
	return nil
}

func (c *cronJobService) AddJob(ctx context.Context, j domain.CronJob) error {
	j.NextTime = j.Next(time.Now())
	return c.repo.AddJob(ctx, j)
}

func (c *cronJobService) refresh(id int64) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := c.repo.UpdateUtime(ctx, id)
	if err != nil {
		c.l.Error("续约失败",
			logger.Int64("jid", id),
			logger.Error(err))
	}
}

func NewCronJobService(
	repo repository.CronJobRepository,
	l logger.Logger) CronJobService {
	return &cronJobService{
		repo:            repo,
		l:               l,
		refreshInterval: time.Second * 10,
	}
}
