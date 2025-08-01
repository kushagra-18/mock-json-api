package com.mock_json.mock_api.services;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;

import com.mock_json.mock_api.helpers.StringHelpers;
import com.mock_json.mock_api.models.MockContent;
import com.mock_json.mock_api.models.Url;
import com.mock_json.mock_api.repositories.MockContentRepository;

import jakarta.servlet.http.HttpServletRequest;

import java.time.LocalDateTime;
import java.util.ArrayList;
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
    public List<MockContent> saveMockContentData(List<MockContent> mockContentList, Url url) {
        LocalDateTime currTime = LocalDateTime.now();

        for (MockContent mockContent : mockContentList) {
            mockContent.setCreatedAt(currTime);
            mockContent.setUpdatedAt(currTime);
            mockContent.setUrlId(url);

            if (mockContent.getLatency() != null) {
                mockContent.setLatency(mockContent.getLatency());
            }
        }

        return mockContentRepository.saveAll(mockContentList);
    }

    public List<MockContent> updateMockContentData(List<MockContent> mockContentList, Url url) {
        LocalDateTime currTime = LocalDateTime.now();

        List<MockContent> updatedMockContents = new ArrayList<>();

        for (MockContent mockContent : mockContentList) {
            MockContent existingContent = mockContentRepository
                    .findById(mockContent.getId()) 
                    .orElse(null);

            if (existingContent != null) {
                existingContent.setData(mockContent.getData());
                existingContent.setDescription(mockContent.getDescription());
                existingContent.setName(mockContent.getName());
                existingContent.setLatency(mockContent.getLatency());
                existingContent.setRandomness(mockContent.getRandomness());
                existingContent.setUpdatedAt(currTime);
                updatedMockContents.add(existingContent);
            } else {
                mockContent.setCreatedAt(currTime);
                mockContent.setUpdatedAt(currTime);
                mockContent.setUrlId(url);
                updatedMockContents.add(mockContent);
            }
        }

        return mockContentRepository.saveAll(updatedMockContents);
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
     * 
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
