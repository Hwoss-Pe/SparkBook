package jwt

import (
	"Webook/pkg/ginx"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

//go:generate mockgen -source=./types.go -package=jwtmocks -destination=./mocks/handler.mock.go Handler
type Handler interface {
	ClearToken(ctx *gin.Context) error
	SetLoginToken(ctx *gin.Context, uid int64) error
	SwtJWTToken(ctx *gin.Context, ssid string, uid int64) error
	CheckSession(ctx *gin.Context, ssid string) error
	ExtractTokenString(ctx *gin.Context) string
}
type RefreshClaims struct {
	Id   int64
	Ssid string
	jwt.RegisteredClaims
}
type UserClaims = ginx.UserClaims
