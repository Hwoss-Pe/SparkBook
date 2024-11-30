package dao

import (
	"context"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func TestGORMTagDAO_GetTagsByBiz(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("gorm.db?mode=memory"), &gorm.Config{
		// 只输出 SQL，不执行查询
		DryRun: true,
	})
	require.NoError(t, err)
	db = db.Debug()
	dao := NewGORMTagDAO(db)
	res, err := dao.GetTagsByBiz(context.Background(), 123, "test", 456)
	if err != nil {
		return
	}
	t.Log(res)
}
