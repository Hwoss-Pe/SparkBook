package dao

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"gorm.io/gorm"
	"time"
)

type Tag struct {
	Id    int64  `gorm:"primaryKey,autoIncrement"`
	Name  string `gorm:"type=varchar(4096)"`
	Uid   int64  `gorm:"index"`
	Ctime int64
	Utime int64
}
type TagBiz struct {
	Id    int64  `gorm:"primaryKey,autoIncrement"`
	BizId int64  `gorm:"index:biz_type_id"`
	Biz   string `gorm:"index:biz_type_id"`
	// 冗余字段，加快查询和删除
	Uid int64 `gorm:"index"`
	Tid int64
	//这个是外键，gorm的操作
	//表示把Tag表里面的Id和当前表的Tid关联，并且设置Tag表id删除的时候对应的tagBiz的数据也会被删除，级联删除
	Tag   *Tag  `gorm:"ForeignKey:Tid;AssociationForeignKey:Id;constraint:OnDelete:CASCADE"`
	Ctime int64 `bson:"ctime,omitempty"`
	Utime int64 `bson:"utime,omitempty"`
}
type TagDAO interface {
	CreateTag(ctx context.Context, tag Tag) (int64, error)
	CreateTagBiz(ctx context.Context, tagBiz []TagBiz) error
	GetTagsByUid(ctx context.Context, uid int64) ([]Tag, error)
	GetTagsByBiz(ctx context.Context, uid int64, biz string, bizId int64) ([]Tag, error)
	GetTags(ctx context.Context, offset, limit int) ([]Tag, error)
	GetTagsById(ctx context.Context, ids []int64) ([]Tag, error)
}

type GORMTagDAO struct {
	db *gorm.DB
}

func NewGORMTagDAO(db *gorm.DB) TagDAO {
	return &GORMTagDAO{db: db}
}

func (G *GORMTagDAO) CreateTag(ctx context.Context, tag Tag) (int64, error) {
	now := time.Now().UnixMilli()
	tag.Ctime = now
	tag.Utime = now
	err := G.db.WithContext(ctx).Create(&tag).Error
	return tag.Id, err
}

// CreateTagBiz 对资源打上标签,并且采用覆盖机制，保证标签顺序
func (G *GORMTagDAO) CreateTagBiz(ctx context.Context, tagBiz []TagBiz) error {
	if len(tagBiz) == 0 {
		return nil
	}
	now := time.Now().UnixMilli()
	for _, t := range tagBiz {
		t.Ctime = now
		t.Utime = now
	}
	//这一堆记录的uid和对应资源都是一样的,所以我第一个解决
	first := tagBiz[0]
	return G.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&TagBiz{}).Delete("uid = ? and biz = ? and biz_id = ?", first.Uid, first.Biz, first.BizId).Error
		if err != nil {
			return err
		}
		return tx.Create(&tagBiz).Error
	})
}

func (G *GORMTagDAO) GetTagsByUid(ctx context.Context, uid int64) ([]Tag, error) {
	var res []Tag
	err := G.db.WithContext(ctx).Where("uid= ?", uid).Find(&res).Error
	return res, err
}

func (G *GORMTagDAO) GetTagsByBiz(ctx context.Context, uid int64, biz string, bizId int64) ([]Tag, error) {
	var res []TagBiz
	err := G.db.WithContext(ctx).Model(&TagBiz{}).
		InnerJoins("Tag", G.db.Model(&Tag{})).
		Where("Tag.uid = ? AND biz = ? AND biz_id = ?", uid, biz, bizId).Find(&res).Error
	return slice.Map(res, func(idx int, src TagBiz) Tag {
		return *src.Tag
	}), err
}

func (G *GORMTagDAO) GetTags(ctx context.Context, offset, limit int) ([]Tag, error) {
	var res []Tag
	err := G.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&res).Error
	return res, err
}

func (G *GORMTagDAO) GetTagsById(ctx context.Context, ids []int64) ([]Tag, error) {
	var res []Tag
	err := G.db.WithContext(ctx).Where("id IN ?", ids).Find(&res).Error
	return res, err
}
