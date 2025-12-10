package dao

import (
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type GORMFollowRelationDAO struct {
	db *gorm.DB
}

func (g *GORMFollowRelationDAO) FollowRelationList(ctx context.Context, follower, offset, limit int64) ([]FollowRelation, error) {
	var res []FollowRelation
	err := g.db.WithContext(ctx).Where("follower = ? and status = ? ", follower, FollowRelationStatusActive).
		Offset(int(offset)).Limit(int(limit)).Find(&res).Error
	return res, err
}
func (g *GORMFollowRelationDAO) FolloweeRelationList(ctx context.Context, followee, offset, limit int64) ([]FollowRelation, error) {
	var res []FollowRelation
	err := g.db.WithContext(ctx).Where("followee = ? and status = ? ", followee, FollowRelationStatusActive).
		Offset(int(offset)).Limit(int(limit)).Find(&res).Error
	return res, err
}

func (g *GORMFollowRelationDAO) FollowerRelationDetail(ctx context.Context, follower int64, followee int64) (FollowRelation, error) {
	var res FollowRelation
	err := g.db.WithContext(ctx).Where("follower = ? AND followee = ? AND status = ?",
		follower, followee, FollowRelationStatusActive).First(&res).Error
	return res, err
}

func (g *GORMFollowRelationDAO) CreateFollowRelation(ctx context.Context, f FollowRelation) error {
	now := time.Now().UnixMilli()
	f.Utime = now
	f.Ctime = now
	f.Status = FollowRelationStatusActive
	return g.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]any{
			"utime":  now,
			"status": FollowRelationStatusActive,
		}),
	}).Create(&f).Error
}

func (g *GORMFollowRelationDAO) UpdateStatus(ctx context.Context, followee int64, follower int64, status uint8) error {
	now := time.Now().UnixMilli()
	return g.db.WithContext(ctx).Model(&FollowRelation{}).Where("follower = ? and followee = ?", follower, followee).
		Updates(map[string]any{
			"status": status,
			"utime":  now,
		}).Error
}

func (g *GORMFollowRelationDAO) CntFollower(ctx context.Context, uid int64) (int64, error) {
	var res int64
	err := g.db.WithContext(ctx).Model(&FollowRelation{}).
		Where("followee = ? and status = ? ", uid, FollowRelationStatusActive).
		Count(&res).Error
	return res, err
}

func (g *GORMFollowRelationDAO) CntFollowee(ctx context.Context, uid int64) (int64, error) {
	var res int64
	err := g.db.WithContext(ctx).Model(&FollowRelation{}).
		Where("follower = ? and status = ? ", uid, FollowRelationStatusActive).
		Count(&res).Error
	return res, err
}

func NewGORMFollowRelationDAO(db *gorm.DB) FollowRelationDao {
	return &GORMFollowRelationDAO{
		db: db,
	}
}
