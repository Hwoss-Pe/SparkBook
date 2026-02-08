package ioc

import (
	"Webook/bff/web"
	jwt2 "Webook/bff/web/jwt"
	"Webook/bff/web/middleware"
	"Webook/pkg/ginx"
	metricsmw "Webook/pkg/ginx/middleware/metrics"
	"Webook/pkg/logger"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
)

var prometheusOnce sync.Once

func InitGinServer(l logger.Logger,
	jwtHdl jwt2.Handler,
	user *web.UserHandler,
	article *web.ArticleHandler,
	reward *web.RewardHandler,
	comment *web.CommentHandler,
	follow *web.FollowHandler,
	notification *web.NotificationHandler) *ginx.Server {
	initPrometheusEndpoint(l)
	engine := gin.Default()
	b := newHTTPMetricsBuilder()
	engine.Use(b.BuildActiveRequest(), b.BuildResponseTime(), b.BuildRequestCount(), corsHdl(), timeout(), middleware.NewJWTLoginMiddlewareBuilder(jwtHdl).Build())
	engine.Static("/static", "./media")
	user.RegisterRoute(engine)
	article.RegisterRoute(engine)
	reward.RegisterRoute(engine)
	comment.RegisterRoute(engine)
	follow.RegisterRoute(engine)
	notification.RegisterRoute(engine)
	web.NewUploadHandler(l).RegisterRoute(engine)
	addr := viper.GetString("http.addr")
	metricsAddr := viper.GetString("prometheus.addr")
	if metricsAddr == "" {
		metricsAddr = ":18081"
	}
	l.Info("bff server initialized",
		logger.String("http_addr", addr),
		logger.String("metrics_addr", metricsAddr),
		logger.String("static_dir", "./media"))
	ginx.InitCounter(prometheus.CounterOpts{
		Namespace: "hwoss",
		Subsystem: "webbok_bff",
		Name:      "http",
	})
	ginx.SetLogger(l)
	return &ginx.Server{
		Engine: engine,
		Addr:   addr,
	}
}

func newHTTPMetricsBuilder() metricsmw.PrometheusBuilder {
	instanceID := viper.GetString("prometheus.instanceId")
	if instanceID == "" {
		h, err := os.Hostname()
		if err == nil && h != "" {
			instanceID = h
		} else {
			instanceID = "unknown"
		}
	}
	return metricsmw.PrometheusBuilder{
		Namespace:  "hwoss",
		Subsystem:  "webbok_bff",
		Name:       "http",
		InstanceID: instanceID,
	}
}

func initPrometheusEndpoint(l logger.Logger) {
	prometheusOnce.Do(func() {
		// 进程/运行时指标（CPU/内存/GC 等）。若已注册则忽略。
		registerCollectorOrIgnore(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
		registerCollectorOrIgnore(collectors.NewGoCollector())

		addr := viper.GetString("prometheus.addr")
		if addr == "" {
			addr = ":18081"
		}
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		go func() {
			err := http.ListenAndServe(addr, mux)
			if err != nil {
				l.Error("启动 Prometheus metrics 端点失败",
					logger.Error(fmt.Errorf("addr=%s, err=%w", addr, err)))
			}
		}()
	})
}

func registerCollectorOrIgnore(c prometheus.Collector) {
	err := prometheus.Register(c)
	if err == nil {
		return
	}
	_, ok := err.(prometheus.AlreadyRegisteredError)
	if ok {
		return
	}
	panic(err)
}
func timeout() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, ok := ctx.Request.Context().Deadline()
		if !ok {
			withTimeout, cancelFunc := context.WithTimeout(ctx.Request.Context(), time.Second*10)
			defer cancelFunc()
			ctx.Request = ctx.Request.Clone(withTimeout)
		}
		ctx.Next()
	}
}
func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowHeaders: []string{"Content-Type", "Authorization"},

		ExposeHeaders: []string{"x-jwt-token", "x-refresh-token"},

		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "yourcompany.com")
		},
		MaxAge: 12 * time.Hour,
	})
}
