package com.auditlogtokafkasourcedemo.demo.controller;

import com.auditlogtokafkasourcedemo.demo.model.SessionModel;
import com.auditlogtokafkasourcedemo.demo.util.fluentbit.KafkaLogger;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;


@Controller
public class DemoAuditLogController {

    @PostMapping("/demo/login")
    public void runDemo(@RequestBody String username) {
        KafkaLogger.log(new SessionModel(username), "user.login", username);
    }
}
