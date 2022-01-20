package com.auditlogtokafkasourcedemo.demo.model;

public class KafkaLog {
    private String topicName;
    private String key;
    private Object kafkaMessage;
    private String kafkaEnabled = "true";

    public KafkaLog() { }

    public KafkaLog(String topicName, String key, Object kafkaMessage) {
        this.topicName = topicName;
        this.key = key;
        this.kafkaMessage = kafkaMessage;
    }

    public String getTopicName() {
        return topicName;
    }

    public String getKey() {
        return key;
    }

    public Object getKafkaMessage() {
        return kafkaMessage;
    }

    public String getKafkaEnabled() {
        return kafkaEnabled;
    }
}
