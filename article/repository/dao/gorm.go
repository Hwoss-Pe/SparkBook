package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type GORMArticleDAO struct {
	db *gorm.DB
}

func (G *GORMArticleDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	err := G.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}

func (G *GORMArticleDAO) UpdateById(ctx context.Context, art Article) error {
	now := time.Now().UnixMilli()
	res := G.db.Model(&Article{}).WithContext(ctx).
		Where("id=? AND author_id = ? ", art.Id, art.AuthorId).
		Updates(map[string]any{
			"title":       art.Title,
			"content":     art.Content,
			"cover_image": art.CoverImage,
			"status":      art.Status,
			"utime":       now,
		})
	err := res.Error
	if err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return errors.New("更新数据失败")
	}
	return nil
}

func (G *GORMArticleDAO) GetByAuthor(ctx context.Context, author int64, offset, limit int) ([]Article, error) {
	var arts []Article
	err := G.db.WithContext(ctx).Model(&Article{}).
		Where("author_id = ?", author).
		Offset(offset).
		Limit(limit).
		Order("utime DESC").
		Find(&arts).Error
	return arts, err
}

func (G *GORMArticleDAO) GetById(ctx context.Context, id int64) (Article, error) {
	var art Article
	err := G.db.WithContext(ctx).Model(&Article{}).
		Where("id = ?", id).
		First(&art).Error
	return art, err
}

func (G *GORMArticleDAO) GetPubById(ctx context.Context, id int64) (PublishedArticle, error) {
	var art PublishedArticle
	err := G.db.WithContext(ctx).Model(&PublishedArticle{}).
		Where("id = ?", id).
		First(&art).Error
	return art, err
}

// Sync 为发表
func (G *GORMArticleDAO) Sync(ctx context.Context, art Article) (int64, error) {
	tx := G.db.WithContext(ctx).Begin()
	now := time.Now().UnixMilli()
	defer tx.Rollback()
	txDAO := NewGORMArticleDAO(tx)
	var (
		id  = art.Id
		err error
	)
	if id == 0 {
		id, err = txDAO.Insert(ctx, art)
	} else {
		err = txDAO.UpdateById(ctx, art)
	}
	if err != nil {
		return 0, err
	}
	//	处理线上库
	art.Id = id
	//偷懒的转移
	publishArt := PublishedArticle(art)
	publishArt.Utime = now
	publishArt.Ctime = now
	err = tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{
				Name: "id",
			},
		},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"title":       art.Title,
			"content":     art.Content,
			"cover_image": art.CoverImage,
			"status":      art.Status,
			"utime":       now,
		}),
	}).Create(&publishArt).Error
	if err != nil {
		return 0, err
	}

	tx.Commit()
	return id, tx.Error
}

func (G *GORMArticleDAO) SyncStatus(ctx context.Context, author, id int64, status uint8) error {
	return G.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&Article{}).
			Where("id = ? and author_id = ?", id, author).
			Update("status", status)

		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return ErrPossibleIncorrectAuthor
		}
		//	线上库也要进行同步
		res = tx.Model(&PublishedArticle{}).
			Where("id = ? and author_id = ?", id, author).
			Update("status", status)

		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return ErrPossibleIncorrectAuthor
		}
		return nil
	})

}

func (G *GORMArticleDAO) ListPubByUtime(ctx context.Context, utime time.Time, offset int, limit int) ([]PublishedArticle, error) {
	var res []PublishedArticle
	err := G.db.WithContext(ctx).
		Order("utime DESC").
		Where("utime < ? ", utime.UnixMilli()).
		Limit(limit).Offset(offset).Find(&res).Error
	return res, err
}

func NewGORMArticleDAO(db *gorm.DB) ArticleDAO {
	return &GORMArticleDAO{
		db: db,
	}
}
