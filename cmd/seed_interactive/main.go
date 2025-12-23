package main

import (
	"Webook/interactive/repository/dao"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"math/rand"
	"time"
)

func main() {
	// DSN for webook_intr database
	dsn := "root:root@tcp(localhost:13316)/webook_intr?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %v", err))
	}

	// Ensure tables exist
	err = db.AutoMigrate(&dao.Interactive{}, &dao.UserLikeBiz{}, &dao.UserCollectionBiz{}, &dao.Collection{})
	if err != nil {
		panic(fmt.Sprintf("failed to migrate tables: %v", err))
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	startID := 200
	count := 200 // Assuming we have 200 articles starting from 200
	userStart := 100
	userEnd := 200

	fmt.Println("Start seeding interactive data...")

	// 1. Ensure each user has a default collection
	userCollectionMap := make(map[int64]int64) // uid -> cid
	for uid := userStart; uid <= userEnd; uid++ {
		var c dao.Collection
		err := db.Where("uid = ? AND name = ?", uid, "My Favorites").First(&c).Error
		if err == gorm.ErrRecordNotFound {
			c = dao.Collection{
				Uid:   int64(uid),
				Name:  "My Favorites",
				Ctime: time.Now().UnixMilli(),
				Utime: time.Now().UnixMilli(),
			}
			if err := db.Create(&c).Error; err != nil {
				fmt.Printf("Failed to create collection for user %d: %v\n", uid, err)
				continue
			}
		} else if err != nil {
			fmt.Printf("Error checking collection for user %d: %v\n", uid, err)
			continue
		}
		userCollectionMap[int64(uid)] = c.Id
	}
	fmt.Printf("Ensured collections for users %d-%d\n", userStart, userEnd)

	// 2. Generate interactive data for articles
	for i := 0; i < count; i++ {
		bizId := int64(startID + i)
		biz := "article"
		now := time.Now().UnixMilli()

		// Random interaction counts
		likeCount := r.Intn(50)    // 0-49 likes
		collectCount := r.Intn(20) // 0-19 collections
		readCount := int64(likeCount + collectCount + r.Intn(1000) + 1)

		// Generate unique users for likes
		likedUsers := randPerm(userEnd-userStart+1, likeCount)
		for _, idx := range likedUsers {
			uid := int64(userStart + idx)
			err := db.Clauses(clause.OnConflict{
				DoUpdates: clause.Assignments(map[string]interface{}{
					"status": 1,
					"utime":  now,
				}),
			}).Create(&dao.UserLikeBiz{
				Uid:    uid,
				Biz:    biz,
				BizId:  bizId,
				Status: 1,
				Ctime:  now,
				Utime:  now,
			}).Error
			if err != nil {
				fmt.Printf("Failed to insert like for bizId %d, uid %d: %v\n", bizId, uid, err)
			}
		}

		// Generate unique users for collections
		collectedUsers := randPerm(userEnd-userStart+1, collectCount)
		for _, idx := range collectedUsers {
			uid := int64(userStart + idx)
			cid, ok := userCollectionMap[uid]
			if !ok {
				continue
			}
			err := db.Clauses(clause.OnConflict{
				DoUpdates: clause.Assignments(map[string]interface{}{
					"utime": now,
				}),
			}).Create(&dao.UserCollectionBiz{
				Uid:   uid,
				Biz:   biz,
				BizId: bizId,
				Cid:   cid,
				Ctime: now,
				Utime: now,
			}).Error
			if err != nil {
				fmt.Printf("Failed to insert collection for bizId %d, uid %d: %v\n", bizId, uid, err)
			}
		}

		// Update or Insert Interactive Aggregation
		err = db.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "biz_id"}, {Name: "biz"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"read_cnt":    readCount,
				"like_cnt":    int64(likeCount),
				"collect_cnt": int64(collectCount),
				"utime":       now,
			}),
		}).Create(&dao.Interactive{
			BizId:      bizId,
			Biz:        biz,
			ReadCnt:    readCount,
			LikeCnt:    int64(likeCount),
			CollectCnt: int64(collectCount),
			Ctime:      now,
			Utime:      now,
		}).Error

		if err != nil {
			fmt.Printf("Failed to upsert interactive for bizId %d: %v\n", bizId, err)
		} else {
			// fmt.Printf("Processed interactive for bizId %d\n", bizId)
		}
	}

	fmt.Println("Seeding interactive data completed.")
}

// randPerm generates n unique random numbers in [0, max-1]
func randPerm(max, n int) []int {
	if n > max {
		n = max
	}
	p := rand.Perm(max)
	return p[:n]
}
