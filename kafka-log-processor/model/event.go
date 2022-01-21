package model

import "encoding/json"

type KafkaLogMessage struct {
	ContainerName string `json:"container_name"`
	Log struct {
		KafkaLog struct {
			TopicName string `json:"topicName"`
			Key string `json:"key"`
			KafkaMessage json.RawMessage `json:"kafkaMessage"`
		} `json:"kafkaLog"`
	} `json:"log"`
}
