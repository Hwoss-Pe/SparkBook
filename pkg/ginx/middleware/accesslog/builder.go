package accesslog

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

// MiddlewareBuilder 可以显式的控制打印请求和响应的日志
type MiddlewareBuilder struct {
	allowReqBody  bool
	allowRespBody bool
	logFunc       func(ctx context.Context, al AccessLog)
}

func NewMiddlewareBuilder(logFunc func(ctx context.Context, al AccessLog)) *MiddlewareBuilder {
	return &MiddlewareBuilder{allowReqBody: true, allowRespBody: false, logFunc: logFunc}
}

func (b *MiddlewareBuilder) AllowReqBody() *MiddlewareBuilder {
	b.allowReqBody = true
	return b
}

func (b *MiddlewareBuilder) AllowRespBody() *MiddlewareBuilder {
	b.allowRespBody = true
	return b
}

type AccessLog struct {
	Method     string `json:"method"`
	Path       string `json:"path"`
	ReqBody    string `json:"req_body"`
	Duration   string `json:"duration"`
	StatusCode int    `json:"status_code"`
	RespBody   string `json:"resp_body"`
}

func (b *MiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		al := AccessLog{
			Method: ctx.Request.Method,
			Path:   ctx.Request.URL.Path,
		}
		if b.allowReqBody && ctx.Request.Body != nil {
			//这个是个stream流，只能读一次，所以读完要给他放进去
			data, _ := ctx.GetRawData()
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(data))
			al.ReqBody = string(data)
		}
		//由于context没有提供对应的可以直接操作resp，因此自己用一些write
		if b.allowRespBody {
			ctx.Writer = responseWriter{
				ResponseWriter: ctx.Writer,
				al:             &al,
			}
		}
		defer func() {
			duration := time.Since(start)
			al.Duration = duration.String()
			b.logFunc(ctx, al)
		}()
		ctx.Next()
	}
}

type responseWriter struct {
	al *AccessLog
	//ResponseWriter是下面这个接口的一个子接口，只需要重写它的方法就行，或者直接实现http.ResponseWriter
	gin.ResponseWriter
}

func (r responseWriter) WriteHeader(statusCode int) {
	r.al.StatusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r responseWriter) Write(data []byte) (int, error) {
	r.al.RespBody = string(data)
	return r.ResponseWriter.Write(data)
}

func (r responseWriter) WriteString(data string) (int, error) {
	r.al.RespBody = data
	return r.ResponseWriter.WriteString(data)
}
