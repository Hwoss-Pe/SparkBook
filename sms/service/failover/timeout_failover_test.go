package failover

import (
	"Webook/sms/service"
	smsmocks "Webook/sms/service/mocks"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTimeoutFailoverSMSService_Send(t *testing.T) {
	testCases := []struct {
		name      string
		mock      func(ctrl *gomock.Controller) []service.Service
		threshold int32
		// 通过控制私有字段的取值，来模拟各种场景
		idx int32
		cnt int32

		wantErr error
		wantIdx int32
		wantCnt int32
	}{
		{
			name: "超时，但是没连续超时",
			mock: func(ctrl *gomock.Controller) []service.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(context.DeadlineExceeded)
				return []service.Service{svc0}
			},
			threshold: 3,
			wantErr:   context.DeadlineExceeded,
			wantCnt:   1,
			wantIdx:   0,
		},
		{
			name: "触发了切换，切换之后成功了",
			mock: func(ctrl *gomock.Controller) []service.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)

				return []service.Service{svc0, svc1}
			},
			threshold: 3,
			cnt:       3,
			// 重置了
			wantCnt: 0,
			// 切换到了 1
			wantIdx: 1,
		},
		{
			name: "触发了切换，切换之后失败",
			mock: func(ctrl *gomock.Controller) []service.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("发送失败"))

				return []service.Service{svc0, svc1}
			},
			threshold: 3,
			cnt:       3,
			// 重置了，因为不是超时错误，所以没有增加
			wantCnt: 0,
			// 切换到了 1
			wantIdx: 1,
			wantErr: errors.New("发送失败"),
		},
		{
			name: "触发了切换，切换之后依旧超时",
			mock: func(ctrl *gomock.Controller) []service.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(context.DeadlineExceeded)

				return []service.Service{svc0, svc1}
			},
			threshold: 3,
			cnt:       3,
			// 重置之后超时
			wantCnt: 1,
			// 切换到了 1
			wantIdx: 1,
			wantErr: context.DeadlineExceeded,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewTimeoutFailoverSMSService(tc.mock(ctrl), tc.threshold)
			svc.idx = tc.idx
			svc.cnt = tc.cnt

			err := svc.Send(context.Background(), "mytpl",
				[]string{}, "152xxx")
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantIdx, svc.idx)
			assert.Equal(t, tc.wantCnt, svc.cnt)
		})
	}
}
