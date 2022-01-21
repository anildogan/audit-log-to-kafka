package service

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	"os"
	"time"
)

type IKafkaProcessedLogPublisher interface {
	Publish(key string, topicName string, rawMessage json.RawMessage ) error
}

type KafkaProcessedLogPublisher struct {
	Sarama sarama.SyncProducer
}

func NewKafkaProcessedLogPublisher(Sarama sarama.SyncProducer) *KafkaProcessedLogPublisher {
	return &KafkaProcessedLogPublisher{
		Sarama: Sarama,
	}
}

func (p *KafkaProcessedLogPublisher) Publish(key string, topicName string, rawMessage json.RawMessage) error {

	message := &sarama.ProducerMessage{
		Topic:     topicName,
		Value:     sarama.StringEncoder(rawMessage),
		Timestamp: time.Now(),
		Key:       sarama.StringEncoder(key),
	}

	partition, offset, err := p.Sarama.SendMessage(message)

	if err != nil {
		zap.S().Error(os.Stderr, "Error occurred when publish Audit Log Event to topic: ", err)
		return err
	}

	zap.S().Infof("Audit Log Event published. Topic: %s Partition:  %d, Offset:  %d", topicName, partition, offset)
	return nil
}