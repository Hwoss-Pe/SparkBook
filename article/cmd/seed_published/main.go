package main

import (
	"Webook/article/repository/dao"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func main() {
	var (
		startID     = flag.Int("start_id", 200, "starting article id")
		startAuthor = flag.Int("start_author_id", 20, "starting author id")
		count       = flag.Int("count", 200, "number of rows to seed")
		//coverImage  = flag.String("cover_image", "https://images.unsplash.com/photo-1555066931-4365d14bab8c?w=800&h=400", "cover image url")
		dsnFlag = flag.String("dsn", "root:root@tcp(localhost:13316)/webook_article", "mysql dsn, e.g. user:pass@tcp(host:port)/db")
	)
	flag.Parse()

	dsn := *dsnFlag
	if dsn == "" {
		if env := os.Getenv("SEED_DSN"); env != "" {
			dsn = env
		}
	}
	if dsn == "" {
		type Config struct {
			DSN string `yaml:"dsn"`
		}
		var cfg Config
		_ = viper.UnmarshalKey("article.db", &cfg)
		if cfg.DSN == "" {
			_ = viper.UnmarshalKey("db", &cfg)
		}
		dsn = cfg.DSN
	}
	if dsn == "" {
		fmt.Println("missing DSN, pass via --dsn or env SEED_DSN or viper config (article.db.dsn / db.dsn)")
		os.Exit(2)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("open db error:", err.Error())
		os.Exit(2)
	}

	now := time.Now()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	topics := []string{
		"Go 并发实践", "微服务架构", "Elasticsearch 向量检索", "Kafka 消费者位点",
		"Redis 缓存策略", "RAG 问答", "gRPC 服务治理", "数据库迁移双写",
		"ETCD 服务发现", "消息驱动设计", "限流与熔断", "可观测性与指标",
	}

	for i := 0; i < *count; i++ {
		id := int64(*startID + i)
		// Random authorID between 101 and 150
		authorID := int64(r.Intn(50) + 101)
		topic := topics[i%len(topics)]
		title := fmt.Sprintf("%s 实战 #%d", topic, id)
		content := fmt.Sprintf("%s 深入讲解与案例集合。编号-%d。包含关键词：%s、Best Practices、Patterns。", topic, id, topic)
		shift := time.Duration(r.Intn(720)) * time.Minute
		ctime := now.Add(-shift).UnixMilli()
		utime := ctime + int64(r.Intn(3600*1000))
		// Use local cover images from 1 to 200
		coverImg := fmt.Sprintf("/static/covers/test/cover_%d.jpg", (i%200)+1)

		row := dao.PublishedArticle{
			Id:         id,
			Title:      title,
			Content:    content,
			AuthorId:   authorID,
			Status:     2,
			Ctime:      ctime,
			Utime:      utime,
			CoverImage: coverImg,
		}

		err = db.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"title":       row.Title,
				"content":     row.Content,
				"author_id":   row.AuthorId,
				"status":      row.Status,
				"utime":       row.Utime,
				"cover_image": row.CoverImage,
			}),
		}).Create(&row).Error
		if err != nil {
			fmt.Println("insert error on id", id, ":", err.Error())
			os.Exit(2)
		}
	}

	fmt.Printf("seeded %d rows into published_articles starting id=%d, author_id=%d\n", *count, *startID, *startAuthor)
}
