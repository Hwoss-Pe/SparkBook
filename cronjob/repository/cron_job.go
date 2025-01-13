package repository

import (
	"Webook/cronjob/domain"
	"Webook/cronjob/repository/dao"
	"golang.org/x/net/context"
	"time"
)

//go:generate mockgen -source=./cron_job.go -package=repomocks -destination=mocks/cron_job.mock.go CronJobRepository
type CronJobRepository interface {
	Preempt(ctx context.Context) (domain.CronJob, error)
	UpdateNextTime(ctx context.Context, id int64, t time.Time) error
	UpdateUtime(ctx context.Context, id int64) error
	Release(ctx context.Context, id int64) error
	AddJob(ctx context.Context, j domain.CronJob) error
}
type PreemptCronJobRepository struct {
	dao dao.JobDAO
}

func (p *PreemptCronJobRepository) Preempt(ctx context.Context) (domain.CronJob, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PreemptCronJobRepository) UpdateNextTime(ctx context.Context, id int64, t time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (p *PreemptCronJobRepository) UpdateUtime(ctx context.Context, id int64) error {
	//TODO implement me
	panic("implement me")
}

func (p *PreemptCronJobRepository) Release(ctx context.Context, id int64) error {
	//TODO implement me
	panic("implement me")
}

func (p *PreemptCronJobRepository) AddJob(ctx context.Context, j domain.CronJob) error {
	//TODO implement me
	panic("implement me")
}

func NewPreemptCronJobRepository(dao dao.JobDAO) CronJobRepository {
	return &PreemptCronJobRepository{dao: dao}
}
