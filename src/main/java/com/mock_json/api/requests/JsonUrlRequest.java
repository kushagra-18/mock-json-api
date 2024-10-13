package com.mock_json.api.requests;

import java.util.List;

import javax.validation.Valid;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.mock_json.api.models.Json;
import com.mock_json.api.models.Url;

public class JsonUrlRequest {
    @Valid
    @JsonProperty("url_data")
    private Url urlData;

    @Valid
    @JsonProperty("json_list")
    private List<Json> jsonList;

    public Url getUrlData() {
        return urlData;
    }

    public void setUrlData(Url urlData) {
        this.urlData = urlData;
    }

    public List<Json> getJsonList() {
        return jsonList;
    }

    public void setJsonList(List<Json> jsonList) {
        this.jsonList = jsonList;
    }   
}
