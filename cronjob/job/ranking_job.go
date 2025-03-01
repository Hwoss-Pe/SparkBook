package job

import (
	rankingv1 "Webook/api/proto/gen/api/proto/ranking/v1"
	"Webook/pkg/logger"
	"context"
	"github.com/google/uuid"
	rlock "github.com/gotomicro/redis-lock"
	"github.com/hashicorp/go-multierror"
	"github.com/redis/go-redis/v9"
	"sync"
	"sync/atomic"
	"time"
)

type RankingJobV1 struct {
	rankingClient rankingv1.RankingServiceClient
	l             logger.Logger
	timeout       time.Duration
	client        *rlock.Client
	key           string

	localLock *sync.Mutex
	lock      *rlock.Lock

	load *atomic.Int32

	nodeID         string
	redisClient    redis.Cmdable
	rankingLoadKey string
	closeSignal    chan struct{}
	loadTicker     *time.Ticker
}

func NewRankingJobV1(
	svc rankingv1.RankingServiceClient,
	l logger.Logger,
	client *rlock.Client,
	timeout time.Duration,
	redisClient redis.Cmdable,
	loadInterval time.Duration,
) *RankingJobV1 {
	res := &RankingJobV1{rankingClient: svc,
		key:            "job:ranking",
		l:              l,
		client:         client,
		localLock:      &sync.Mutex{},
		timeout:        timeout,
		nodeID:         uuid.New().String(),
		redisClient:    redisClient,
		rankingLoadKey: "ranking_job_nodes_load",
		load:           &atomic.Int32{},
		closeSignal:    make(chan struct{}),
		loadTicker:     time.NewTicker(loadInterval),
	}
	// 开启
	res.loadCycle()
	return res
}

func (r *RankingJobV1) Name() string {
	return "ranking"
}

func (r *RankingJobV1) Run() error {
	r.localLock.Lock()
	lock := r.lock
	r.localLock.Unlock()
	if lock == nil {
		// 是不是负载最低的，如果是，就尝试获取分布式锁
		// 抢分布式锁
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
		defer cancel()
		lock, err := r.client.Lock(ctx, r.key, r.timeout,
			&rlock.FixIntervalRetry{
				Interval: time.Millisecond * 100,
				Max:      3,
				// 重试的超时
			}, time.Second)
		if err != nil {
			r.l.Warn("获取分布式锁失败", logger.Error(err))
			return nil
		}
		r.l.Debug(r.nodeID + "获得了分布式锁 ")
		r.lock = lock
		go func() {
			er := lock.AutoRefresh(r.timeout/2, r.timeout)
			if er != nil {
				// 续约失败了
				// 没办法中断当下正在调度的热榜计算（如果有）
				r.localLock.Lock()
				r.lock = nil
				r.localLock.Unlock()
			}
		}()
	}
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	_, err := r.rankingClient.TopN(ctx, &rankingv1.TopNRequest{})
	return err
}

func (r *RankingJobV1) loadCycle() {
	go func() {
		for range r.loadTicker.C {
			// 上报负载
			r.reportLoad()
			r.releaseLockIfNeed()
		}
	}()
}

func (r *RankingJobV1) releaseLockIfNeed() {
	// 检测自己是不是负载最低，如果不是，那么就直接释放分布式锁。
	r.localLock.Lock()
	lock := r.lock
	r.localLock.Unlock()
	if lock != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		// 最低负载的
		res, err := r.redisClient.ZPopMin(ctx, r.rankingLoadKey).Result()
		if err != nil {
			// 记录日志
			return
		}
		head := res[0]
		if head.Member.(string) != r.nodeID {
			// 不是自己，释放锁
			r.l.Debug(r.nodeID+" 不是负载最低的节点，释放分布式锁",
				logger.Field{Key: "head", Value: head})
			r.localLock.Lock()
			r.lock = nil
			r.localLock.Unlock()
			err := lock.Unlock(ctx)
			if err != nil {
				return
			}
		}
	}
}

// 上报负载
func (r *RankingJobV1) reportLoad() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	load := r.load.Load()
	r.l.Debug(r.nodeID+" 上报负载: ", logger.Int32("load", load))
	r.redisClient.ZAdd(ctx, r.rankingLoadKey,
		redis.Z{Member: r.nodeID, Score: float64(load)})
	cancel()
	return
}

func (r *RankingJobV1) Close() error {
	r.localLock.Lock()
	lock := r.lock
	r.localLock.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var err *multierror.Error
	if lock != nil {
		err = multierror.Append(err, lock.Unlock(ctx))
	}
	if r.loadTicker != nil {
		r.loadTicker.Stop()
	}
	// 删除自己的负载
	err = multierror.Append(err, r.redisClient.ZRem(ctx, r.rankingLoadKey, redis.Z{Member: r.nodeID}).Err())
	return err
}
