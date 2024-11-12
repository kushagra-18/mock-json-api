package com.mock_json.testing_api.entities;

import lombok.Data;
import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.mapping.Document;
import org.springframework.data.mongodb.core.mapping.Field;

import java.time.LocalDateTime;
import java.util.Map;

@Data
@Document(collection = "test_results")
public class TestResultEntity {

    @Id
    @Field("id")
    private String id;

    @Field("test_id")
    private String testId;

    @Field("request_url")
    private String requestUrl;

    @Field("method")
    private String method;

    @Field("request_headers")
    private Map<String, String> requestHeaders;

    @Field("request_body")
    private String requestBody;

    @Field("response_status")
    private int responseStatus;

    @Field("response_headers")
    private Map<String, String> responseHeaders;

    @Field("response_body")
    private String responseBody;

    @Field("response_time_ms")
    private long responseTimeMs;

    @Field("success")
    private boolean success;

    @Field("timestamp")
    private LocalDateTime timestamp;
}
