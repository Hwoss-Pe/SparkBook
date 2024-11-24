package wechat

import (
	"Webook/oauth2/domain/wechat"
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type PrometheusDecorator struct {
	Service
	sum prometheus.Summary
}

func NewPrometheusDecorator(svc Service,
	namespace string,
	subsystem string,
	instanceId string,
	name string) *PrometheusDecorator {
	sum := prometheus.NewSummary(prometheus.SummaryOpts{
		Name:      name,
		Namespace: namespace,
		Subsystem: subsystem,
		ConstLabels: map[string]string{
			"instance_id": instanceId,
		},
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.9:   0.01,
			0.95:  0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	})
	prometheus.MustRegister(sum)
	return &PrometheusDecorator{
		Service: svc,
		sum:     sum,
	}
}
func (p *PrometheusDecorator) VerifyCode(ctx context.Context, code string) (wechat.WechatInfo, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		p.sum.Observe(float64(duration.Milliseconds()))
	}()
	return p.Service.VerifyCode(ctx, code)
}
