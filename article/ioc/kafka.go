package ioc

import (
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

func InitProducer() sarama.SyncProducer {
	type Config struct {
		Addrs []string `yaml:"addrs"`
	}
	newConfig := sarama.NewConfig()
	newConfig.Producer.Return.Successes = true
	var cfg Config
	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}
	producer, err := sarama.NewSyncProducer(cfg.Addrs, newConfig)
	if err != nil {
		panic(err)
	}
	return producer
}
