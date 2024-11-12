package com.mock_json.mock_api.dtos;

import java.util.List;


import com.fasterxml.jackson.annotation.JsonProperty;
import com.mock_json.mock_api.models.MockContent;
import com.mock_json.mock_api.models.Url;

import jakarta.validation.Valid;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Size;
import lombok.Data;

@Data
public class MockContentUrlDto {
    
    @Valid
    @NotNull(message = "URL data cannot be null")
    @JsonProperty("url_data")
    private Url urlData;

    @Valid
    @NotNull(message = "test data cannot be null")
    private String test;

    @Valid
    @NotNull(message = "Mock content list cannot be null")
    @Size(min = 1, message = "Mock content list must contain at least one item")
    @JsonProperty("mock_content_list")
    private List<@Valid MockContent> mockContentList;

    public Url getUrlData() {
        return urlData;
    }

    public void setUrlData(Url urlData) {
        this.urlData = urlData;
    }

    public List<MockContent> getMockContentList() {
        return mockContentList;
    }

    public void setMockContentList(List<MockContent> mockContentList) {
        this.mockContentList = mockContentList;
    }   
}
