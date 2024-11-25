package ratelimit

import (
	"Webook/pkg/ratelimit"
	"Webook/sms/service"
	"context"
	"errors"
	"fmt"
)

// 假设针对是腾讯sms
const key = "sms_tencent"

var errLimited = errors.New("短信服务触发限流")

type LimitSmsService struct {
	svc     service.Service
	limiter ratelimit.Limiter
}

func NewLimitSMSService(svc service.Service, limiter ratelimit.Limiter) *LimitSmsService {
	return &LimitSmsService{
		svc:     svc,
		limiter: limiter,
	}
}
func (l *LimitSmsService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	limit, err := l.limiter.Limit(ctx, key)
	if err != nil {
		return fmt.Errorf("短信服务判断是否限流异常 %w", err)
	}
	if limit {
		//限流直接返回
		return errLimited
	}
	return l.svc.Send(ctx, tplId, args, numbers...)
}
