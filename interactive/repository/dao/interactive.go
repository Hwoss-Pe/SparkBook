package dao

import (
	"Webook/pkg/migrator"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

//go:generate mockgen -source=./interactive.go -package=daomocks -destination=mocks/interactive.mock.go InteractiveDAO
type InteractiveDAO interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	InsertLikeInfo(ctx context.Context, biz string, bizId, uid int64) error
	GetLikeInfo(ctx context.Context, biz string, bizId, uid int64) (UserLikeBiz, error)
	DeleteLikeInfo(ctx context.Context, biz string, bizId, uid int64) error
	Get(ctx context.Context, biz string, bizId int64) (Interactive, error)
	InsertCollectionBiz(ctx context.Context, cb UserCollectionBiz) error
	GetCollectionInfo(ctx context.Context, biz string, bizId, uid int64) (UserCollectionBiz, error)
	BatchIncrReadCnt(ctx context.Context, bizs []string, ids []int64) error
	GetByIds(ctx context.Context, biz string, ids []int64) ([]Interactive, error)
}

type GORMInteractiveDAO struct {
	db *gorm.DB
}

func (G *GORMInteractiveDAO) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	return G.incrReadCnt(G.db.WithContext(ctx), biz, bizId)
}

// InsertLikeInfo 点赞对于用户可见的操作用的是status字段
func (G *GORMInteractiveDAO) InsertLikeInfo(ctx context.Context, biz string, bizId, uid int64) error {
	now := time.Now().UnixMilli()
	err := G.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{
				"status": 1,
				"utime":  now,
			}),
		}).Create(&UserLikeBiz{
			Uid:    uid,
			Ctime:  now,
			Utime:  now,
			Biz:    biz,
			BizId:  bizId,
			Status: 1,
		}).Error
		if err != nil {
			return err
		}
		return tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{
				"like_cnt": gorm.Expr("`like_cnt`+1"),
				"utime":    now,
			}),
		}).Create(&Interactive{
			LikeCnt: 1,
			Ctime:   now,
			Utime:   now,
			Biz:     biz,
			BizId:   bizId,
		}).Error
	})

	return err
}

func (G *GORMInteractiveDAO) GetLikeInfo(ctx context.Context, biz string, bizId, uid int64) (UserLikeBiz, error) {
	var res UserLikeBiz
	err := G.db.WithContext(ctx).Where("biz=? AND biz_id = ? AND uid = ? AND status = ?",
		biz, bizId, uid, 1).First(&res).Error
	return res, err
}

func (G *GORMInteractiveDAO) DeleteLikeInfo(ctx context.Context, biz string, bizId, uid int64) error {
	now := time.Now().UnixMilli()
	err := G.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&UserLikeBiz{}).Where("biz =? AND biz_id = ? AND uid = ?", biz, bizId, uid).Updates(map[string]any{
			"status": 0,
			"utime":  now,
		}).Error
		if err != nil {
			return err
		}
		return G.db.WithContext(ctx).Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{
				"like_cnt": gorm.Expr("`like_cnt` - 1 "),
				"utime":    now,
			}),
		}).Create(&Interactive{
			LikeCnt: 1,
			Ctime:   now,
			Utime:   now,
			Biz:     biz,
			BizId:   bizId,
		}).Error
	})
	return err
}

func (G *GORMInteractiveDAO) Get(ctx context.Context, biz string, bizId int64) (Interactive, error) {
	var res Interactive
	err := G.db.WithContext(ctx).Where("biz = ? AND biz_id = ?", biz, bizId).
		First(&res).Error
	return res, err
}

func (G *GORMInteractiveDAO) InsertCollectionBiz(ctx context.Context, cb UserCollectionBiz) error {
	now := time.Now().UnixMilli()
	cb.Utime = now
	cb.Ctime = now
	return G.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := G.db.WithContext(ctx).Create(&cb).Error
		if err != nil {
			return err
		}
		return tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{
				"collect_cnt": gorm.Expr("`collect_cnt`+1"),
				"utime":       now,
			}),
		}).Create(&Interactive{
			CollectCnt: 1,
			Ctime:      now,
			Utime:      now,
			Biz:        cb.Biz,
			BizId:      cb.BizId,
		}).Error
	})
}

func (G *GORMInteractiveDAO) GetCollectionInfo(ctx context.Context, biz string, bizId, uid int64) (UserCollectionBiz, error) {
	var res UserCollectionBiz
	err := G.db.WithContext(ctx).Where("biz=? AND biz_id = ? AND uid = ?", biz, bizId, uid).
		First(&res).Error
	return res, err
}

func (G *GORMInteractiveDAO) BatchIncrReadCnt(ctx context.Context, bizs []string, ids []int64) error {
	return G.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 让调用者保证两者是相等的
		for i := 0; i < len(bizs); i++ {
			err := G.IncrReadCnt(ctx, bizs[i], ids[i])
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (G *GORMInteractiveDAO) GetByIds(ctx context.Context, biz string, ids []int64) ([]Interactive, error) {
	var res []Interactive
	err := G.db.WithContext(ctx).Where("biz = ? and id in ?", biz, ids).Find(&res).Error
	return res, err
}

func (G *GORMInteractiveDAO) incrReadCnt(tx *gorm.DB, biz string, bizId int64) error {
	now := time.Now().UnixMilli()
	return tx.Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]any{
			"utime":    now,
			"read_cnt": gorm.Expr("read_cnt +1 "),
		}),
	}).Create(&Interactive{
		ReadCnt: 1,
		Ctime:   now,
		Utime:   now,
		Biz:     biz,
		BizId:   bizId,
	}).Error
}

type Interactive struct {
	Id         int64  `gorm:"primaryKey,autoIncrement"`
	BizId      int64  `gorm:"uniqueIndex:biz_type_id"`
	Biz        string `gorm:"type:varchar(128);uniqueIndex:biz_type_id"`
	ReadCnt    int64
	CollectCnt int64
	LikeCnt    int64
	Ctime      int64
	Utime      int64
}

func (i Interactive) ID() int64 {
	return i.Id
}

func (i Interactive) CompareTo(dst migrator.Entity) bool {
	if di, ok := dst.(Interactive); ok {
		return di == i
	}
	return false
}

type UserLikeBiz struct {
	Id     int64  `gorm:"primaryKey,autoIncrement"`
	BizId  int64  `gorm:"uniqueIndex:biz_type_id_uid"`
	Biz    string `gorm:"type:varchar(128);uniqueIndex:biz_type_id_uid"`
	Uid    int64  `gorm:"uniqueIndex:biz_type_id_uid"`
	Status uint8
	Ctime  int64
	Utime  int64
}

// Collection 收藏夹
type Collection struct {
	Id   int64  `gorm:"primaryKey,autoIncrement"`
	Name string `gorm:"type=varchar(1024)"`
	Uid  int64  `gorm:""`

	Ctime int64
	Utime int64
}

type UserCollectionBiz struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 收藏夹 ID
	// 作为关联关系中的外键，我们这里需要索引
	Cid   int64  `gorm:"index"`
	BizId int64  `gorm:"uniqueIndex:biz_type_id_uid"`
	Biz   string `gorm:"type:varchar(128);uniqueIndex:biz_type_id_uid"`
	// 这算是一个冗余，因为正常来说，
	// 只需要在 Collection 中维持住 Uid 就可以
	Uid   int64 `gorm:"uniqueIndex:biz_type_id_uid"`
	Ctime int64
	Utime int64
}
