package consumer

import (
	"context"
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	"math"
	"os"
	"strings"
	"time"
)

type ConnectionParameters struct {
	Brokers         string
	ClientId        string
	Version         string
	Topic           []string
	ErrorTopic      string
	ConsumerGroupId string
	KafkaUsername   string
	KafkaPassword   string
	ApplicationPort string
	FromBeginning   bool
}

type Consumer interface {
	Subscribe(handler EventHandler)
	Unsubscribe()
}

type EventHandler interface {
	Setup(sarama.ConsumerGroupSession) error
	Cleanup(sarama.ConsumerGroupSession) error
	ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error
}

type kafkaConsumer struct {
	topic         []string
	errorTopic    string
	consumerGroup sarama.ConsumerGroup
}

func NewConsumer(connectionParams ConnectionParameters) (Consumer, error) {
	v, err := sarama.ParseKafkaVersion(connectionParams.Version)
	if err != nil {
		return nil, err
	}

	config := sarama.NewConfig()
	config.Version = v
	config.Consumer.Return.Errors = true
	config.ClientID = connectionParams.ClientId
	if connectionParams.FromBeginning {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}
	config.Metadata.Retry.Max = math.MaxInt32
	config.Metadata.Retry.Backoff = 5 * time.Second
	config.Consumer.Group.Session.Timeout = 20 * time.Second
	config.Consumer.Group.Heartbeat.Interval = 6 * time.Second
	config.Consumer.MaxProcessingTime = 500 * time.Millisecond

	cg, err := sarama.NewConsumerGroup(strings.Split(connectionParams.Brokers, ","), connectionParams.ConsumerGroupId, config)
	if err != nil {
		return nil, err
	}

	return &kafkaConsumer{
		topic:         connectionParams.Topic,
		errorTopic:    connectionParams.ErrorTopic,
		consumerGroup: cg,
	}, nil
}

const ErrorKey = "ErrorMessage"
const RetryCount = "RetryCount"

func (c *kafkaConsumer) Subscribe(handler EventHandler) {
	ctx := context.Background()
	topics := func() []string {
		result := make([]string, 0)
		if c.errorTopic != "" {
			result = append(result, c.errorTopic)
		}
		result = append(result, c.topic...)
		return result
	}

	go func() {
		for {
			if err := c.consumerGroup.Consume(ctx, topics(), handler); err != nil {
				zap.S().Errorf("Error from consumer: %+v", err.Error())
				os.Exit(1)
				return
			}

			if ctx.Err() != nil {
				zap.S().Errorf("Error from consumer: %+v", ctx.Err().Error())
				return
			}
		}
	}()
	go func() {
		for err := range c.consumerGroup.Errors() {
			zap.S().Errorf("Error from consumer group : %+v", err.Error())
		}
	}()
	zap.S().Infof("Kafka  consumer listens topic : %+v", c.topic)
}

func (c *kafkaConsumer) Unsubscribe() {
	if err := c.consumerGroup.Close(); err != nil {
		zap.S().Errorf("Client wasn't closed :%+v", err)
	}
	zap.S().Info("Kafka consumer closed")
}
