package dao

import (
	"Webook/account/domain"
	"gorm.io/gorm"
	"time"
)

func InitTables(db *gorm.DB) error {
	err := db.AutoMigrate(&Account{}, &AccountActivity{})
	if err != nil {
		return err
	}
	// 补充一个初始化系统账号的代码
	now := time.Now().UnixMilli()
	_ = db.Create(&Account{
		//没有系统账号就行
		Type:     domain.AccountTypeSystem,
		Currency: "CNY",
		Ctime:    now,
		Utime:    now,
	}).Error
	return nil
}
