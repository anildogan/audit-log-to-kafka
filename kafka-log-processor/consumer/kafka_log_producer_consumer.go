package consumer

import (
	"github.com/Shopify/sarama"
	kafka "github.com/anildogan/audit-log-to-kafka/kafka-log-processor/kafka/consumer"
	"github.com/anildogan/audit-log-to-kafka/kafka-log-processor/service"
	"go.uber.org/zap"
)

type EventHandler struct {
	handlerService service.IKafkaLogProducerHandlerService
	producer       sarama.SyncProducer
	errorTopic     string
}

func NewEventHandler(handlerService service.IKafkaLogProducerHandlerService, producer sarama.SyncProducer, errorTopic string) kafka.EventHandler {
	return &EventHandler{
		handlerService: handlerService,
		errorTopic:     errorTopic,
		producer:       producer,
	}
}

func (e *EventHandler) Setup(session sarama.ConsumerGroupSession) error {
	zap.S().Info("Kafka audit log producer consumer is started.")
	return nil
}

func (e *EventHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (e *EventHandler) ConsumeClaim(session sarama.ConsumerGroupSession, c sarama.ConsumerGroupClaim) error {
	for message := range c.Messages() {
		if c.Topic() != e.errorTopic {
			err := e.process(message)
			if err != nil {
				zap.S().Errorf("Error executing Kafka audit log consumer, message: %s , err: %+v", string(message.Value), err)
			}
			session.MarkMessage(message, "")
		}
	}
	return nil
}

func (e *EventHandler) process(message *sarama.ConsumerMessage) error {
	zap.S().Infof("Message will be executed. message: %s", string(message.Value))
	err := e.handlerService.Handle(message)
	if err == nil {
		zap.S().Infof("Message is executed successfully, message: %s", string(message.Value))
	} else {
		zap.S().Infof("Message is not executed successfully. so is routing to error topic: %s, message: %s", e.errorTopic, string(message.Value))
		zap.S().Errorf("Error is : %s", err)
		err = e.sendToErrorTopic(message, e.errorTopic, err.Error())
		if err != nil {
			zap.S().Infof("Message is not published to error topic: %s", e.errorTopic)
		}
	}
	return err
}

func (e *EventHandler) sendToErrorTopic(message *sarama.ConsumerMessage, errorTopic string, errorMessage string) error {
	_, _, err := e.producer.SendMessage(&sarama.ProducerMessage{
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte(kafka.ErrorKey),
				Value: []byte(errorMessage),
			},
		},
		Topic: errorTopic,
		Value: sarama.StringEncoder(message.Value),
	})
	return err
}
