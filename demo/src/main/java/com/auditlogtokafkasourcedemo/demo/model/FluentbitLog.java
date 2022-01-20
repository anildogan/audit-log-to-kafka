package com.auditlogtokafkasourcedemo.demo.model;

public class FluentbitLog {
    private KafkaLog kafkaLog;

    public FluentbitLog(KafkaLog kafkaLog) {
        this.kafkaLog = kafkaLog;
    }

    public KafkaLog getKafkaLog() {
        return kafkaLog;
    }
}

