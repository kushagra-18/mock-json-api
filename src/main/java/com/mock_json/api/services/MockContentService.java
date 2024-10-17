package com.mock_json.api.services;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;

import com.mock_json.api.helpers.StringHelpers;
import com.mock_json.api.models.MockContent;
import com.mock_json.api.models.Url;
import com.mock_json.api.repositories.MockContentRepository;

import jakarta.servlet.http.HttpServletRequest;

import java.time.LocalDateTime;
import java.util.List;
import java.util.Optional;
import java.util.Random;
import java.util.concurrent.TimeUnit;

@Service
public class MockContentService {

    private final MockContentRepository mockContentRepository;

    private static final Logger logger = LoggerFactory.getLogger(MockContentService.class);

    public MockContentService(MockContentRepository mockContentRepository) {
        this.mockContentRepository = mockContentRepository;
    }

    /**
     * Saves the JSON data to the database.
     * 
     * @param json
     * @param project
     * @return
     */
    public MockContent saveMockContentData(MockContent json, Url url) {

        LocalDateTime currTime = LocalDateTime.now();

        json.setCreatedAt(currTime);
        json.setUpdatedAt(currTime);
        json.setUrlId(url);

        if (json.getLatency() != null) {
            json.setLatency(json.getLatency());
        }

        return mockContentRepository.save(json);

    }

    /**
     * @description: This method simulates the latency of the JSON data.
     * @param json
     */
    public void simulateLatency(MockContent json) {

        Optional.ofNullable(json)
                .map(MockContent::getLatency)
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
     * @param mockContentList
     * @return
     */
    public MockContent selectRandomJson(List<MockContent> mockContentList) {
        int totalWeight = 0;
        
        for (MockContent json : mockContentList) {
            totalWeight += json.getRandomness();
        }

        Random random = new Random();
        
        int randomNumber = random.nextInt(totalWeight);

        for (MockContent json : mockContentList) {
            randomNumber -= json.getRandomness();
            if (randomNumber < 0) {
                return json;
            }
        }

        return null;
    }
}
