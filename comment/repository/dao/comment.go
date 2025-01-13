package dao

import (
	"database/sql"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

// ErrDataNotFound 通用的数据没找到
var ErrDataNotFound = gorm.ErrRecordNotFound

//go:generate mockgen -source=./comment.go -package=daomocks -destination=mocks/comment.mock.go CommentDAO
type CommentDAO interface {
	Insert(ctx context.Context, u Comment) error
	FindByBiz(ctx context.Context, biz string, bizId, minID, limit int64) ([]Comment, error)
	FindCommentList(ctx context.Context, comment Comment) ([]Comment, error)

	FindRepliesByPid(ctx context.Context, pid int64, offset, limit int) ([]Comment, error)

	Delete(ctx context.Context, u Comment) error
	FindOneByIDs(ctx context.Context, id []int64) ([]Comment, error)
	FindRepliesByRid(ctx context.Context, rid int64, id int64, limit int64) ([]Comment, error)
}

type Comment struct {
	Id int64 `json:"id" gorm:"column:id;primaryKey"`
	//发表评论的用户
	Uid int64  `json:"uid" gorm:"column:uid;index"`
	Biz string `gorm:"column:biz;index:biz_type_id" json:"biz"`
	// 对应的业务ID
	BizID int64 `gorm:"column:biz_id;index:biz_type_id" json:"bizID"`

	//	根评论为0就是一级
	RootID sql.NullInt64 `json:"rootID" gorm:"column:root_id;index"`
	// 父级评论
	PID sql.NullInt64 `gorm:"column:pid;index" json:"pid"`
	//设置级联删除，删除的时候会把父评论对应的给一起删了,可以理解成PID找不到的id的时候就会把自己也删掉
	ParentComment *Comment `gorm:"ForeignKey:PID;AssociationForeignKey:ID;constraint:OnDelete:CASCADE"`
	// 评论内容
	Content string `gorm:"type:text;column:content" json:"content"`
	// 创建时间
	Ctime int64 `gorm:"column:ctime;" json:"ctime"`
	// 更新时间
	Utime int64 `gorm:"column:utime;" json:"utime"`
}

func (*Comment) TableName() string {
	return "comments"
}

type GORMCommentDAO struct {
	db *gorm.DB
}

func NewGORMCommentDAO(db *gorm.DB) CommentDAO {
	return &GORMCommentDAO{db: db}
}

func (G *GORMCommentDAO) Insert(ctx context.Context, u Comment) error {
	return G.db.
		WithContext(ctx).
		Create(u).
		Error
}

func (G *GORMCommentDAO) FindByBiz(ctx context.Context, biz string, bizId, minID, limit int64) ([]Comment, error) {
	var res []Comment
	err := G.db.WithContext(ctx).
		Where("biz = ? AND biz_id = ? AND id > ? AND pid IS NULL", biz, bizId, minID).
		Limit(int(limit)).
		Find(&res).Error
	return res, err
}

func (G *GORMCommentDAO) FindCommentList(ctx context.Context, u Comment) ([]Comment, error) {
	var res []Comment
	builder := G.db.WithContext(ctx)
	//根据id是否为0，如果只返回主评论，否则返回主评论和其所有子评论
	if u.Id == 0 {
		builder = builder.
			Where("biz=?", u.Biz).
			Where("biz_id=?", u.BizID).
			Where("root_id is null")
	} else {
		builder = builder.Where("root_id=? or id =?", u.Id, u.Id)
	}
	err := builder.Find(&res).Error
	return res, err
}

func (G *GORMCommentDAO) FindRepliesByPid(ctx context.Context, pid int64, offset, limit int) ([]Comment, error) {
	var res []Comment
	err := G.db.WithContext(ctx).Where("pid = ?", pid).
		Order("id DESC").
		Offset(offset).Limit(limit).Find(&res).Error
	return res, err
}

func (G *GORMCommentDAO) Delete(ctx context.Context, u Comment) error {
	return G.db.WithContext(ctx).Delete(&Comment{
		Id: u.Id,
	}).Error
}

func (G *GORMCommentDAO) FindOneByIDs(ctx context.Context, ids []int64) ([]Comment, error) {
	var res []Comment
	err := G.db.WithContext(ctx).
		Where("id in ?", ids).
		First(&res).
		Error
	return res, err
}

func (G *GORMCommentDAO) FindRepliesByRid(ctx context.Context, rid int64, id int64, limit int64) ([]Comment, error) {
	var res []Comment
	err := G.db.WithContext(ctx).
		Where("root_id = ? AND id < ?", rid, id).
		Order("id ASC").
		Limit(int(limit)).Find(&res).Error
	return res, err
}
