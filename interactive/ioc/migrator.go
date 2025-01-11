package ioc

import (
	"Webook/interactive/repository/dao"
	"Webook/pkg/ginx"
	"Webook/pkg/gormx/connpool"
	"Webook/pkg/logger"
	"Webook/pkg/migrator/event"
	"Webook/pkg/migrator/event/fixer"
	"Webook/pkg/migrator/scheduler"
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
)

const topic = "migrator_interactives"

func InitFixDataConsumer(l logger.Logger, src SrcDB, dst DstDB, client sarama.Client) *fixer.Consumer[dao.Interactive] {
	consumer, err := fixer.NewConsumer[dao.Interactive](client, l, src, dst, topic)
	if err != nil {
		panic(err)
	}
	return consumer
}

func InitMigratorProducer(p sarama.SyncProducer) event.Producer {
	return event.NewSaramaProducer(p, topic)
}
func InitMigratorWeb(l logger.Logger, src SrcDB, dst DstDB,
	pool *connpool.DoubleWritePool, producer event.Producer) *ginx.Server {
	engine := gin.Default()
	ginx.InitCounter(prometheus.CounterOpts{
		Namespace: "hwoss",
		Subsystem: "webook_intr",
		Name:      "http_biz_code",
		ConstLabels: map[string]string{
			"instance_id": "my-instance -1 ",
		},
	})

	intrs := scheduler.NewScheduler[dao.Interactive](l, src, dst, pool, producer)
	intrs.RegisterRoutes(engine.Group("/intr"))
	addr := viper.GetString("migrator.http.addr")
	return &ginx.Server{
		Engine: engine,
		Addr:   addr,
	}
}

func InitDoubleWritePool(src SrcDB, dst DstDB) *connpool.DoubleWritePool {
	return connpool.NewDoubleWritePool(src, dst)
}
