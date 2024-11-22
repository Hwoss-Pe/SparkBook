package connpool

import (
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

type DoubleWriteSuite struct {
	suite.Suite
	db *gorm.DB
}

func (d *DoubleWriteSuite) SetupSuite() {
	t := d.T()
	src, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	require.NoError(t, err)
	err = src.AutoMigrate(&Interactive{})
	require.NoError(t, err)
	dst, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook_intr"))
	require.NoError(t, err)
	err = dst.AutoMigrate(&Interactive{})
	require.NoError(t, err)
	//需要这么代表双写的数据库
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: &DoubleWritePool{
			src:     src.ConnPool,
			dst:     dst.ConnPool,
			pattern: atomicx.NewValueOf(PatternSrcFirst),
		},
	}))
	require.NoError(t, err)
	d.db = db
}
func (d *DoubleWriteSuite) TearDownTest() {
	//d.db.Exec("TRUNCATE TABLE interactives")
}
func (d *DoubleWriteSuite) TestDoubleWriteTransaction() {
	err := d.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&Interactive{
			Biz:   "test",
			BizId: 10087,
		}).Error
	})
	require.NoError(d.T(), err)
}
func (d *DoubleWriteSuite) TestDoubleWriteTest() {
	t := d.T()
	err := d.db.Create(&Interactive{
		Biz:   "test",
		BizId: 10086,
	}).Error
	assert.NoError(t, err)
	// 查询数据库就可以看到对应的数据
}
func TestDoubleWrite(t *testing.T) {
	suite.Run(t, new(DoubleWriteSuite))
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
