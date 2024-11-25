package failover

import (
	"Webook/sms/service"
	"context"
	"errors"
	"sync/atomic"
)

type TimeoutFailoverSMSService struct {
	//lock sync.Mutex
	svcs []service.Service
	idx  int32
	// 连续超时次数
	cnt int32
	// 连续超时次数阈值
	threshold int32
}

func (t *TimeoutFailoverSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	cnt := atomic.LoadInt32(&t.cnt)
	idx := atomic.LoadInt32(&t.idx)
	//这里就是判断是否触发切换
	if cnt >= t.threshold {
		//为了防止并发切换，这里用CAS
		newIdx := (idx + 1) % int32(len(t.svcs))
		//CAS方法回自动检测并且进行赋值
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
			//说明切换的时候重置失败计数
			atomic.StoreInt32(&t.cnt, 0)
		}
		idx = newIdx
	}
	svc := t.svcs[idx]
	err := svc.Send(ctx, tplId, args, numbers...)
	switch {
	case err == nil:
		// 没有任何错误，重置计数器
		atomic.StoreInt32(&t.cnt, 0)
	case errors.Is(err, context.DeadlineExceeded):
		atomic.AddInt32(&t.cnt, 1)
	default:
	}
	return err
}

func NewTimeoutFailoverSMSService(svcs []service.Service, threshold int32) *TimeoutFailoverSMSService {
	return &TimeoutFailoverSMSService{
		svcs:      svcs,
		threshold: threshold,
	}
}
