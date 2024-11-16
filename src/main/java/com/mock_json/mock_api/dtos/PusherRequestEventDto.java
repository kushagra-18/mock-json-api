package com.mock_json.mock_api.dtos;

import com.fasterxml.jackson.annotation.JsonProperty;

import lombok.AllArgsConstructor;
import lombok.Data;

@Data
@AllArgsConstructor
public class PusherRequestEventDto {

    private String method;

    @JsonProperty("url_string")
    private String urlString;

    @JsonProperty("url_id")
    private Long urlId;

    private String ip;
    
    private int status;
}
