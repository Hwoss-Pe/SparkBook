package ioc

import (
	"Webook/pkg/logger"
	"Webook/user/repository/dao"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

func InitDB(l logger.Logger) *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	c := Config{}
	err := viper.UnmarshalKey("db", &c)
	if err != nil {
		panic(fmt.Errorf("初始化配置失败 %v1, 原因 %w", c, err))
	}
	//在gorm接入自己拓展的日志
	db, err := gorm.Open(mysql.Open(c.DSN), &gorm.Config{
		Logger: glogger.New(gormLoggerFunc(l.Debug),
			glogger.Config{
				SlowThreshold: 0,
				LogLevel:      glogger.Info,
			}),
	})
	if err != nil {
		panic(err)
	}

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}

	//后面接入prometheus

	return db
}

type gormLoggerFunc func(msg string, fields ...logger.Field)

func (g gormLoggerFunc) Printf(msg string, args ...interface{}) {
	g(msg, logger.Field{Key: "args", Value: args})
}
