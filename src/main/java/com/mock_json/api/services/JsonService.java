package com.mock_json.api.services;

import org.springframework.stereotype.Service;

import com.mock_json.api.helpers.StringHelpers;
import com.mock_json.api.models.Json;
import com.mock_json.api.models.Project;
import com.mock_json.api.repositories.JsonRepository;

import jakarta.servlet.http.HttpServletRequest;

import java.net.URI;
import java.net.URISyntaxException;
import java.time.LocalDateTime;

@Service
public class JsonService {

    private final JsonRepository jsonRepository;

    public JsonService(JsonRepository jsonRepository) {
        this.jsonRepository = jsonRepository;
    }

    public boolean checkURLExists(String url) {
        // TODO: Implement this method to check if the URL exists
        return true;
    }

    /**
     * Saves the JSON data to the database.
     * 
     * @param json
     * @param project
     * @return
     */
    public Json saveJsonData(Json json, Project project) {

        LocalDateTime currTime = LocalDateTime.now();
        json.setCreatedAt(currTime);
        json.setUpdatedAt(currTime);
        json.setProject(project);

        return jsonRepository.save(json);

    }

    public Json findJsonById(Long jsonId) {
        return jsonRepository.findById(jsonId)
                .orElseThrow(() -> new RuntimeException("Json with ID " + jsonId + " not found"));
    }

    public Json findJsonByUrl(String url) {
        return jsonRepository.findByUrl(url)
                .orElseThrow(() -> new RuntimeException("Json with URL " + url + " not found"));
    }

    /**
     * Returns the full URL with query string.
     * which was recieved in the request
     * 
     * @param request
     */
    public String getUrl(HttpServletRequest request) {
       
        StringBuilder fullPathWithQuery = new StringBuilder(request.getRequestURI());
        
        String queryString = request.getQueryString();

        if (queryString != null) {
            fullPathWithQuery.append("?").append(queryString);
        }

        return StringHelpers.ltrim(fullPathWithQuery.toString(), '/');
    }
}
