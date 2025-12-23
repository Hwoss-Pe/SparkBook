package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:root@tcp(localhost:13316)/webook_user?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("connect db error: %v", err))
	}
	var cnt int64
	err = db.Raw("SELECT COUNT(*) FROM users WHERE id BETWEEN ? AND ?", 100, 199).Scan(&cnt).Error
	if err != nil {
		panic(fmt.Sprintf("check target range error: %v", err))
	}
	if cnt > 0 {
		panic(fmt.Sprintf("target IDs 100-199 already exist, count=%d", cnt))
	}
	tx := db.Begin()
	if tx.Error != nil {
		panic(tx.Error)
	}
	err = tx.Exec("UPDATE users SET id = id + ? WHERE id BETWEEN ? AND ?", 100000, 203, 302).Error
	if err != nil {
		_ = tx.Rollback()
		panic(fmt.Sprintf("phase1 update error: %v", err))
	}
	err = tx.Exec("UPDATE users SET id = id - ? WHERE id BETWEEN ? AND ?", 100103, 100203, 100302).Error
	if err != nil {
		_ = tx.Rollback()
		panic(fmt.Sprintf("phase2 update error: %v", err))
	}
	if err = tx.Commit().Error; err != nil {
		panic(fmt.Sprintf("commit error: %v", err))
	}
	fmt.Println("Reindex users done: moved IDs 203-302 to 100-199")
}
