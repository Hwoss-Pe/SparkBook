package web

import (
	"Webook/pkg/ginx"
	"github.com/gin-gonic/gin"
)

type Result = ginx.Result

type Page struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type Handler interface {
	RegisterRoute(s *gin.Engine)
}
