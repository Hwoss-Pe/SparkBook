package ioc

import (
	"Webook/pkg/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

func InitLogger() logger.Logger {
	// 这里   用一个小技巧，
	// 就是直接使用 zap 本身的配置结构体来处理
	cfg := zap.NewDevelopmentConfig()
	logDir := os.Getenv("WEBOOK_LOG_DIR")
	if logDir == "" {
		logDir = "./logs"
	}
	logFile := filepath.Join(logDir, "bff.log")
	_ = os.MkdirAll(filepath.Dir(logFile), 0o755)
	cfg.OutputPaths = []string{"stdout", logFile}
	cfg.ErrorOutputPaths = []string{"stderr", logFile}
	err := viper.UnmarshalKey("log", &cfg)
	if err != nil {
		panic(err)
	}
	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	// Make zap.L() and stdlib log write to the same sink.
	zap.ReplaceGlobals(l)
	_ = zap.RedirectStdLog(l)

	res := logger.NewZapLogger(l)
	res.Info("bff logger initialized", logger.String("log_file", logFile))
	return res
}
