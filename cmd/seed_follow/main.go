package main

import (
	"Webook/follow/repository/dao"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"math/rand"
	"time"
)

func main() {
	dsn := "root:root@tcp(localhost:13316)/webook_follow?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %v", err))
	}
	err = dao.InitTables(db)
	if err != nil {
		panic(fmt.Sprintf("failed to migrate tables: %v", err))
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	start := 100
	end := 199
	now := time.Now().UnixMilli()

	type pair struct {
		a int64
		b int64
	}
	created := make(map[pair]bool)

	for u := start; u <= end; u++ {
		k := 5 + r.Intn(6)
		for i := 0; i < k; i++ {
			v := start + r.Intn(end-start+1)
			if v == u {
				continue
			}
			a := int64(u)
			b := int64(v)
			p := pair{a: min64(a, b), b: max64(a, b)}
			if created[p] {
				continue
			}
			err = upsertFollow(db, a, b, now)
			if err != nil {
				fmt.Printf("failed to follow %d -> %d: %v\n", a, b, err)
				continue
			}
			err = upsertFollow(db, b, a, now)
			if err != nil {
				fmt.Printf("failed to follow %d -> %d: %v\n", b, a, err)
				continue
			}
			created[p] = true
		}
	}
	fmt.Println("Seeding mutual follows completed.")
}

func upsertFollow(db *gorm.DB, follower, followee int64, now int64) error {
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "follower"}, {Name: "followee"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"status": dao.FollowRelationStatusActive,
			"utime":  now,
		}),
	}).Create(&dao.FollowRelation{
		Follower: follower,
		Followee: followee,
		Status:   dao.FollowRelationStatusActive,
		Ctime:    now,
		Utime:    now,
	}).Error
}

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
