package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var ErrUserDuplicate = errors.New("用户邮箱或者手机号冲突")

// ErrDataNotFound 通用的数据没找到
var ErrDataNotFound = gorm.ErrRecordNotFound

//go:generate mockgen -source=./user.go -package=daomocks -destination=mocks/user.mock.go UserDAO
type UserDAO interface {
	Insert(ctx context.Context, u User) error
	UpdateNonZeroFields(ctx context.Context, u User) error
	FindByPhone(ctx context.Context, phone string) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByWechat(ctx context.Context, openId string) (User, error)
	FindById(ctx context.Context, id int64) (User, error)
}

type GORMUserDAO struct {
	db *gorm.DB
}

func (G *GORMUserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Utime = now
	u.Ctime = now
	err := G.db.WithContext(ctx).Create(&u).Error
	var me *mysql.MySQLError
	if errors.As(err, &me) {
		const uniqueIndexErrNo uint16 = 1062
		if me.Number == uniqueIndexErrNo {
			return ErrUserDuplicate
		}
	}
	return err
}

func (G *GORMUserDAO) UpdateNonZeroFields(ctx context.Context, u User) error {
	//这里直接传对象引用的会只修改非零值或者非空字符串
	return G.db.Updates(&u).Error
}

func (G *GORMUserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := G.db.WithContext(ctx).First(&u, "phone = ? ", phone).Error
	return u, err
}

func (G *GORMUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := G.db.WithContext(ctx).First(&u, "email = ? ", email).Error
	return u, err
}

func (G *GORMUserDAO) FindByWechat(ctx context.Context, openId string) (User, error) {
	var u User
	err := G.db.WithContext(ctx).First(&u, "wechat_open_id = ?", openId).Error
	return u, err
}

func (G *GORMUserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := G.db.WithContext(ctx).First(&u, "id = ?", id).Error
	return u, err
}

func NewGORMUserDAO(db *gorm.DB) UserDAO {
	return &GORMUserDAO{db: db}

}

type User struct {
	Id       int64          `gorm:"primaryKey,autoIncrement"`
	Email    sql.NullString `gorm:"unique"`
	Password string
	Phone    sql.NullString `gorm:"unique"`
	Birthday sql.NullInt64
	Nickname sql.NullString `gorm:"type=varchar(128)"`
	AboutMe  sql.NullString `gorm:"type=varchar(1024)"`
	Avatar   sql.NullString `gorm:"type=varchar(512)"` // 头像URL

	WechatOpenId  sql.NullString `gorm:"type=varchar(256);unique"`
	WechatUnionId sql.NullString `gorm:"type=varchar(256)"`
	Ctime         int64
	Utime         int64
}
