package metric

import (
	"context"
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"net"
	"strconv"
	"time"
)

type PrometheusHook struct {
	vector *prometheus.SummaryVec
}

func NewPrometheusHook(
	namespace string,
	subsystem string,
	instanceId string,
	name string,
) *PrometheusHook {
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      name,
		ConstLabels: map[string]string{
			"instance_id": instanceId,
		},
	}, []string{"cmd", "key_exist"})
	return &PrometheusHook{
		vector: vector,
	}
}
func (p PrometheusHook) DialHook(next redis.DialHook) redis.DialHook {
	//建立链接的时候的进行操作
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}

func (p PrometheusHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		start := time.Now()
		var err error
		defer func() {
			duration := time.Since(start)
			keyExist := errors.Is(err, redis.Nil)
			p.vector.WithLabelValues(cmd.Name(), strconv.FormatBool(keyExist)).
				Observe(float64(duration.Milliseconds()))
		}()
		err = next(ctx, cmd)
		return err
	}
}

func (p PrometheusHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		return next(ctx, cmds)
	}
}
