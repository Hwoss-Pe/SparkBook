//go:build wireinject

package startup

import (
	accountv1 "Webook/api/proto/gen/api/proto/account/v1"
	pmtv1 "Webook/api/proto/gen/api/proto/payment/v1"
	"Webook/reward/repository"
	"Webook/reward/repository/cache"
	"Webook/reward/repository/dao"
	"Webook/reward/service"
	"github.com/google/wire"
)

var thirdPartySet = wire.NewSet(InitTestDB, InitLogger, InitRedis)

func InitWechatNativeSvc(client pmtv1.WechatPaymentServiceClient, client2 accountv1.AccountServiceClient) *service.WechatNativeRewardService {
	wire.Build(service.NewWechatNativeRewardService,
		thirdPartySet,
		cache.NewRewardRedisCache,
		repository.NewRewardRepository, dao.NewRewardGORMDAO)
	return new(service.WechatNativeRewardService)
}
