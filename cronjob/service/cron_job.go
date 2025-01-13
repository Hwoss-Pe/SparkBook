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
	//TODO implement me
	panic("implement me")
}

func (c *cronJobService) ResetNextTime(ctx context.Context, job domain.CronJob) error {
	//TODO implement me
	panic("implement me")
}

func (c *cronJobService) AddJob(ctx context.Context, j domain.CronJob) error {
	//TODO implement me
	panic("implement me")
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
