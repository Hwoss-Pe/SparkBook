package wechat

import (
	"Webook/oauth2/domain/wechat"
	"Webook/pkg/logger"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

const authURLPattern = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redire"

// 微信回调地址
var redirectURL = url.PathEscape("https://demo.com/oauth2/wechat/callback")

type service struct {
	appId     string
	appSecret string
	client    *http.Client
	logger    logger.Logger
}

func (s *service) AuthURL(ctx context.Context, state string) (string, error) {
	//利用格式化操作进行拼接
	return fmt.Sprintf(authURLPattern, s.appId, redirectURL, state), nil
}

func (s *service) VerifyCode(ctx context.Context, code string) (wechat.WechatInfo, error) {
	const baseURL = "https://api.weixin.qq.com/sns/oauth2/access_token"
	param := url.Values{}
	param.Set("appid", s.appId)
	param.Set("secret", s.appSecret)
	param.Set("code", code)
	param.Set("grant_type", "authorization_code")
	accessTokenURL := baseURL + "?" + param.Encode()
	req, err := http.NewRequest("GET", accessTokenURL, nil)
	if err != nil {
		return wechat.WechatInfo{}, err
	}
	req = req.WithContext(ctx)
	resp, err := s.client.Do(req)
	if err != nil {
		return wechat.WechatInfo{}, err
	}
	var res Result
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return wechat.WechatInfo{}, err
	}
	if res.ErrCode != 0 {
		return wechat.WechatInfo{}, errors.New("换取 access_token 失败")
	}
	defer resp.Body.Close()
	return wechat.WechatInfo{
		OpenId:  res.OpenId,
		UnionId: res.UnionId,
	}, nil
}

func NewService(appId, appSecret string,
	logger logger.Logger) Service {
	return &service{
		appId:     appId,
		appSecret: appSecret,
		client:    http.DefaultClient,
		logger:    logger,
	}
}

// Result 是微信返回的东西
type Result struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errMsg"`

	Scope string `json:"scope"`

	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`

	OpenId  string `json:"openid"`
	UnionId string `json:"unionid"`
}
