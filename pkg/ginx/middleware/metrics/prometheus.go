package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

type PrometheusBuilder struct {
	//命名空间和子系统
	Namespace string
	Subsystem string
	Name      string
	//可有可无
	Help string
	// 这一个实例名字，   可以考虑使用 本地 IP，
	// 又或者在启动的时候配置一个 ID
	InstanceID string
}

// BuildResponseTime 统计返回的请求时间
func (p *PrometheusBuilder) BuildResponseTime() gin.HandlerFunc {
	labels := []string{"method", "pattern", "status"}
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:   p.Namespace,
		Subsystem:   p.Subsystem,
		Name:        p.Name + "_resp_time_ms",
		ConstLabels: map[string]string{"instance_id": p.InstanceID},
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.90:  0.01,
			0.99:  0.001,
			0.999: 0.0001,
		}}, labels)
	prometheus.MustRegister(vector)
	return func(context *gin.Context) {
		method := context.Request.Method
		start := time.Now()
		defer func() {
			duration := time.Since(start).Milliseconds()
			pattern := context.FullPath()
			if pattern == "" {
				pattern = "__unknown__"
			}
			vector.WithLabelValues(
				method, pattern,
				strconv.Itoa(context.Writer.Status())).
				Observe(float64(duration))
		}()
		context.Next()
	}
}

// BuildRequestCount 统计请求总量，用于计算 QPS，并按状态码聚合。
// PromQL 示例：sum(rate(<namespace>_<subsystem>_<name>_requests_total[1m])) by (method, pattern)
func (p *PrometheusBuilder) BuildRequestCount() gin.HandlerFunc {
	labels := []string{"method", "pattern", "status"}
	vector := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   p.Namespace,
		Subsystem:   p.Subsystem,
		Name:        p.Name + "_requests_total",
		Help:        "HTTP requests total",
		ConstLabels: map[string]string{"instance_id": p.InstanceID},
	}, labels)
	prometheus.MustRegister(vector)
	return func(context *gin.Context) {
		context.Next()
		pattern := context.FullPath()
		if pattern == "" {
			pattern = "__unknown__"
		}
		vector.WithLabelValues(
			context.Request.Method,
			pattern,
			strconv.Itoa(context.Writer.Status()),
		).Inc()
	}
}

// BuildActiveRequest 统计活跃的请求数
func (p *PrometheusBuilder) BuildActiveRequest() gin.HandlerFunc {
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   p.Namespace,
		Subsystem:   p.Subsystem,
		Name:        p.Name + "_active_req",
		ConstLabels: map[string]string{"instance_id": p.InstanceID}},
	)

	prometheus.MustRegister(gauge)
	return func(context *gin.Context) {
		gauge.Inc()
		defer gauge.Dec()
		context.Next()
	}
}
