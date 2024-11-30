package ioc

import (
	"Webook/tag/events"
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

func InitKafka() sarama.Client {
	type Config struct {
		Addrs []string `yaml:"addrs"`
	}
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.Return.Successes = true
	var cfg Config
	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}
	client, err := sarama.NewClient(cfg.Addrs, saramaCfg)
	if err != nil {
		panic(err)
	}

	return client
}

func InitProducer(c sarama.Client) events.Producer {
	client, _ := sarama.NewSyncProducerFromClient(c)
	return events.NewSaramaSyncProducer(client)
}
