package integration

import (
	grpc2 "Webook/account/grpc"
	"Webook/account/integration/startup"
	"Webook/account/repository/dao"
	accountv1 "Webook/api/proto/gen/api/proto/account/v1"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
	"time"
)

type AccountServiceServerTestSuite struct {
	suite.Suite
	db     *gorm.DB
	server *grpc2.AccountServiceServer
}

func (s *AccountServiceServerTestSuite) SetupSuite() {
	s.db = startup.InitTestDB()
	s.server = startup.InitAccountService()
}

func (s *AccountServiceServerTestSuite) TearDownTest() {
	s.db.Exec("TRUNCATE TABLE `accounts`")
	//s.db.Exec("TRUNCATE TABLE `account_activities`")
}

func (s *AccountServiceServerTestSuite) TestCredit() {
	testCases := []struct {
		name    string
		before  func(t *testing.T)
		after   func(t *testing.T)
		req     *accountv1.CreditRequest
		wantErr error
	}{
		{
			name: "用户账号不存在",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				var sysAccount dao.Account
				err := s.db.WithContext(ctx).Where("type = ?", uint8(accountv1.AccountType_AccountTypeSystem)).
					First(&sysAccount).Error
				assert.NoError(t, err)
				assert.Equal(t, int64(10), sysAccount.Balance)
				var usrAccount dao.Account
				err = s.db.WithContext(ctx).Where("uid = ?", 1024).
					First(&usrAccount).Error
				require.NoError(t, err)
				usrAccount.Id = 0
				assert.True(t, usrAccount.Ctime > 0)
				usrAccount.Ctime = 0
				assert.True(t, usrAccount.Utime > 0)
				usrAccount.Utime = 0
				assert.Equal(t, dao.Account{
					Account:  123,
					Uid:      1024,
					Type:     uint8(accountv1.AccountType_AccountTypeReward),
					Balance:  100,
					Currency: "CNY",
				}, usrAccount)
			},
			req: &accountv1.CreditRequest{
				Biz:   "test",
				BizId: 123,
				Items: []*accountv1.CreditItem{
					{
						Account:     123,
						AccountType: accountv1.AccountType_AccountTypeReward,
						Amt:         100,
						Currency:    "CNY",
						Uid:         1024,
					},
					{
						AccountType: accountv1.AccountType_AccountTypeSystem,
						Amt:         10,
						Currency:    "CNY",
					},
				},
			},
		},
		{
			name: "用户账号存在",
			before: func(t *testing.T) {
				err := s.db.Create(&dao.Account{
					Uid:      1025,
					Account:  123,
					Type:     uint8(accountv1.AccountType_AccountTypeReward),
					Balance:  300,
					Currency: "CNY",
					Ctime:    1111,
					Utime:    2222,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				var usrAccount dao.Account
				err := s.db.WithContext(ctx).Where("uid = ?", 1025).
					First(&usrAccount).Error
				require.NoError(t, err)
				usrAccount.Id = 0
				assert.True(t, usrAccount.Ctime > 0)
				usrAccount.Ctime = 0
				assert.True(t, usrAccount.Utime > 0)
				usrAccount.Utime = 0
				assert.Equal(t, dao.Account{
					Account:  123,
					Uid:      1025,
					Type:     uint8(accountv1.AccountType_AccountTypeReward),
					Balance:  400,
					Currency: "CNY",
				}, usrAccount)
			},
			req: &accountv1.CreditRequest{
				Biz:   "test",
				BizId: 123,
				Items: []*accountv1.CreditItem{
					{
						Account:     123,
						AccountType: accountv1.AccountType_AccountTypeReward,
						Amt:         100,
						Currency:    "CNY",
						Uid:         1025,
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			tc.before(t)
			_, err := s.server.Credit(context.Background(), tc.req)
			assert.Equal(t, tc.wantErr, err)
			tc.after(t)
		})
	}
}

func TestAccountServiceServer(t *testing.T) {
	suite.Run(t, new(AccountServiceServerTestSuite))
}
