package com.auditlogtokafkasourcedemo.demo.util.fluentbit;

import com.auditlogtokafkasourcedemo.demo.model.FluentbitLog;
import com.auditlogtokafkasourcedemo.demo.model.KafkaLog;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public final class KafkaLogger {

    private static final Logger logger = LoggerFactory.getLogger(KafkaLogger.class);

    public static void log(Object message, String topic, String key) {
        FluentbitLog fluentbitLog = new FluentbitLog(new KafkaLog(topic, key, message));

        try {
            logger.info(new ObjectMapper().writeValueAsString(fluentbitLog));
        } catch (JsonProcessingException ignored) { }
    }
}