package web

import (
	"Webook/payment/service/wechat"
	"Webook/pkg/ginx"
	"Webook/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
)

type WechatHandler struct {
	handler *notify.Handler
	l       logger.Logger
	svc     *wechat.NativePaymentService
}

func NewWechatHandler(handler *notify.Handler,
	nativeSvc *wechat.NativePaymentService,
	l logger.Logger) *WechatHandler {
	return &WechatHandler{
		handler: handler,
		svc:     nativeSvc,
		l:       l}
}

func (h *WechatHandler) RegisterRoutes(server *gin.Engine) {
	server.GET("hello", func(context *gin.Context) {
		context.String(200, "Hello, WeChat!")
	})
	server.Any("pay/callback", ginx.Wrap(h.HandleNative))
}

func (h *WechatHandler) HandleNative(ctx *gin.Context) (ginx.Result, error) {
	p := &payments.Transaction{}
	_, err := h.handler.ParseNotifyRequest(ctx, ctx.Request, p)
	if err != nil {
		return ginx.Result{}, err
	}
	err = h.svc.HandleCallback(ctx, p)
	return ginx.Result{}, err
}
