package wechat

import (
	"Webook/oauth2/domain/wechat"
	"context"
)

//go:generate mockgen -source=./oauth2.go -package=wechatmocks -destination=mocks/svc.mock.go Service
type Service interface {
	AuthURL(ctx context.Context, state string) (string, error)
	// VerifyCode 目前大部分公司的 OAuth2 平台都差不多的设计
	// 返回一个 unionId。可以理解为，在第三方平台上的 unionId
	VerifyCode(ctx context.Context, code string) (wechat.WechatInfo, error)
}
