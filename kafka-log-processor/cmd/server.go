package cmd

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/anildogan/audit-log-to-kafka/kafka-log-processor/constant"
	"github.com/anildogan/audit-log-to-kafka/kafka-log-processor/consumer"
	kafka "github.com/anildogan/audit-log-to-kafka/kafka-log-processor/kafka/consumer"
	"github.com/anildogan/audit-log-to-kafka/kafka-log-processor/kafka/producer"
	"github.com/anildogan/audit-log-to-kafka/kafka-log-processor/service"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
	"time"
)

type Server struct {
	Router *gin.Engine
}

func NewServer() *Server {
	return &Server{}
}

func (srv *Server) InitRouter() *Server {
	producer, _ := producer.NewProducer()
	kafkaProcessedLogPublisher := service.NewKafkaProcessedLogPublisher(producer)
	handlerService := service.NewKafkaLogProducerHandlerService(kafkaProcessedLogPublisher)

	runKafkaLogProducerConsumer(handlerService, producer)

	return srv
}

func runKafkaLogProducerConsumer(handlerService service.IKafkaLogProducerHandlerService, producer sarama.SyncProducer) {
	topic := viper.GetString(constant.KafkaLogProducerConsumerTopic)
	errorTopic := viper.GetString(constant.KafkaLogProducerConsumerErrorTopic)
	kafkaLogProducerConsumerTopicConfig := kafka.ConnectionParameters{
		ConsumerGroupId: viper.GetString(constant.KafkaConsumerGroupId),
		ClientId:        viper.GetString(constant.KafkaClientId),
		Brokers:         viper.GetString(constant.KafkaBrokers),
		Version:         viper.GetString(constant.KafkaVersion),
		Topic:           []string{topic},
		ErrorTopic:      errorTopic,
		FromBeginning:   viper.GetBool(constant.KafkaFromBeginning),
	}

	kafkaConsumer, err := kafka.NewConsumer(kafkaLogProducerConsumerTopicConfig)
	if err != nil {
		panic(err)
	}

	eventHandler := consumer.NewEventHandler(handlerService, producer, errorTopic)
	kafkaConsumer.Subscribe(eventHandler)

	zap.S().Info("Kafka log producer consumer is running successfully.")
}

func (srv *Server) Run() *Server {
	srv.Router.Run()
	return srv
}

func (srv *Server) InitLogger() *Server {
	logger, err := zap.Config{
		Level:       zap.NewAtomicLevelAt(getLogLevel()),
		Encoding:    "json",
		OutputPaths: []string{"stdout"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:    "message",
			LevelKey:      "level",
			TimeKey:       "time",
			CallerKey:     "source",
			StacktraceKey: "stack-trace",
			EncodeLevel:   zapcore.CapitalLevelEncoder,
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format("2006-01-02 15:04:05,000"))
			},
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}.Build()

	if err != nil {
		panic(fmt.Sprintf("error : %s , message : logger can not initialize", err.Error()))
	}

	zap.ReplaceGlobals(logger)

	defer logger.Sync()

	return srv
}

func getLogLevel() zapcore.Level {
	logLevel := strings.ToLower(viper.GetString(constant.LogLevel))
	switch logLevel {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}
