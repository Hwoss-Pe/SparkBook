package async

import (
	"Webook/pkg/logger"
	"Webook/sms/domain"
	"Webook/sms/repository"
	"Webook/sms/service"
	"context"
	"errors"
	"sync"
	"time"
)

type AsyncService struct {
	svc  service.Service
	l    logger.Logger
	repo repository.AsyncSmsRepository
}

func NewAsyncService(svc service.Service, l logger.Logger, repo repository.AsyncSmsRepository) service.Service {
	res := &AsyncService{svc: svc, l: l, repo: repo}
	//在创建的时候就直接启动这个异步处理进程，
	go func() {
		res.StartAsyncCycle()
	}()
	return res
}

func (s *AsyncService) StartAsyncCycle() {
	for {
		s.AsyncSend()
	}
}
func (s *AsyncService) AsyncSend() {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	//这个方法需要保证全实例只有一个能获取不重复的
	sms, err := s.repo.PreemptWaitingSMS(ctx)
	cancelFunc()
	switch {
	case err == nil:
		ctx, cancelFunc = context.WithTimeout(context.Background(), time.Second)
		defer cancelFunc()
		err := s.svc.Send(ctx, sms.TplId, sms.Args, sms.Numbers...)
		if err != nil {
			s.l.Error("执行异步发送短信失败",
				logger.Error(err),
				logger.Int64("id", sms.Id))
		}
		//如果是抢占到后发送成功
		res := err == nil
		// 通知 repository  这一次的执行结果
		err = s.repo.ReportScheduleResult(ctx, sms.Id, res)
		if err != nil {
			s.l.Error("执行异步发送短信成功，但是标记数据库失败",
				logger.Error(err),
				logger.Bool("res", res),
				logger.Int64("id", sms.Id))
		}
	case errors.Is(err, repository.ErrWaitingSMSNotFound):
		//	没抢到或者资源不存在
		time.Sleep(time.Second * 3)
	default:
		s.l.Error("抢占异步发送短信任务失败",
			logger.Error(err))
		//数据库出现问题，重试
		time.Sleep(time.Second * 3)
	}
}
func (s *AsyncService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	// 先判断是否需要异步
	if m.asyncMode {
		// 转储到数据库
		return s.repo.Add(ctx, domain.AsyncSms{
			TplId:    tplId,
			Args:     args,
			Numbers:  numbers,
			RetryMax: 3, // 设置可以重试三次
		})
	}
	// 同步发送
	start := time.Now()
	err := s.svc.Send(ctx, tplId, args, numbers...)
	duration := time.Since(start)
	// 记录统计数据
	m.recordMetrics(duration, err != nil)
	return err
}

func (s *AsyncService) needAsync() bool {
	// 1.1 使用绝对阈值，比如说直接发送的时候，（连续一段时间，或者连续N个请求）响应时间超过了 500ms，然后后续请求转异步
	// 1.2 或者send返回的错误率大于某个阈值
	// 2. 退出异步的方式就是连续请求和错误率都低于了阈值
	m.mu.Lock()
	defer m.mu.Unlock()
	m.shouldSwitchToAsync()
	return m.asyncMode
}

type metrics struct {
	mu              sync.Mutex
	responseTimes   []time.Duration // 最近请求的响应时间
	timestamps      []time.Time     // 对应的请求时间戳
	errorTimestamps []time.Time     // 错误请求的时间戳
	totalTimestamps []time.Time     // 所有请求的时间戳
	windowDuration  time.Duration   // 统计时间窗口
	asyncMode       bool            // 当前是否需要异步发送
	timeoutSize     int             // 连续超时 N 个请求的统计阈值
}

var m = metrics{
	responseTimes:   make([]time.Duration, 0, 100), // 假设滑动窗口大小为 100
	timestamps:      make([]time.Time, 0, 100),
	errorTimestamps: make([]time.Time, 0, 100),
	totalTimestamps: make([]time.Time, 0, 100),
	windowDuration:  time.Second * 60, // 统计过去 60 秒的数据
	asyncMode:       false,
	timeoutSize:     10,
}

// 记录请求的响应时间和错误状态
func (m *metrics) recordMetrics(responseTime time.Duration, isError bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 添加响应时间和时间戳
	m.responseTimes = append(m.responseTimes, responseTime)
	m.timestamps = append(m.timestamps, time.Now())
	m.totalTimestamps = append(m.totalTimestamps, time.Now())

	// 如果发生错误，记录错误的时间戳
	if isError {
		m.errorTimestamps = append(m.errorTimestamps, time.Now())
	}

	// 清理过期数据
	m.cleanup()
}

// 清理超出时间窗口的数据
func (m *metrics) cleanup() {
	now := time.Now()

	// 清理响应时间和时间戳
	validResponseTimes := make([]time.Duration, 0, len(m.responseTimes))
	validTimestamps := make([]time.Time, 0, len(m.timestamps))
	for i, ts := range m.timestamps {
		if now.Sub(ts) <= m.windowDuration {
			validResponseTimes = append(validResponseTimes, m.responseTimes[i])
			validTimestamps = append(validTimestamps, ts)
		}
	}
	m.responseTimes = validResponseTimes
	m.timestamps = validTimestamps

	// 清理错误请求时间戳
	validErrorTimestamps := make([]time.Time, 0, len(m.errorTimestamps))
	for _, ts := range m.errorTimestamps {
		if now.Sub(ts) <= m.windowDuration {
			validErrorTimestamps = append(validErrorTimestamps, ts)
		}
	}
	m.errorTimestamps = validErrorTimestamps

	// 清理总请求时间戳
	validTotalTimestamps := make([]time.Time, 0, len(m.totalTimestamps))
	for _, ts := range m.totalTimestamps {
		if now.Sub(ts) <= m.windowDuration {
			validTotalTimestamps = append(validTotalTimestamps, ts)
		}
	}
	m.totalTimestamps = validTotalTimestamps
}

// 判定是否需要切换到异步
func (m *metrics) shouldSwitchToAsync() {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查连续超时请求
	threshold := time.Millisecond * 500
	cnt := 0
	for _, t := range m.responseTimes {
		if t > threshold {
			cnt++
			if cnt >= m.timeoutSize {
				m.asyncMode = true
				return
			}
		} else {
			cnt = 0
		}
	}

	// 检查错误率是否超过阈值
	errorRateThreshold := 0.3 // 错误率阈值 30%
	if len(m.totalTimestamps) > 0 && float64(len(m.errorTimestamps))/float64(len(m.totalTimestamps)) > errorRateThreshold {
		m.asyncMode = true
		return
	}

	// 如果没有满足切换到异步的条件，则退出异步模式
	m.asyncMode = false
}

func (m *metrics) needAsync() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.asyncMode
}
