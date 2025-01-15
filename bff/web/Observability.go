package web

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"time"
)

// ObservabilityHandler 用于可观测性的 Handler
type ObservabilityHandler struct {
}

func NewObservabilityHandler() *ObservabilityHandler {
	return &ObservabilityHandler{}
}

func (o *ObservabilityHandler) RegisterRoutes(s *gin.Engine) {
	tg := s.Group("/test")
	tg.GET("/random", o.Random)
}

func (o *ObservabilityHandler) Random(ctx *gin.Context) {
	num := rand.Int31n(1000)
	// 模拟响应时间
	time.Sleep(time.Millisecond * time.Duration(num))
	ctx.String(http.StatusOK, "OK")
}
