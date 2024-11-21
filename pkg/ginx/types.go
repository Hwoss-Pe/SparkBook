package ginx

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Server struct {
	*gin.Engine
	Addr string
}
type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}
type UserClaims struct {
	Id        int64
	UserAgent string
	Ssid      string
	jwt.RegisteredClaims
}

func (s *Server) Start() error {
	return s.Engine.Run(s.Addr)
}
