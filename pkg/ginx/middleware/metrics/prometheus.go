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
	// 这一个实例名字，你可以考虑使用 本地 IP，
	// 又或者在启动的时候配置一个 ID
	InstanceID string
}

// BuildResponseTime 统计返回的请求时间
func (p *PrometheusBuilder) BuildResponseTime() gin.HandlerFunc {
	labels := []string{"method", "pattern", "status"}
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:   p.Namespace,
		Subsystem:   p.Subsystem,
		Name:        p.Name + "_resp_time",
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
			vector.WithLabelValues(
				method, context.FullPath(),
				strconv.Itoa(context.Writer.Status())).
				Observe(float64(duration))
		}()
		context.Next()
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
		defer gauge.Desc()
		context.Next()
	}
}
