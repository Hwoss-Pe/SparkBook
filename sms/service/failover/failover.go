package failover

import (
	"Webook/sms/service"
	"context"
	"errors"
	"sync/atomic"
)

type FailoverSMSService struct {
	//可选择的 SMS Service 实现
	svcs []service.Service

	idx uint64
}

func NewFailoverSMSService(svcs []service.Service) *FailoverSMSService {
	return &FailoverSMSService{
		svcs: svcs,
	}
}

// SendV1 走的是遍历的数组的方式，这种会把大量流量打在切片第一个上面不采用
func (f *FailoverSMSService) SendV1(ctx context.Context, tplId string, args []string, numbers ...string) error {
	for _, svc := range f.svcs {
		err := svc.Send(ctx, tplId, args, numbers...)
		if err == nil {
			return nil
		}
	}
	return errors.New("发送失败，所有服务商都尝试过了")
}

func (f *FailoverSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	//全局的增长计数
	idx := atomic.LoadUint64(&f.idx)
	length := uint64(len(f.svcs))
	//用的是求余数方法进行遍历发送
	for i := idx; i < idx+length; i++ {
		index := i % length
		svc := f.svcs[index]
		err := svc.Send(ctx, tplId, args, numbers...)
		switch {
		case err == nil:
			return nil
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			//	调用者设置的超时时间到了
			// 调用者主动取消了
			return err
		}
	}
	atomic.AddUint64(&f.idx, 1)
	return errors.New("发送失败，所有服务商都尝试过了")
}
