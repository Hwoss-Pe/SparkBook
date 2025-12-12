package web

import (
	"Webook/pkg/ginx"
	"Webook/pkg/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type UploadHandler struct {
	l logger.Logger
}

func NewUploadHandler(l logger.Logger) *UploadHandler {
	return &UploadHandler{l: l}
}

func (h *UploadHandler) RegisterRoute(s *gin.Engine) {
	g := s.Group("/upload")
	g.POST("/avatar", ginx.WrapClaims(h.UploadAvatar))
	g.POST("/cover", ginx.WrapClaims(h.UploadCover))
}

func (h *UploadHandler) UploadAvatar(ctx *gin.Context, claims ginx.UserClaims) (ginx.Result, error) {
	url, err := h.save(ctx, "avatars", fmt.Sprintf("%d", claims.Id))
	if err != nil {
		h.l.Error("上传头像失败", logger.Error(err))
		return ginx.Result{Code: 5, Msg: "系统错误"}, err
	}
	return ginx.Result{Data: url}, nil
}

func (h *UploadHandler) UploadCover(ctx *gin.Context, claims ginx.UserClaims) (ginx.Result, error) {
	aid := strings.TrimSpace(ctx.PostForm("articleId"))
	sub := aid
	if sub == "" {
		sub = fmt.Sprintf("%d", claims.Id)
	}
	url, err := h.save(ctx, "covers", sub)
	if err != nil {
		h.l.Error("上传封面失败", logger.Error(err))
		return ginx.Result{Code: 5, Msg: "系统错误"}, err
	}
	return ginx.Result{Data: url}, nil
}

func (h *UploadHandler) save(ctx *gin.Context, kind string, subdir string) (string, error) {
	file, err := ctx.FormFile("file")
	if err != nil {
		return "", err
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowExt(ext) {
		ctx.Status(http.StatusBadRequest)
		return "", http.ErrNotSupported
	}
	now := time.Now()
	dir := filepath.Join("media", kind, subdir, now.Format("2006"), now.Format("01"))
	err = os.MkdirAll(dir, 0o755)
	if err != nil {
		return "", err
	}
	name := uuid.NewString() + ext
	full := filepath.Join(dir, name)
	err = ctx.SaveUploadedFile(file, full)
	if err != nil {
		return "", err
	}
	rel := filepath.ToSlash(filepath.Join("/static", kind, subdir, now.Format("2006"), now.Format("01"), name))
	return rel, nil
}

func allowExt(ext string) bool {
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
		return true
	default:
		return false
	}
}
