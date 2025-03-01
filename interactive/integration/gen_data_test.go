package integration

import (
	"Webook/interactive/integration/startup"
	"Webook/interactive/repository/dao"
	_ "embed"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

//go:embed init.sql
var initSQL string

func TestGenSQL(t *testing.T) {
	file, err := os.OpenFile("data.sql", os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0666)
	require.NoError(t, err)
	defer file.Close()

	//	创建数据库和表语句
	_, err = file.WriteString(initSQL)
	require.NoError(t, err)
	const prefix = "INSERT INTO `interactives`(`biz_id`, `biz`, `read_cnt`, `collect_cnt`, `like_cnt`, `ctime`, `utime`)\nVALUES"
	const rowNum = 10
	now := time.Now().UnixMilli()
	_, err = file.WriteString(prefix)
	for i := 0; i < rowNum; i++ {
		if i > 0 {
			file.Write([]byte{',', '\n'})
		}
		file.Write([]byte{'('})
		// biz_id
		file.WriteString(strconv.Itoa(i + 1))
		// biz
		file.WriteString(`,"test",`)
		// read_cnt
		file.WriteString(strconv.Itoa(int(rand.Int31n(10000))))
		file.Write([]byte{','})

		// collect_cnt
		file.WriteString(strconv.Itoa(int(rand.Int31n(10000))))
		file.Write([]byte{','})
		// like_cnt
		file.WriteString(strconv.Itoa(int(rand.Int31n(10000))))
		file.Write([]byte{','})

		// ctime
		file.WriteString(strconv.FormatInt(now, 10))
		file.Write([]byte{','})

		// utime
		file.WriteString(strconv.FormatInt(now, 10))

		file.Write([]byte{')'})
	}
}
func TestGenData(t *testing.T) {
	// 这个是批量插入，数据量不是特别大的时候，可以用这个
	// GenData 要比 GenSQL 慢
	db := startup.InitTestDB()
	db.DryRun = true
	// 1000 批
	for i := 0; i < 10; i++ {
		// 每次 100 条
		//    可以考虑直接用 CreateInBatches，GORM 帮   分批次
		// 自己分可以控制内存消耗
		const batchSize = 100
		data := make([]dao.Interactive, 0, batchSize)
		now := time.Now().UnixMilli()
		for j := 0; j < batchSize; j++ {
			data = append(data, dao.Interactive{
				Biz:        "test",
				BizId:      int64(i*batchSize + j + 1),
				ReadCnt:    rand.Int63(),
				LikeCnt:    rand.Int63(),
				CollectCnt: rand.Int63(),
				Utime:      now,
				Ctime:      now,
			})
		}

		err := db.Transaction(func(tx *gorm.DB) error {
			return tx.Create(data).Error
		})
		require.NoError(t, err)
	}
}
