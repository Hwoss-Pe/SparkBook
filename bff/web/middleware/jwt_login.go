package middleware

import (
	jwt2 "Webook/bff/web/jwt"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ecodeclub/ekit/set"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

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
	s.Add("/articles/pub/list")
	s.Add("/articles/pub/ranking")
	s.Add("/search")
	s.Add("/users/recommend_authors")
	s.Add("/comment/list")
	s.Add("/comment/replies")
	s.Add("/follow/statics")
	return &JWTLoginMiddlewareBuilder{
		publicPaths: s,
		Handler:     hdl,
	}
}

func (j *JWTLoginMiddlewareBuilder) Build() gin.HandlerFunc {
	//逻辑是先进行url过滤,否则获取Authorization对应的jwt串,结合uc和key进行解析
	//对时间进行校验,在对ssid进行校验是否已经退出登录
	return func(ctx *gin.Context) {
		//目前先放开所有接口后续更新 TODO
		// 调试模式下设置一个默认用户
		// ctx.Set("user", jwt2.UserClaims{Id: 1})
		// return
		// 公开路径支持前缀匹配，例如 /articles/pub/:id
		if j.publicPaths.Exist(ctx.Request.URL.Path) || strings.HasPrefix(ctx.Request.URL.Path, "/articles/pub/") || strings.HasPrefix(ctx.Request.URL.Path, "/articles/author/") {
			// 公开路径也尝试解析 JWT，但不强制要求
			tokenStr := j.ExtractTokenString(ctx)
			if tokenStr != "" {
				uc := jwt2.UserClaims{}
				token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
					return jwt2.AccessTokenKey, nil
				})
				if err == nil && token.Valid {
					ctx.Set("user", uc)
				}
			}
			return
		}
		//提取jwt
		tokenStr := j.ExtractTokenString(ctx)
		if tokenStr == "" {
			fmt.Printf("JWT中间件: 未找到token, URL: %s\n", ctx.Request.URL.Path)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		fmt.Printf("JWT中间件: 提取到token: %s...\n", tokenStr[:min(len(tokenStr), 20)])

		uc := jwt2.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return jwt2.AccessTokenKey, nil
		})
		if err != nil || !token.Valid {
			// 不正确的 token
			fmt.Printf("JWT中间件: token解析失败, error: %v, valid: %v\n", err, token != nil && token.Valid)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		fmt.Printf("JWT中间件: token解析成功, 用户ID: %d, 会话ID: %s\n", uc.Id, uc.Ssid)
		expireTime, err := uc.GetExpirationTime()
		if err != nil {
			// 拿不到过期时间
			fmt.Printf("JWT中间件: 获取过期时间失败, error: %v\n", err)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if expireTime.Before(time.Now()) {
			// 已经过期
			fmt.Printf("JWT中间件: token已过期, 过期时间: %v, 当前时间: %v\n", expireTime, time.Now())
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
			fmt.Printf("JWT中间件: 会话检查失败, error: %v\n", err)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		fmt.Printf("JWT中间件: 验证成功, 设置用户信息到上下文\n")
		ctx.Set("user", uc)
	}
}
