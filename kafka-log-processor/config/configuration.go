package config

import (
	"github.com/anildogan/audit-log-to-kafka/kafka-log-processor/constant"
	"github.com/spf13/viper"
)

func Init() {
	viper.SetDefault(constant.LogLevel, "DEBUG")
	viper.SetDefault(constant.KafkaBrokers, "localhost:9092")
	viper.SetDefault(constant.KafkaVersion, "6.0.0")
	viper.SetDefault(constant.KafkaClientId, "kafka-log-processor-client")
	viper.SetDefault(constant.KafkaProducerRetryMax, 3)
	viper.SetDefault(constant.KafkaConsumerGroupId, "kafka-log-processor")
	viper.SetDefault(constant.KafkaFromBeginning, false)
	viper.SetDefault(constant.KafkaLogProducerConsumerTopic, "application.log")
	viper.SetDefault(constant.KafkaLogProducerConsumerRetryTopic, "application.log.RETRY")
	viper.SetDefault(constant.KafkaLogProducerConsumerErrorTopic, "application.log.ERROR")
	viper.AutomaticEnv()
}
