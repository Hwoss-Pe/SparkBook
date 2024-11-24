package cronjobx

import (
	"Webook/pkg/logger"
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"strconv"
	"time"
)

type CronJobBuilder struct {
	l      logger.Logger
	tracer trace.Tracer
	p      *prometheus.SummaryVec
}

func NewCronJobBuilder(l logger.Logger) *CronJobBuilder {
	vec := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "Hwoss",
		Subsystem: "webook",
		Help:      "统计定时任务的执行情况",
		Name:      "cron_job",
	}, []string{"name", "success"})
	prometheus.MustRegister(vec)
	return &CronJobBuilder{l: l, tracer: otel.GetTracerProvider().Tracer("webook/internal/job"), p: vec}
}
func (c *CronJobBuilder) Build(job Job) cron.Job {
	name := job.Name()
	//	这里需要返回一个接口，可以定义一个结构体去
	return cronJobFuncAdapter(
		func() error {
			//这里实现具体逻辑,开始链路追踪和打印日志
			_, span := c.tracer.Start(context.Background(), name)
			defer span.End()
			start := time.Now()
			c.l.Info("任务开始",
				logger.String("job", name))
			var success bool
			defer func() {
				c.l.Info("任务结束",
					logger.String("job", name))
				duration := time.Since(start).Milliseconds()
				c.p.WithLabelValues(name,
					strconv.FormatBool(success)).Observe(float64(duration))
			}()
			err := job.Run()
			success = err == nil
			if err != nil {
				span.RecordError(err)
				c.l.Error("运行任务失败", logger.Error(err),
					logger.String("job", name))
			}
			return nil
		})
}

type cronJobFuncAdapter func() error

func (c cronJobFuncAdapter) Run() {
	_ = c()
}
