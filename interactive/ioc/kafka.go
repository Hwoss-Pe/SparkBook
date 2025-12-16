package ioc

import (
	"Webook/interactive/events"
	"Webook/interactive/repository/dao"
	"Webook/pkg/migrator/event/fixer"
	"Webook/pkg/saramax"

	"time"

	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

func InitKafka() sarama.Client {
	type Config struct {
		Addrs []string `yaml:"addrs"`
	}
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	// 确保新加入的消费者在没有已提交位点时从最早开始消费，避免错过历史消息
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = time.Second
	var cfg Config
	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}
	client, err := sarama.NewClient(cfg.Addrs, config)
	if err != nil {
		panic(err)
	}
	return client
}

func InitProducer(client sarama.Client) sarama.SyncProducer {
	res, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	return res
}

func NewConsumers(intrs *events.InteractiveReadEventConsumer,
	like *events.LikeEventConsumer,
	collect *events.CollectEventConsumer,
	notif *events.NotificationEventConsumer,
	fix *fixer.Consumer[dao.Interactive]) []saramax.Consumer {
	return []saramax.Consumer{
		intrs,
		like,
		collect,
		notif,
		fix,
	}
}
