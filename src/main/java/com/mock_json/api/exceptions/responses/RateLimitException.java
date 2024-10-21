package com.mock_json.api.exceptions.responses;

public class RateLimitException extends RuntimeException {
   
    public RateLimitException(String message) {
        super(message);
    }

    public RateLimitException(String message, Throwable cause) {
        super(message, cause);
    }
}
