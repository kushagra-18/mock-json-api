package com.mock_json.api.services;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;

import com.mock_json.api.helpers.StringHelpers;
import com.mock_json.api.models.Json;
import com.mock_json.api.models.Url;
import com.mock_json.api.repositories.JsonRepository;

import jakarta.servlet.http.HttpServletRequest;

import java.time.LocalDateTime;
import java.util.List;
import java.util.Optional;
import java.util.Random;
import java.util.concurrent.TimeUnit;

@Service
public class JsonService {

    private final JsonRepository jsonRepository;

    private static final Logger logger = LoggerFactory.getLogger(JsonService.class);

    public JsonService(JsonRepository jsonRepository) {
        this.jsonRepository = jsonRepository;
    }

    /**
     * Saves the JSON data to the database.
     * 
     * @param json
     * @param project
     * @return
     */
    public Json saveJsonData(Json json, Url url) {

        LocalDateTime currTime = LocalDateTime.now();

        json.setCreatedAt(currTime);
        json.setUpdatedAt(currTime);
        json.setUrlId(url);

        if (json.getLatency() != null) {
            json.setLatency(json.getLatency());
        }

        return jsonRepository.save(json);

    }

    /**
     * @description: This method simulates the latency of the JSON data.
     * @param json
     */
    public void simulateLatency(Json json) {

        Optional.ofNullable(json)
                .map(Json::getLatency)
                .ifPresent(latency -> {
                    if (latency > 0) {
                        try {
                            TimeUnit.MILLISECONDS.sleep(latency);
                        } catch (InterruptedException e) {
                            Thread.currentThread().interrupt();
                            logger.error("Thread was interrupted during sleep: " + e.getMessage());
                        }
                    }
                });
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

    /**
     * Selects a random JSON object from the list of JSON objects.
     * based on the randomness of each JSON object.
     * @param jsonList
     * @return
     */
    public Json selectRandomJson(List<Json> jsonList) {
        int totalWeight = 0;
        
        for (Json json : jsonList) {
            totalWeight += json.getRandomness();
        }

        Random random = new Random();
        
        int randomNumber = random.nextInt(totalWeight);

        for (Json json : jsonList) {
            randomNumber -= json.getRandomness();
            if (randomNumber < 0) {
                return json;
            }
        }

        return null;
    }
}
