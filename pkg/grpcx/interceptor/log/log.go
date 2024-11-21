package log

import (
	"Webook/pkg/grpcx/interceptor"
	"Webook/pkg/logger"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"runtime"
	"time"
)

type LoggerInterceptorBuilder struct {
	l logger.Logger
	interceptor.Builder
}

func NewLoggerInterceptorBuilder(l logger.Logger) *LoggerInterceptorBuilder {
	return &LoggerInterceptorBuilder{l: l}
}
func (l *LoggerInterceptorBuilder) BuildUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		//过滤掉grpc的探活日志
		if info.FullMethod == "/grpc.health.v1.Health/Check" {
			return handler(ctx, req)
		}
		var start = time.Now()
		var fields = make([]logger.Field, 0, 20)
		var event = "normal" //区分是否出现panic
		cost := time.Since(start)
		if rec := recover(); rec != nil {
			switch recType := rec.(type) {
			case error:
				rec = recType
			default:
				err = fmt.Errorf("%v", rec)
			}
			stack := make([]byte, 4096)
			stack = stack[:runtime.Stack(stack, true)]
			event = "recover"
			err = status.New(codes.Internal, "panic, err "+err.Error()).Err()
		}
		st, _ := status.FromError(err)
		defer func() {
			fields = append(fields,
				logger.String("type", "unary"),
				logger.String("code", st.Code().String()),
				logger.String("code_msg", st.Message()),
				logger.String("event", event),
				logger.String("method", info.FullMethod),
				logger.Int64("cost", cost.Milliseconds()),
				logger.String("peer", l.PeerName(ctx)),
				logger.String("peer_ip", l.PeerIP(ctx)),
			)
			l.l.Info("RPC调用", fields...)
		}()
		return handler(ctx, req)
	}
}
