package ioc

import (
	"Webook/interactive/events"
	"Webook/interactive/repository/dao"
	"Webook/pkg/migrator/event/fixer"
	"Webook/pkg/saramax"
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

func InitKafka() sarama.Client {
	type Config struct {
		Addrs []string `yaml:"addrs"`
	}
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
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
	fix *fixer.Consumer[dao.Interactive]) []saramax.Consumer {
	return []saramax.Consumer{
		intrs,
		fix,
	}
}
