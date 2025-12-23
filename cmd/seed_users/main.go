package main

import (
	"Webook/user/repository/dao"
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

func main() {
	// DSN for webook_user database
	dsn := "root:root@tcp(localhost:13316)/webook_user?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %v", err))
	}

	// Ensure table exists
	err = db.AutoMigrate(&dao.User{})
	if err != nil {
		panic(fmt.Sprintf("failed to migrate table: %v", err))
	}

	fmt.Println("Start seeding users...")
	for i := 1; i <= 100; i++ {
		now := time.Now().UnixMilli()

		// Use FirstOrCreate to avoid duplicate errors if run multiple times
		var u dao.User
		email := fmt.Sprintf("user_%d@example.com", i)
		phone := fmt.Sprintf("15000000%03d", i) // 15000000001 - 15000000100

		// Search criteria
		err := db.Where("email = ?", email).First(&u).Error
		if err == nil {
			// User exists, update avatar just in case
			u.Avatar = sql.NullString{String: fmt.Sprintf("/static/avatars/test/avatar_%d.png", i), Valid: true}
			u.AboutMe = sql.NullString{String: fmt.Sprintf("我是用户%d，热爱Go、微服务、RAG等技术。", i), Valid: true}
			u.Birthday = sql.NullInt64{Int64: randomBirthdayMillis(i), Valid: true}
			db.Save(&u)
			fmt.Printf("Updated user %d\n", i)
			continue
		}

		u = dao.User{
			Email:    sql.NullString{String: email, Valid: true},
			Password: "password",
			Phone:    sql.NullString{String: phone, Valid: true},
			Nickname: sql.NullString{String: fmt.Sprintf("User_%d", i), Valid: true},
			AboutMe:  sql.NullString{String: fmt.Sprintf("我是用户%d，热爱Go、微服务、RAG等技术。", i), Valid: true},
			Birthday: sql.NullInt64{Int64: randomBirthdayMillis(i), Valid: true},
			Avatar:   sql.NullString{String: fmt.Sprintf("/static/avatars/test/avatar_%d.png", i), Valid: true},
			Ctime:    now,
			Utime:    now,
		}

		err = db.Create(&u).Error
		if err != nil {
			fmt.Printf("Failed to create user %d: %v\n", i, err)
		} else {
			fmt.Printf("Created user %d\n", i)
		}
	}
	fmt.Println("Seeding completed.")
}

func randomBirthdayMillis(seed int) int64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano() + int64(seed)))
	year := 1990 + r.Intn(16)
	month := time.Month(1 + r.Intn(12))
	day := 1 + r.Intn(28)
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local).UnixMilli()
}
