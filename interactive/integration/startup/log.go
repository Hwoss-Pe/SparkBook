package startup

import (
	"Webook/pkg/logger"
)

func InitLog() logger.Logger {
	return logger.NewNoOpLogger()
}
