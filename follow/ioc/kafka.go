package ioc

import (
	"Webook/follow/events"
	"Webook/pkg/saramax"
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

func InitKafka() sarama.Client {
	type Config struct {
		Addrs []string `yaml:"addrs"`
	}
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	var c Config
	err := viper.UnmarshalKey("kafka", &c)
	if err != nil {
		panic(err)
	}
	cli, err := sarama.NewClient(c.Addrs, cfg)
	if err != nil {
		panic(err)
	}
	return cli
}

func InitProducer(c sarama.Client) sarama.SyncProducer {
	p, err := sarama.NewSyncProducerFromClient(c)
	if err != nil {
		panic(err)
	}
	return p
}

func NewConsumers(follow *events.FollowEventConsumer) []saramax.Consumer {
	return []saramax.Consumer{follow}
}
