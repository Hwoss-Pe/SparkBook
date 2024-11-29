package integration

import (
	pmtv1 "Webook/api/proto/gen/api/proto/payment/v1"
	pmtmocks "Webook/api/proto/gen/api/proto/payment/v1/mocks"
	"Webook/reward/domain"
	"Webook/reward/integration/startup"
	"Webook/reward/repository/dao"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
	"testing"
	"time"
)

type WechatNativeRewardServiceTestSuite struct {
	suite.Suite
	rdb redis.Cmdable
	db  *gorm.DB
}

func (s *WechatNativeRewardServiceTestSuite) TestPreReward() {
	t := s.T()
	testCases := []struct {
		name string
		// 实在不想真的跟微信打交道，先保证自己这边没问题
		mock   func(ctrl *gomock.Controller) pmtv1.WechatPaymentServiceClient
		before func(t *testing.T)
		after  func(t *testing.T)

		r domain.Reward

		wantData string
		wantErr  error
	}{
		{
			name: "直接创建成功",
			mock: func(ctrl *gomock.Controller) pmtv1.WechatPaymentServiceClient {
				client := pmtmocks.NewMockWechatPaymentServiceClient(ctrl)
				client.EXPECT().NativePrePay(gomock.Any(), gomock.Any()).
					Return(&pmtv1.NativePrePayResponse{
						CodeUrl: "test_url",
					}, nil)
				return client
			},
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				// 验证数据和缓存
				var r dao.Reward
				err := s.db.Where("biz = ? AND biz_id = ?", "test", 1).First(&r).Error
				assert.NoError(t, err)
				assert.True(t, r.Id > 0)
				r.Id = 0
				assert.True(t, r.Ctime > 0)
				r.Ctime = 0
				assert.True(t, r.Utime > 0)
				r.Utime = 0
				assert.Equal(t, dao.Reward{
					Biz:       "test",
					BizId:     1,
					BizName:   "测试项目",
					TargetUid: 1234,
					Uid:       123,
					Amount:    1,
				}, r)

				codeURL, err := s.rdb.GetDel(ctx, s.codeURLKey("test", 1, 123)).Result()
				require.NoError(t, err)
				assert.Equal(t, "test_url", codeURL)
			},
			r: domain.Reward{
				Uid: 123,
				Target: domain.Target{
					Biz:     "test",
					BizId:   1,
					BizName: "测试项目",
					Uid:     1234,
				},
				Amt: 1,
			},
			wantData: "test_url",
		},
		{
			name: "拿到缓存",
			mock: func(ctrl *gomock.Controller) pmtv1.WechatPaymentServiceClient {
				client := pmtmocks.NewMockWechatPaymentServiceClient(ctrl)
				return client
			},
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				err := s.rdb.Set(ctx, s.codeURLKey("test", 2, 123), "test_url_1", time.Minute).Err()
				require.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				codeURL, err := s.rdb.GetDel(ctx, s.codeURLKey("test", 2, 123)).Result()
				require.NoError(t, err)
				assert.Equal(t, "test_url_1", codeURL)
			},
			r: domain.Reward{
				Uid: 123,
				Target: domain.Target{
					Biz:     "test",
					BizId:   2,
					BizName: "测试项目",
					Uid:     1234,
				},
				Amt: 1,
			},
			wantData: "test_url_1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := startup.InitWechatNativeSvc(tc.mock(ctrl))
			codeURL, err := svc.PreReward(context.Background(), tc.r)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantData, codeURL)
			tc.after(t)
		})
	}
}

func (s *WechatNativeRewardServiceTestSuite) SetupSuite() {
	s.rdb = startup.InitRedis()
	s.db = startup.InitTestDB()
}

func (s *WechatNativeRewardServiceTestSuite) TearDownTest() {
	s.db.Exec("TRUNCATE TABLE rewards")
}

func (s *WechatNativeRewardServiceTestSuite) codeURLKey(biz string, bizId, uid int64) string {
	return fmt.Sprintf("reward:code_url:%s:%d:%d",
		biz, bizId, uid)
}

func TestWechatNativeRewardService(t *testing.T) {
	suite.Run(t, new(WechatNativeRewardServiceTestSuite))
}
