package repository

import (
	"Webook/user/domain"
	"Webook/user/repository/cache"
	"Webook/user/repository/dao"
	"context"
	"database/sql"
	"errors"
	"time"
)

var ErrUserDuplicate = dao.ErrUserDuplicate
var ErrUserNotFound = dao.ErrDataNotFound

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	// Update 更新数据，只有非 0 值才会更新
	Update(ctx context.Context, u domain.User) error
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindById(ctx context.Context, id int64) (domain.User, error)
	// FindByWechat 暂时可以认为按照 openId来查询
	// 将来可能需要按照 unionId 来查询
	FindByWechat(ctx context.Context, openId string) (domain.User, error)
}

type CachedUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func NewCachedUserRepository(dao dao.UserDAO, cache cache.UserCache) UserRepository {
	return &CachedUserRepository{dao: dao, cache: cache}
}

func (c *CachedUserRepository) Create(ctx context.Context, u domain.User) error {
	return c.dao.Insert(ctx, dao.User{
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Password: u.Password,
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		WechatOpenId: sql.NullString{
			String: u.WechatInfo.OpenId,
			Valid:  u.WechatInfo.OpenId != "",
		},
		WechatUnionId: sql.NullString{
			String: u.WechatInfo.UnionId,
			Valid:  u.WechatInfo.UnionId != "",
		},
	})
}

func (c *CachedUserRepository) Update(ctx context.Context, u domain.User) error {
	err := c.dao.UpdateNonZeroFields(ctx, c.domainToEntity(u))
	if err != nil {
		return err
	}
	//更新那些非敏感字段需要更新缓存，这里折中的方法是直接删缓存，等查的时候重新写入
	return c.cache.Delete(ctx, u.Id)
}

func (c *CachedUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := c.dao.FindByPhone(ctx, phone)
	return c.entityToDomain(u), err
}

func (c *CachedUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := c.dao.FindByEmail(ctx, email)
	return c.entityToDomain(u), err
}

func (c *CachedUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	//这里采用查数据库后写进缓存，采用降级的快慢路径操作
	// 先走缓存，缓存找不到就从数据库找，找完重新set进去
	user, err := c.cache.Get(ctx, id)
	switch {
	case err == nil:
		return user, nil
		//	这里返回的是redis不存在这个key，那我就去数据库找，找完写进缓存
	case errors.Is(err, cache.ErrKeyNotExist):
		if ctx.Value("downgrade") == "true" {
			return domain.User{}, errors.New("缓存中没有数据，并且触发了降级，放弃查询数据库")
		}
		ur, err := c.dao.FindById(ctx, id)
		if err != nil {
			return domain.User{}, err
		}
		du := c.entityToDomain(ur)
		_ = c.cache.Set(ctx, du)
		return du, nil
	default:
		//这里如果redis出现错误，我就不让他流量直接打到数据库，因此直接返回错误
		return domain.User{}, err
	}
}

func (c *CachedUserRepository) FindByWechat(ctx context.Context, openId string) (domain.User, error) {
	u, err := c.dao.FindByPhone(ctx, openId)
	return c.entityToDomain(u), err
}

func (c *CachedUserRepository) domainToEntity(user domain.User) dao.User {
	return dao.User{
		Id: user.Id,
		Email: sql.NullString{
			String: user.Email,
			Valid:  user.Email != "",
		},
		Password: user.Password,
		Phone: sql.NullString{
			String: user.Phone,
			Valid:  user.Phone != "",
		},
		Birthday: sql.NullInt64{
			Int64: user.Birthday.UnixMilli(),
			Valid: !user.Birthday.IsZero(),
		},
		Nickname: sql.NullString{
			String: user.Nickname,
			Valid:  user.Nickname != "",
		},
		AboutMe: sql.NullString{
			String: user.AboutMe,
			Valid:  user.AboutMe != "",
		},
	}
}

func (c *CachedUserRepository) entityToDomain(ue dao.User) domain.User {
	var birthday time.Time
	//在这里进行一个字段判断合法
	if ue.Birthday.Valid {
		birthday = time.UnixMilli(ue.Birthday.Int64)
	}

	return domain.User{
		Id:       ue.Id,
		Email:    ue.Email.String,
		Password: ue.Password,
		Phone:    ue.Phone.String,
		Nickname: ue.Nickname.String,
		AboutMe:  ue.AboutMe.String,
		Birthday: birthday,
		Ctime:    time.UnixMilli(ue.Ctime),
		WechatInfo: domain.WechatInfo{
			OpenId:  ue.WechatOpenId.String,
			UnionId: ue.WechatUnionId.String,
		},
	}
}
