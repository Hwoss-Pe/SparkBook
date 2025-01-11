package startup

import (
	"Webook/interactive/repository/dao"
	"database/sql"
	"golang.org/x/net/context"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

var db *gorm.DB

func InitTestDB() *gorm.DB {
	if db == nil {
		dsn := "root:root@tcp(localhost:13316)/webook"
		sqlDB, err := sql.Open("mysql", dsn)
		if err != nil {
			panic(err)
		}
		for {
			ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
			err := sqlDB.PingContext(ctx)
			cancelFunc()
			if err == nil {
				break
			}
			log.Println("等待数据库连接")
		}
		db, err = gorm.Open(mysql.Open(dsn))
		if err != nil {
			panic(err)
		}
		err = dao.InitTables(db)
		if err != nil {
			panic(err)
		}

	}
	return db
}
