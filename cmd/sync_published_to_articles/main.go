package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Article struct {
	Id         int64  `gorm:"primaryKey,autoIncrement"`
	Title      string `gorm:"type=varchar(4096)"`
	Content    string `gorm:"type:BLOB"`
	CoverImage string `gorm:"type=varchar(1024)"`
	AuthorId   int64  `gorm:"index"`
	Status     uint8
	Ctime      int64
	Utime      int64 `gorm:"index"`
}

type PublishedArticle struct {
	Id         int64  `gorm:"primaryKey,autoIncrement"`
	Title      string `gorm:"type=varchar(4096)"`
	Content    string `gorm:"type:BLOB"`
	CoverImage string `gorm:"type=varchar(1024)"`
	AuthorId   int64  `gorm:"index"`
	Status     uint8
	Ctime      int64
	Utime      int64
}

func main() {
	dsn := "root:root@tcp(localhost:13316)/webook_article?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("open db error: %v", err))
	}
	const batch = 500
	offset := 0
	for {
		var rows []PublishedArticle
		err := db.Model(&PublishedArticle{}).
			Limit(batch).Offset(offset).
			Find(&rows).Error
		if err != nil {
			panic(fmt.Sprintf("query published_articles error: %v", err))
		}
		if len(rows) == 0 {
			break
		}
		for _, r := range rows {
			a := Article{
				Id:         r.Id,
				Title:      r.Title,
				Content:    r.Content,
				CoverImage: r.CoverImage,
				AuthorId:   r.AuthorId,
				Status:     2,
				Ctime:      r.Ctime,
				Utime:      r.Utime,
			}
			err = db.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "id"}},
				DoUpdates: clause.Assignments(map[string]interface{}{
					"title":       a.Title,
					"content":     a.Content,
					"cover_image": a.CoverImage,
					"author_id":   a.AuthorId,
					"status":      a.Status,
					"utime":       a.Utime,
				}),
			}).Create(&a).Error
			if err != nil {
				panic(fmt.Sprintf("upsert articles id=%d error: %v", a.Id, err))
			}
		}
		offset += len(rows)
	}
	fmt.Println("sync done")
}
