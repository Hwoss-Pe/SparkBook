package ioc

import (
	"Webook/pkg/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"time"
)

func InitLogger() logger.Logger {
	// 就是直接使用 zap 本身的配置结构体来处理
	// 配置Lumberjack以支持日志文件的滚动
	lumberjackLogger := &lumberjack.Logger{
		// 注意有没有权限
		Filename:   "/var/log/user.log", // 指定日志文件路径
		MaxSize:    50,                  // 每个日志文件的最大大小，单位：MB
		MaxBackups: 3,                   // 保留旧日志文件的最大个数
		MaxAge:     28,                  // 保留旧日志文件的最大天数
		Compress:   true,                // 是否压缩旧的日志文件
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(lumberjackLogger),
		zapcore.DebugLevel)

	l := zap.New(core, zap.AddCaller())
	res := logger.NewZapLogger(l)
	go func() {
		// 测试 ELK，我直接输出日志
		ticker := time.NewTicker(time.Millisecond * 1000)
		for t := range ticker.C {
			res.Info("模拟输出日志", logger.String("time", t.String()))
		}
	}()
	return res
}
