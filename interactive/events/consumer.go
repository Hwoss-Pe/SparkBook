package events

import (
	"Webook/article/repository"
	"Webook/pkg/canalx"
	"Webook/pkg/logger"
	"Webook/pkg/migrator"
	"Webook/pkg/migrator/event"
	"Webook/pkg/migrator/validator"
	"Webook/pkg/saramax"
	"github.com/IBM/sarama"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"sync/atomic"
	"time"
)

type MYSQLBinlogConsumer[T migrator.Entity] struct {
	client   sarama.Client
	l        logger.Logger
	table    string
	repo     *repository.CachedArticleRepository
	srcToDst *validator.CanalIncrValidator[T]
	dstToSrc *validator.CanalIncrValidator[T]
	desFirst *atomic.Bool
}

func (r *MYSQLBinlogConsumer[T]) Start() error {
	client, err := sarama.NewConsumerGroupFromClient("migrator", r.client)
	if err != nil {
		return err
	}
	go func() {
		err2 := client.Consume(context.Background(), []string{"webook_binlog"}, saramax.NewHandler[canalx.Message[T]](r.l, r.Consume))
		if err2 != nil {
			r.l.Error("退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}

// Consume canal对应的增量同步发送给kafka
func (r *MYSQLBinlogConsumer[T]) Consume(msg *sarama.ConsumerMessage,
	val canalx.Message[T]) error {
	dstFirst := r.desFirst.Load()
	var v *validator.CanalIncrValidator[T]
	// db:
	//  src:
	//    dsn: "root:root@tcp(localhost:13316)/webook"
	//  dst:
	//    dsn: "root:root@tcp(localhost:13316)/webook_intr"
	if dstFirst && val.Database == "webook_intr" {
		// 校验，用 dst 的来校验
		v = r.dstToSrc
	} else if !dstFirst && val.Database == "webook" {
		v = r.srcToDst
	}
	if v != nil {
		for _, data := range val.Data {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			err := v.Validate(ctx, data.ID())
			cancel()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// NewMYSQLBinlogConsumer 数据迁移传入对应的*gorm.DB ，其余代码载pkg里面
func NewMYSQLBinlogConsumer[T migrator.Entity](client sarama.Client, l logger.Logger,
	table string, repo *repository.CachedArticleRepository,
	src *gorm.DB,
	dst *gorm.DB,
	producer event.Producer) *MYSQLBinlogConsumer[T] {

	srcToDst := validator.NewCanalIncrValidator[T](src, dst, "SRC", l, producer)
	dstToSrc := validator.NewCanalIncrValidator[T](src, dst, "DST", l, producer)

	return &MYSQLBinlogConsumer[T]{
		client: client, l: l, table: table, repo: repo,
		srcToDst: srcToDst,
		dstToSrc: dstToSrc,
		desFirst: &atomic.Bool{}}
}
func (r *MYSQLBinlogConsumer[T]) DstFirst() {
	r.desFirst.Store(true)
}
