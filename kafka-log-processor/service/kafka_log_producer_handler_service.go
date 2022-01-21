package service

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/anildogan/audit-log-to-kafka/kafka-log-processor/model"
	"go.uber.org/zap"
)

type IKafkaLogProducerHandlerService interface {
	Handle(message *sarama.ConsumerMessage) error
}

type kafkaLogProducerHandlerService struct {
	kafkaProcessedLogPublisher IKafkaProcessedLogPublisher
}

func NewKafkaLogProducerHandlerService(KafkaProcessedLogPublisher IKafkaProcessedLogPublisher) *kafkaLogProducerHandlerService {
	return &kafkaLogProducerHandlerService{
		kafkaProcessedLogPublisher: KafkaProcessedLogPublisher,
	}
}

func (s *kafkaLogProducerHandlerService) Handle(consumedMessage *sarama.ConsumerMessage) error {
	if consumedMessage.Value == nil {
		return nil
	}

	var rawEvent model.KafkaLogMessage
	err := json.Unmarshal(consumedMessage.Value, &rawEvent)
	if err != nil {
		return err
	}

	topicName := rawEvent.Log.KafkaLog.TopicName
	key := rawEvent.Log.KafkaLog.Key
	rawMessage := rawEvent.Log.KafkaLog.KafkaMessage

	zap.S().Infof("Sending consumedMessage to topic: %s from container: %s, with key: %s, with consumedMessage: %s",
		topicName,
		rawEvent.ContainerName,
		key,
		string(rawMessage))

	err = s.kafkaProcessedLogPublisher.Publish(key, topicName, rawMessage)
	if err != nil {
		return err
	}

	return nil
}
