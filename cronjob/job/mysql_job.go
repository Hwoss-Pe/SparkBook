package job

import (
	rankingv1 "Webook/api/proto/gen/api/proto/ranking/v1"
	"Webook/cronjob/domain"
	"Webook/cronjob/service"
	"Webook/pkg/logger"
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"time"
)

// Executor 执行器，任务执行器
type Executor interface {
	Name() string
	// Exec ctx 这个是全局控制，Executor 的实现者注意要正确处理 ctx 超时或者取消
	Exec(ctx context.Context, j domain.CronJob) error
}

// LocalFuncExecutor 调用本地方法的
type LocalFuncExecutor struct {
	funcs  map[string]func(ctx context.Context, j domain.CronJob) error
	client rankingv1.RankingServiceClient
}

func NewLocalFuncExecutor(client rankingv1.RankingServiceClient) *LocalFuncExecutor {
	return &LocalFuncExecutor{
		funcs:  map[string]func(ctx context.Context, j domain.CronJob) error{},
		client: client,
	}
}

func (l *LocalFuncExecutor) Name() string {
	return "local"
}
func (l *LocalFuncExecutor) Ranking(ctx context.Context, job domain.CronJob) error {
	//只是计算后存储到redis
	_, err := l.client.TopN(ctx, &rankingv1.TopNRequest{})
	if err != nil {
		return err
	}
	return nil
}

func (l *LocalFuncExecutor) RegisterFunc(name string, fn func(ctx context.Context, j domain.CronJob) error) {
	l.funcs[name] = fn
}

func (l *LocalFuncExecutor) Exec(ctx context.Context, j domain.CronJob) error {
	fn, ok := l.funcs[j.Name]
	if !ok {
		return fmt.Errorf("未注册本地方法 %s", j.Name)
	}
	return fn(ctx, j)
}

type Scheduler struct {
	dbTimeout time.Duration

	svc service.CronJobService

	executors map[string]Executor
	l         logger.Logger

	limiter *semaphore.Weighted
}

func NewScheduler(svc service.CronJobService, l logger.Logger) *Scheduler {
	return &Scheduler{
		svc:       svc,
		dbTimeout: time.Second,
		limiter:   semaphore.NewWeighted(100),
		l:         l,
		executors: map[string]Executor{},
	}
}

func (s *Scheduler) RegisterExecutor(exec Executor) {
	s.executors[exec.Name()] = exec
}

func (s *Scheduler) Schedule(ctx context.Context) error {
	for {
		// 放弃调度
		if ctx.Err() != nil {
			return ctx.Err()
		}
		err := s.limiter.Acquire(ctx, 1)
		if err != nil {
			return err
		}
		dbCtx, cancel := context.WithTimeout(ctx, s.dbTimeout)
		j, err := s.svc.Preempt(dbCtx)
		cancel()
		if err != nil {
			continue
		}
		exec, ok := s.executors[j.Executor]
		if !ok {
			s.l.Error("找不到执行器",
				logger.Int64("jid", j.Id),
				logger.String("executor", j.Executor))
			continue
		}
		go func() {
			defer func() {
				s.limiter.Release(1)
				// 这边要释放掉
				j.CancelFunc()
			}()
			err1 := exec.Exec(ctx, j)
			if err1 != nil {
				s.l.Error("执行任务失败",
					logger.Int64("jid", j.Id),
					logger.Error(err1))
				return
			}
			err1 = s.svc.ResetNextTime(ctx, j)
			if err1 != nil {
				s.l.Error("重置下次执行时间失败",
					logger.Int64("jid", j.Id),
					logger.Error(err1))
			}
		}()
	}
}
