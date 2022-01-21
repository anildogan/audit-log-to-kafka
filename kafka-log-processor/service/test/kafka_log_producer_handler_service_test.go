package test

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/anildogan/audit-log-to-kafka/kafka-log-processor/model"
	"github.com/anildogan/audit-log-to-kafka/kafka-log-processor/service"
	"github.com/anildogan/audit-log-to-kafka/kafka-log-processor/service/mock"
	_ "github.com/golang/mock/gomock"
	"testing"
)

func TestKafkaLogProducerHandlerService_Handle(t *testing.T) {
	//given
	kafkaLogMessage := new(model.KafkaLogMessage)
	kafkaLogMessage.ContainerName = "container-name"
	kafkaLogMessage.Log.KafkaLog.TopicName = "topic-name"
	kafkaLogMessage.Log.KafkaLog.Key = "123"
	kafkaLogMessage.Log.KafkaLog.KafkaMessage = json.RawMessage(`{"foo":"bar"}`)

	eventStream, _ := json.Marshal(kafkaLogMessage)
	message := sarama.ConsumerMessage{Value: eventStream}
	publisher := new(mock.IKafkaProcessedLogPublisher)
	publisher.On("Publish","123", "topic-name", json.RawMessage(`{"foo":"bar"}`)).Once()
	kafkaLogProducerHandlerService := service.NewKafkaLogProducerHandlerService(publisher)

	//when
	kafkaLogProducerHandlerService.Handle(&message)

	//then
	publisher.AssertCalled(t, "Publish", "123", "topic-name", json.RawMessage(`{"foo":"bar"}`))
}