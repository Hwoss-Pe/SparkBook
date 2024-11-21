package prometheus

import (
	"Webook/pkg/grpcx/interceptor"
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"strings"
	"time"
)

type InterceptorBuilder struct {
	Namespace string
	Subsystem string
	interceptor.Builder
}

func (b *InterceptorBuilder) buildUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	summaryVec := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: b.Namespace,
		Subsystem: b.Subsystem,
		Name:      "server_handle_seconds",
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.9:   0.01,
			0.95:  0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	}, []string{"type", "service", "method", "peer", "code"})
	prometheus.MustRegister(summaryVec)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		start := time.Now()
		defer func() {
			serviceName, method := b.splitMethodName(info.FullMethod)
			st, _ := status.FromError(err)
			code := "OK"
			if st != nil {
				code = st.Code().String()
			}
			summaryVec.WithLabelValues("unary", serviceName, method, b.PeerName(ctx), code).
				Observe(float64(time.Since(start)))
		}()
		return handler(ctx, req)
	}
}
func (b *InterceptorBuilder) splitMethodName(fullMethodName string) (string, string) {
	fullMethodName = strings.TrimPrefix(fullMethodName, "/") // 前导斜杠
	if i := strings.Index(fullMethodName, "/"); i >= 0 {
		//分成服务名和对应的方法名
		return fullMethodName[:i], fullMethodName[i+1:]
	}
	return "unknown", "unknown"
}
