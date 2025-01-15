package middleware

import (
	jwt2 "Webook/bff/web/jwt"
	"github.com/ecodeclub/ekit/set"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

type JWTLoginMiddlewareBuilder struct {
	publicPaths set.Set[string]
	jwt2.Handler
}

func NewJWTLoginMiddlewareBuilder(hdl jwt2.Handler) *JWTLoginMiddlewareBuilder {
	s := set.NewMapSet[string](3)
	s.Add("/users/signup")
	s.Add("/users/login_sms/code/send")
	s.Add("/users/login_sms")
	s.Add("/users/refresh_token")
	s.Add("/users/login")
	s.Add("/oauth2/wechat/authurl")
	s.Add("/oauth2/wechat/callback")
	s.Add("/test/random")
	return &JWTLoginMiddlewareBuilder{
		publicPaths: s,
		Handler:     hdl,
	}
}

func (j *JWTLoginMiddlewareBuilder) Build() gin.HandlerFunc {
	//逻辑是先进行url过滤,否则获取Authorization对应的jwt串,结合uc和key进行解析
	//对时间进行校验,在对ssid进行校验是否已经退出登录
	return func(ctx *gin.Context) {
		if j.publicPaths.Exist(ctx.Request.URL.Path) {
			return
		}
		//提取jwt
		tokenStr := j.ExtractTokenString(ctx)
		uc := jwt2.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, uc, func(token *jwt.Token) (interface{}, error) {
			return jwt2.AccessTokenKey, nil
		})
		if err != nil || !token.Valid {
			// 不正确的 token
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		expireTime, err := uc.GetExpirationTime()
		if err != nil {
			// 拿不到过期时间
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if expireTime.Before(time.Now()) {
			// 已经过期
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//if ctx.GetHeader("User-Agent") != uc.UserAgent {
		//	// 换了一个 User-Agent，可能是攻击者
		//	ctx.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}
		err = j.CheckSession(ctx, uc.Ssid)
		if err != nil {
			// 系统错误或者用户已经主动退出登录了
			// 这里也可以考虑说，如果在 Redis 已经崩溃的时候，
			// 就不要去校验是不是已经主动退出登录了。
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("user", uc)
	}
}
