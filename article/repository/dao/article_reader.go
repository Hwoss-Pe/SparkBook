package dao

import (
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ArticleReaderDAO interface {
	// Upsert 将会更新标题和内容，但是不会更新别的内容
	// 这个要求 Reader 和 Author 是不同库
	Upsert(ctx context.Context, art Article) error
	// UpsertV2 版本用于同库不同表
	UpsertV2(ctx context.Context, art PublishedArticle) error
}

type GORMArticleReaderDAO struct {
	db *gorm.DB
}

func NewGORMArticleReaderDAO(db *gorm.DB) ArticleReaderDAO {
	return &GORMArticleReaderDAO{
		db: db,
	}
}

func (dao *GORMArticleReaderDAO) Upsert(ctx context.Context, art Article) error {
	return dao.db.Clauses(clause.OnConflict{
		// ID 冲突的时候。实际上，在 MYSQL 里面   写不写都可以
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"title":       art.Title,
			"content":     art.Content,
			"cover_image": art.CoverImage,
		}),
	}).Create(&art).Error
}

// UpsertV2 同库不同表
func (dao *GORMArticleReaderDAO) UpsertV2(ctx context.Context, art PublishedArticle) error {
	return dao.db.Clauses(clause.OnConflict{
		// ID 冲突的时候。实际上，在 MYSQL 里面写不写都可以
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"title":       art.Title,
			"content":     art.Content,
			"cover_image": art.CoverImage,
		}),
	}).Create(&art).Error
}
