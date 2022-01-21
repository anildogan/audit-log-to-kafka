package mock

import (
	"encoding/json"
	"github.com/stretchr/testify/mock"
)

type IKafkaProcessedLogPublisher struct {
	mock.Mock
}

func (_p *IKafkaProcessedLogPublisher) Publish(key string, topicName string, rawMessage json.RawMessage) error {
	_p.Called(key, topicName, rawMessage)
	return nil
}
