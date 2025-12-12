package main

import (
	artdao "Webook/article/repository/dao"
	searchdao "Webook/search/repository/dao"
	tagdao "Webook/tag/repository/dao"
	userdao "Webook/user/repository/dao"
	"context"
	"encoding/json"
	"flag"
	"log"
	"strconv"

	"github.com/olivere/elastic/v7"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	var (
		userDSN    string
		articleDSN string
		tagDSN     string
		esURL      string
		batch      int
	)
	flag.StringVar(&userDSN, "user_dsn", "root:root@tcp(localhost:13316)/webook_user", "")
	flag.StringVar(&articleDSN, "article_dsn", "root:root@tcp(localhost:13316)/webook_article", "")
	flag.StringVar(&tagDSN, "tag_dsn", "root:root@tcp(localhost:13316)/webook_article", "")
	flag.StringVar(&esURL, "es", "http://localhost:9200", "")
	flag.IntVar(&batch, "batch", 500, "")
	flag.Parse()

	ec, err := elastic.NewClient(
		elastic.SetURL(esURL),
		elastic.SetSniff(false),
	)
	if err != nil {
		log.Fatal(err)
	}
	if err = searchdao.InitES(ec); err != nil {
		log.Fatal(err)
	}
	userES := searchdao.NewUserElasticDAO(ec)
	artES := searchdao.NewArticleElasticDAO(ec)
	anyES := searchdao.NewAnyESDAO(ec)

	udb, err := gorm.Open(mysql.Open(userDSN), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	adb, err := gorm.Open(mysql.Open(articleDSN), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	tdb, err := gorm.Open(mysql.Open(tagDSN), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	if err = backfillUsers(context.Background(), udb, userES, batch); err != nil {
		log.Fatal(err)
	}
	if err = backfillArticles(context.Background(), adb, artES, batch); err != nil {
		log.Fatal(err)
	}
	if err = backfillTags(context.Background(), tdb, anyES, batch); err != nil {
		log.Fatal(err)
	}
}

func backfillUsers(ctx context.Context, db *gorm.DB, es searchdao.UserDAO, batch int) error {
	offset := 0
	for {
		var rows []userdao.User
		err := db.WithContext(ctx).Model(&userdao.User{}).Limit(batch).Offset(offset).Find(&rows).Error
		if err != nil {
			return err
		}
		if len(rows) == 0 {
			break
		}
		for _, u := range rows {
			du := searchdao.User{
				Id:       u.Id,
				Email:    u.Email.String,
				Nickname: u.Nickname.String,
				Phone:    u.Phone.String,
				Avatar:   u.Avatar.String,
			}
			if err = es.InputUser(ctx, du); err != nil {
				return err
			}
		}
		offset += len(rows)
	}
	return nil
}

func backfillArticles(ctx context.Context, db *gorm.DB, es searchdao.ArticleDAO, batch int) error {
	offset := 0
	for {
		var rows []artdao.PublishedArticle
		err := db.WithContext(ctx).Model(&artdao.PublishedArticle{}).Order("utime ASC").Limit(batch).Offset(offset).Find(&rows).Error
		if err != nil {
			return err
		}
		if len(rows) == 0 {
			break
		}
		for _, a := range rows {
			da := searchdao.Article{
				Id:      a.Id,
				Title:   a.Title,
				Status:  int32(a.Status),
				Content: a.Content,
				Tags:    nil,
			}
			if err = es.InputArticle(ctx, da); err != nil {
				return err
			}
		}
		offset += len(rows)
	}
	return nil
}

type tagKey struct {
	Uid   int64
	Biz   string
	BizId int64
}

func backfillTags(ctx context.Context, db *gorm.DB, es searchdao.AnyDAO, batch int) error {
	var keys []tagKey
	err := db.WithContext(ctx).Model(&tagdao.TagBiz{}).
		Select("uid, biz, biz_id").
		Distinct("uid, biz, biz_id").
		Find(&keys).Error
	if err != nil {
		return err
	}
	tdao := tagdao.NewGORMTagDAO(db)
	for i := 0; i < len(keys); i += batch {
		j := i + batch
		if j > len(keys) {
			j = len(keys)
		}
		for _, k := range keys[i:j] {
			tags, err1 := tdao.GetTagsByBiz(ctx, k.Uid, k.Biz, k.BizId)
			if err1 != nil {
				return err1
			}
			names := make([]string, 0, len(tags))
			for _, t := range tags {
				names = append(names, t.Name)
			}
			payload := struct {
				Uid   int64    `json:"uid"`
				Biz   string   `json:"biz"`
				BizId int64    `json:"biz_id"`
				Tags  []string `json:"tags"`
			}{Uid: k.Uid, Biz: k.Biz, BizId: k.BizId, Tags: names}
			data, err1 := json.Marshal(payload)
			if err1 != nil {
				return err1
			}
			docID := strconv.FormatInt(k.Uid, 10) + "_" + k.Biz + "_" + strconv.FormatInt(k.BizId, 10)
			if err1 = es.Input(ctx, searchdao.TagIndexName, docID, string(data)); err1 != nil {
				return err1
			}
		}
	}
	return nil
}
