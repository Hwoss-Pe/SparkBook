package web

import (
	oauth2v1 "Webook/api/proto/gen/api/proto/oauth2/v1"
	userv1 "Webook/api/proto/gen/api/proto/user/v1"
	jwt2 "Webook/bff/web/jwt"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
)

//这个做的是扫码登录

type OAuth2WechatHandler struct {
	wechatSvc       oauth2v1.Oauth2ServiceClient
	userSvc         userv1.UsersServiceClient
	stateCookieName string
	stateTokenKey   []byte
	jwt2.Handler
}

func NewOAuth2WechatHandler(service oauth2v1.Oauth2ServiceClient,
	userSvc userv1.UsersServiceClient,
	jwthdl jwt2.Handler) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		wechatSvc:       service,
		userSvc:         userSvc,
		stateCookieName: "jwt-state",
		stateTokenKey:   []byte("moyn8y9abnd7q4zkq2m73yw8tu9j5ixB"),
		Handler:         jwthdl,
	}
}

func (h *OAuth2WechatHandler) RegisterRoute(s *gin.Engine) {
	g := s.Group("/oauth2/wechat")
	g.GET("/authurl", h.OAuth2URL)
	// 这边用 Any 万无一失
	g.Any("/callback", h.Callback)
}

func (h *OAuth2WechatHandler) OAuth2URL(ctx *gin.Context) {
	state := uuid.New().String()
	url, err := h.wechatSvc.AuthURL(ctx, &oauth2v1.AuthURLRequest{
		State: state,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误，请稍后再试",
		})
		return
	}
	err = h.setStateCookie(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误，请稍后再试",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: url,
	})
	return
}

func (h *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	err := h.verifyState(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统异常，请重试",
		})
		return
	}
	code := ctx.Query("code")
	info, err := h.wechatSvc.VerifyCode(ctx, &oauth2v1.VerifyCodeRequest{
		Code: code,
	})
	if err != nil {
		// 实际上这个错误，也有可能是 code 不对
		// 但是给前端的信息没有太大的必要区分究竟是代码不对还是系统本身有问题
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	u, err := h.userSvc.FindOrCreateByWechat(ctx,
		&userv1.FindOrCreateByWechatRequest{
			Info: &userv1.WechatInfo{
				OpenId:  info.OpenId,
				UnionId: info.UnionId,
			},
		})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	err = h.SetLoginToken(ctx, u.User.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "登录成功",
	})
}

func (h *OAuth2WechatHandler) verifyState(ctx *gin.Context) error {
	state := ctx.Query("state")
	ck, err := ctx.Cookie(h.stateCookieName)
	if err != nil {
		return fmt.Errorf("%w, 无法获得 cookie", err)
	}
	var sc StateClaims
	_, err = jwt.ParseWithClaims(ck, &sc, func(token *jwt.Token) (interface{}, error) {
		return h.stateTokenKey, nil
	})
	if err != nil {
		return fmt.Errorf("%w, cookie 不是合法 JWT token", err)
	}
	if sc.State != state {
		return errors.New("state 被篡改了")
	}
	return nil
}
func (h *OAuth2WechatHandler) setStateCookie(ctx *gin.Context, state string) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, StateClaims{
		State: state,
	})
	tokenStr, err := token.SignedString(h.stateTokenKey)
	if err != nil {
		return err
	}
	ctx.SetCookie(h.stateCookieName, tokenStr,
		600,
		"/oauth2/wechat/callback",
		// 这边把 HTTPS 协议禁止了。不过在生产环境中要开启。
		"", false, true)
	return nil
}

type StateClaims struct {
	State string
	jwt.RegisteredClaims
}
