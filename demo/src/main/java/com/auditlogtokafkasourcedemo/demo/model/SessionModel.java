package com.auditlogtokafkasourcedemo.demo.model;

import java.time.LocalDateTime;

public class SessionModel {

    private String username;
    private String loginDate;

    public SessionModel(String username) {
        this.username = username;
        this.loginDate = LocalDateTime.now().toString();
    }

    public String getUsername() {
        return username;
    }

    public void setUsername(String username) {
        this.username = username;
    }

    public String getLoginDate() {
        return loginDate;
    }
}
