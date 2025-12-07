package ginx

import (
	"Webook/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
)

var log logger.Logger = logger.NewNoOpLogger()

var vector *prometheus.CounterVec

func InitCounter(opt prometheus.CounterOpts) {
	vector = prometheus.NewCounterVec(opt, []string{"code"})
	prometheus.MustRegister(vector)
}
func SetLogger(l logger.Logger) {
	log = l
}

// WrapClaimsAndReq 统一解决handleFunc里面多数需要解析请求和进行jwt校验和返回数据的内容
func WrapClaimsAndReq[Req any](fn func(*gin.Context, Req, UserClaims) (Result, error)) gin.HandlerFunc {
	return func(context *gin.Context) {
		var req Req
		err := context.Bind(&req)
		if err != nil {
			log.Error("解析请求失败", logger.Error(err))
		}
		//放进上下文的UserClaim不带指针的
		value, ok := context.Get("user")
		if !ok {
			context.AbortWithStatus(http.StatusUnauthorized)
			log.Error("无法获取对应Claim", logger.String("path", context.Request.URL.Path))
			return
		}
		claims, ok := value.(UserClaims)
		if !ok {
			context.AbortWithStatus(http.StatusUnauthorized)
			log.Error("无法获取对应Claim", logger.String("path", context.Request.URL.Path))
			return
		}
		result, err := fn(context, req, claims)
		//	记录code日志
		vector.WithLabelValues(strconv.Itoa(result.Code)).Inc()
		if err != nil {
			log.Error("执行业务逻辑失败",
				logger.Error(err))
		}
		context.JSON(http.StatusOK, result)
	}
}

// WrapReq  用户不需要校验的请求
func WrapReq[Req any](fn func(*gin.Context, Req) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Req
		if err := ctx.Bind(&req); err != nil {
			log.Error("解析请求失败", logger.Error(err))
			return
		}
		res, err := fn(ctx, req)
		if err != nil {
			log.Error("执行业务逻辑失败",
				logger.Error(err))
			ctx.JSON(http.StatusUnauthorized, res)
			return
		}
		vector.WithLabelValues(strconv.Itoa(res.Code)).Inc()
		ctx.JSON(http.StatusOK, res)
	}
}

// WrapClaims 需要校验的get请求
func WrapClaims(fn func(*gin.Context, UserClaims) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 可以用包变量来配置，还是那句话，因为泛型的限制，这里只能用包变量
		rawVal, ok := ctx.Get("user")
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			log.Error("无法获得 claims",
				logger.String("path", ctx.Request.URL.Path))
			return
		}
		// 注意，这里要求放进去 ctx 的不能是*UserClaims，这是常见的一个错误
		claims, ok := rawVal.(UserClaims)
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			log.Error("无法获得 claims",
				logger.String("path", ctx.Request.URL.Path))
			return
		}
		res, err := fn(ctx, claims)
		if err != nil {
			log.Error("执行业务逻辑失败",
				logger.Error(err))
		}
		vector.WithLabelValues(strconv.Itoa(res.Code)).Inc()
		ctx.JSON(http.StatusOK, res)
	}
}

// Wrap 装饰器模式
func Wrap(fn func(*gin.Context) (Result, error)) gin.HandlerFunc {
	return func(context *gin.Context) {
		result, err := fn(context)
		if err != nil {
			log.Error("执行业务逻辑失败",
				logger.Error(err))
		}
		vector.WithLabelValues(strconv.Itoa(result.Code)).Inc()
		context.JSON(http.StatusOK, result)
	}
}
