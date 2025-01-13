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
	j, err := p.dao.Preempt(ctx)
	if err != nil {
		return domain.CronJob{}, err
	}
	return p.toDomain(j), nil
}

func (p *PreemptCronJobRepository) UpdateNextTime(ctx context.Context, id int64, t time.Time) error {
	return p.dao.UpdateNextTime(ctx, id, t)
}

func (p *PreemptCronJobRepository) UpdateUtime(ctx context.Context, id int64) error {
	return p.dao.UpdateUtime(ctx, id)
}

func (p *PreemptCronJobRepository) Release(ctx context.Context, id int64) error {
	return p.dao.Release(ctx, id)
}

func (p *PreemptCronJobRepository) AddJob(ctx context.Context, j domain.CronJob) error {
	return p.dao.Insert(ctx, p.toEntity(j))
}

func NewPreemptCronJobRepository(dao dao.JobDAO) CronJobRepository {
	return &PreemptCronJobRepository{dao: dao}
}

func (p *PreemptCronJobRepository) toEntity(j domain.CronJob) dao.Job {
	return dao.Job{
		Id:         j.Id,
		Name:       j.Name,
		Expression: j.Expression,
		Cfg:        j.Cfg,
		Executor:   j.Executor,
		NextTime:   j.NextTime.UnixMilli(),
	}
}

func (p *PreemptCronJobRepository) toDomain(j dao.Job) domain.CronJob {
	return domain.CronJob{
		Id:         j.Id,
		Name:       j.Name,
		Expression: j.Expression,
		Cfg:        j.Cfg,
		Executor:   j.Executor,
		NextTime:   time.UnixMilli(j.NextTime),
	}
}
