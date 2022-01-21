package producer

import (
	"github.com/Shopify/sarama"
	"github.com/anildogan/audit-log-to-kafka/kafka-log-processor/constant"
	"github.com/spf13/viper"
	"strings"
	"time"
)

func NewProducer() (sarama.SyncProducer, error) {
	v, err := sarama.ParseKafkaVersion(viper.GetString(constant.KafkaVersion))
	if err != nil {
		return nil, err
	}

	config := sarama.NewConfig()
	config.Version = v
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Idempotent = true
	config.Net.MaxOpenRequests = 1

	config.Metadata.Retry.Max = 1
	config.Metadata.Retry.Backoff = 5 * time.Second

	syncProducer, err := sarama.NewSyncProducer(strings.Split(viper.GetString(constant.KafkaBrokers), ","), config)
	if err != nil {
		return nil, err
	}
	return syncProducer, err
}
