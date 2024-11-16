package com.mock_json.mock_api.exceptions.responses;

import java.time.LocalDateTime;

public class ErrorResponse {

    private LocalDateTime timestamp;
    private int status_code;
    private String error;
    private String message;
    private String path;
    
    public ErrorResponse(
        int status_code, String error, String message, String path, String appEnv) {
        this.timestamp = LocalDateTime.now();
        this.status_code = status_code;
        this.error = error;
        if (status_code == 500 && "production".equalsIgnoreCase(appEnv)) {
            this.message = "Internal Server Error"; 
        } else {
            this.message = message; 
        }
        this.path = path;
    }

    public LocalDateTime getTimestamp() {
        return timestamp;
    }

    public int getStatus_code() {
        return status_code;
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
