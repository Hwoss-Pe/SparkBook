package service

import (
	"Webook/cronjob/domain"
	"Webook/cronjob/repository"
	repomocks "Webook/cronjob/repository/mocks"
	"Webook/pkg/logger"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func TestCronJobService_Preempt(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) repository.CronJobRepository
		wantErr  error
		wantJob  domain.CronJob
		interval time.Duration
	}{
		{
			name: "抢占并且续约",
			mock: func(ctrl *gomock.Controller) repository.CronJobRepository {
				repo := repomocks.NewMockCronJobRepository(ctrl)
				repo.EXPECT().Preempt(gomock.Any()).Return(domain.CronJob{
					Id: 1,
				}, nil)
				// interval 设置为三秒多，所以会续约三次
				repo.EXPECT().UpdateUtime(gomock.Any(), int64(1)).Times(3).
					Return(nil)
				repo.EXPECT().Release(gomock.Any(), int64(1)).Return(nil)
				return repo
			},
			interval: time.Second*3 + time.Millisecond*100,
			wantErr:  nil,
			wantJob: domain.CronJob{
				Id: 1,
			},
		},
		{
			name: "抢占失败",
			mock: func(ctrl *gomock.Controller) repository.CronJobRepository {
				repo := repomocks.NewMockCronJobRepository(ctrl)
				repo.EXPECT().Preempt(gomock.Any()).
					Return(domain.CronJob{}, errors.New("db error"))
				return repo
			},
			interval: time.Second*3 + time.Millisecond*100,
			wantErr:  errors.New("db error"),
			wantJob:  domain.CronJob{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//并行测试
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewCronJobService(tc.mock(ctrl), logger.NewNoOpLogger())
			svc.(*cronJobService).refreshInterval = time.Second
			job, err := svc.Preempt(context.Background())
			assert.Equal(t, tc.wantErr, err)
			// 因为后面还要处理，所以在 err != nil 的时候要返回
			if err != nil {
				return
			}
			assert.NotNil(t, job.CancelFunc)
			cancelFunc := job.CancelFunc
			job.CancelFunc = nil
			assert.Equal(t, tc.wantJob, job)
			time.Sleep(tc.interval)
			// 模拟运行之后取消续约
			cancelFunc()
			// 再次 sleep，借助 mock 确定真的退出了续约循环
			time.Sleep(tc.interval)
		})
	}
}
