package scheduler

import (
	"Webook/pkg/ginx"
	"Webook/pkg/gormx/connpool"
	"Webook/pkg/logger"
	"Webook/pkg/migrator"
	"Webook/pkg/migrator/event"
	"Webook/pkg/migrator/validator"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sync"
	"time"
)

// Scheduler 统一管理，需要实现对应的接口
type Scheduler[T migrator.Entity] struct {
	src     *gorm.DB
	dst     *gorm.DB
	lock    sync.Mutex
	l       logger.Logger
	pattern string
	//全量校验和增量校验
	cancelFull func()
	cancelIncr func()
	//双写逻辑
	pool *connpool.DoubleWritePool
	//异步消费逻辑
	producer event.Producer
}

func NewScheduler[T migrator.Entity](
	l logger.Logger,
	src *gorm.DB,
	dst *gorm.DB,
	pool *connpool.DoubleWritePool,
	producer event.Producer) *Scheduler[T] {
	return &Scheduler[T]{
		l:          l,
		src:        src,
		dst:        dst,
		pattern:    connpool.PatternSrcOnly,
		cancelFull: func() {},
		cancelIncr: func() {},
		pool:       pool,
		producer:   producer,
	}
}

func (s *Scheduler[T]) RegisterRoutes(server *gin.RouterGroup) {
	// 将这个暴露为 HTTP 接口修改对应的状态
	server.POST("/src_only", ginx.Wrap(s.SrcOnly))
	server.POST("/src_first", ginx.Wrap(s.SrcFirst))
	server.POST("/dst_first", ginx.Wrap(s.DstFirst))
	server.POST("/dst_only", ginx.Wrap(s.DstOnly))
	server.POST("/full/start", ginx.Wrap(s.StartFullValidation))
	server.POST("/full/stop", ginx.Wrap(s.StopFullValidation))
	server.POST("/incr/stop", ginx.Wrap(s.StopIncrementValidation))
	server.POST("/incr/start", ginx.WrapReq[StartIncrRequest](s.StartIncrementValidation))
}

// ---- 下面是四个阶段 ---- //
func (s *Scheduler[T]) SrcOnly(c *gin.Context) (ginx.Result, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.pattern = connpool.PatternSrcOnly
	s.pool.ChangePattern(connpool.PatternSrcOnly)
	return ginx.Result{
		Msg: "OK",
	}, nil
}

func (s *Scheduler[T]) SrcFirst(c *gin.Context) (ginx.Result, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.pattern = connpool.PatternSrcFirst
	s.pool.ChangePattern(connpool.PatternSrcFirst)
	return ginx.Result{
		Msg: "OK",
	}, nil
}

func (s *Scheduler[T]) DstFirst(c *gin.Context) (ginx.Result, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.pattern = connpool.PatternDstFirst
	s.pool.ChangePattern(connpool.PatternDstFirst)
	return ginx.Result{
		Msg: "OK",
	}, nil
}

func (s *Scheduler[T]) DstOnly(c *gin.Context) (ginx.Result, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.pattern = connpool.PatternDstOnly
	s.pool.ChangePattern(connpool.PatternDstOnly)
	return ginx.Result{
		Msg: "OK",
	}, nil
}
func (s *Scheduler[T]) newValidator() (*validator.Validator[T], error) {
	switch s.pattern {
	case connpool.PatternDstOnly, connpool.PatternDstFirst:
		return validator.NewValidator[T](s.dst, s.src, "DST", s.l, s.producer), nil
	case connpool.PatternSrcFirst, connpool.PatternSrcOnly:
		return validator.NewValidator[T](s.src, s.dst, "SRC", s.l, s.producer), nil
	default:
		return nil, fmt.Errorf("未知的 pattern %s", s.pattern)
	}
}

// StartFullValidation 全量校验相关
func (s *Scheduler[T]) StartFullValidation(c *gin.Context) (ginx.Result, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	cancel := s.cancelFull
	v, err := s.newValidator()
	if err != nil {
		return ginx.Result{}, err
	}
	var ctx context.Context
	//把ctx赋值回去可以控制取消
	ctx, s.cancelFull = context.WithCancel(context.Background())

	go func() {
		//	这里记得需要取消上一次的全量
		cancel()
		err := v.Validate(ctx)
		if err != nil {
			s.l.Warn("全量校验退出", logger.Error(err))
		}
	}()
	return ginx.Result{
		Msg: "OK",
	}, nil
}
func (s *Scheduler[T]) StopFullValidation(c *gin.Context) (ginx.Result, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.cancelFull()
	return ginx.Result{
		Msg: "OK",
	}, nil
}

// StartIncrementValidation 增量校验相关
//前端传入单位为毫秒的增量时间，和对应的utime

type StartIncrRequest struct {
	Utime int64 `json:"utime"`
	// 毫秒数
	// json 不能正确处理 time.Duration 类型
	Interval int64 `json:"interval"`
}

func (s *Scheduler[T]) StartIncrementValidation(c *gin.Context,
	req StartIncrRequest) (ginx.Result, error) {
	// 开启增量校验
	s.lock.Lock()

	defer s.lock.Unlock()
	// 取消上一次的
	cancel := s.cancelIncr
	v, err := s.newValidator()
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统异常",
		}, nil
	}
	v.SleepInterval(time.Duration(req.Interval) * time.Millisecond).Utime(req.Utime)
	var ctx context.Context
	ctx, s.cancelIncr = context.WithCancel(context.Background())

	go func() {
		cancel()
		err := v.Validate(ctx)
		s.l.Warn("退出增量校验", logger.Error(err))
	}()
	return ginx.Result{
		Msg: "启动增量校验成功",
	}, nil
}
func (s *Scheduler[T]) StopIncrementValidation(c *gin.Context) (ginx.Result, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.cancelIncr()
	return ginx.Result{
		Msg: "OK",
	}, nil
}
