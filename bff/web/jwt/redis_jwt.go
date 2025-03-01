package jwt

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

var RefreshTokenKey = []byte("moyn8y9abnd7q4zkq2m73yw8tu9j5ixA")
var AccessTokenKey = []byte("moyn8y9abnd7q4zkq2m73yw8tu9j5ixm")

type RedisHandler struct {
	cmd redis.Cmdable
	//长token的过期时间
	rtExpiration time.Duration
}

func (r *RedisHandler) ClearToken(ctx *gin.Context) error {
	// 正常用户的这两个 token 都会被前端更新
	// 也就是说在登录校验里面，走不到 redis 那一步就返回了
	ctx.Header("x-jwt-token", "")
	ctx.Header("x-refresh-token", "")

	uc := ctx.MustGet("user").(UserClaims)
	return r.cmd.Set(ctx, r.key(uc.Ssid), "", r.rtExpiration).Err()
}
func (r *RedisHandler) key(ssid string) string {
	return fmt.Sprintf("users:Ssid:%s", ssid)
}
func (r *RedisHandler) setRefreshToken(ctx *gin.Context, ssid string, uid int64) error {
	rc := RefreshClaims{
		Id:   uid,
		Ssid: ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 7 * 24)),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rc)
	refreshTokenStr, err := refreshToken.SignedString(RefreshTokenKey)
	if err != nil {
		return err
	}
	ctx.Header("x-refresh-token", refreshTokenStr)
	return nil
}
func (r *RedisHandler) SetLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := r.SetJWTToken(ctx, ssid, uid)
	if err != nil {
		return err
	}
	err = r.setRefreshToken(ctx, ssid, uid)
	return err
}

func (r *RedisHandler) SetJWTToken(ctx *gin.Context, ssid string, uid int64) error {
	//根据key，ssid和user-agent进行生成token，并且存储到x-jwt-token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, UserClaims{
		Id:        uid,
		Ssid:      ssid,
		UserAgent: ctx.GetHeader("User-Agent"),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
	})
	tokenStr, err := token.SignedString(AccessTokenKey)
	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

func (r *RedisHandler) CheckSession(ctx *gin.Context, ssid string) error {
	//在缓存里面检查,反向思维， 在redis里面只去记录已经退出登录的
	result, err := r.cmd.Exists(ctx, r.key(ssid)).Result()
	if err != nil {
		return err
	}
	if result > 0 {
		return errors.New("用户已经退出登录")
	}
	return nil
}

func (r *RedisHandler) ExtractTokenString(ctx *gin.Context) string {
	authCode := ctx.GetHeader("Authorization")
	if authCode == "" {
		return ""
	}
	//这里的格式是Bearer your_access_token
	seg := strings.SplitN(authCode, " ", 2)
	if len(seg) != 2 {
		//	格式不对
		return ""
	}
	return seg[1]
}

func NewRedisHandler(cmd redis.Cmdable) Handler {
	return &RedisHandler{
		cmd:          cmd,
		rtExpiration: time.Hour * 24 * 7,
	}
}
