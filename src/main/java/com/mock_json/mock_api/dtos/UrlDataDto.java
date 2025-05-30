package com.mock_json.mock_api.dtos;

import com.fasterxml.jackson.annotation.JsonProperty;

import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Size;
import lombok.Data;

@Data
public class UrlDataDto {

    @NotNull(message = "URL cannot be null")
    private String url;

    @JsonProperty("description")
    @Size(min = 1, message = "Description must be at least 1 character")
    @Size(max = 255, message = "Description must be at most 255 characters")
    private String description;

    @JsonProperty("name")
    private String name;

    @JsonProperty("requests")
    private Integer requests;

    @JsonProperty("time")
    private Long time;
    
}
