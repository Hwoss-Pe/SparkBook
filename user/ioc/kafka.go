package ioc

import (
	userEvents "Webook/user/events"
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

func InitKafka() sarama.SyncProducer {
	type Config struct {
		Addrs []string `yaml:"addrs"`
	}
	var cfg Config
	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}
	sc := sarama.NewConfig()
	sc.Producer.Return.Successes = true
	p, err := sarama.NewSyncProducer(cfg.Addrs, sc)
	if err != nil {
		panic(err)
	}
	return p
}

func InitProducer(p sarama.SyncProducer) userEvents.Producer {
	return userEvents.NewSaramaSyncProducer(p)
}
