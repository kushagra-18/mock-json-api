package com.mock_json.mock_api.exceptions.responses;

import java.time.LocalDateTime;

public class ErrorResponse {

    private LocalDateTime timestamp;
    private int status;
    private String error;
    private String message;
    private String path;
    
    public ErrorResponse(
            int status, String error, String message, String path, String appEnv) {
        this.timestamp = LocalDateTime.now();
        this.status = status;
        this.error = error;
        if (status == 500 && "production".equalsIgnoreCase(appEnv)) {
            this.message = "Internal Server Error"; 
        } else {
            this.message = message; 
        }
        this.path = path;
    }

    public LocalDateTime getTimestamp() {
        return timestamp;
    }

    public int getStatus() {
        return status;
    }

    public String getError() {
        return error;
    }

    public String getMessage() {

        return message;
    }

    public String getPath() {
        return path;
    }
}
