package com.mock_json.mock_api.requests;

import java.util.List;

import javax.validation.Valid;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.mock_json.mock_api.models.MockContent;
import com.mock_json.mock_api.models.Url;

public class JsonUrlRequest {
    @Valid
    @JsonProperty("url_data")
    private Url urlData;

    @Valid
    @JsonProperty("mock_content_list")
    private List<MockContent> mockContentList;

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
