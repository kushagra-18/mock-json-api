package com.mock_json.testing_api.services;

import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;

import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;

@Service
public class WebTestingService {

    private static final int THREAD_COUNT = 10; // Number of concurrent threads
    private static final int REQUEST_COUNT = 100; // Total number of requests to send
    private static final String API_URL = "https://stage.vipfortunes.com/";

    private final RestTemplate restTemplate = new RestTemplate();

    public void hitURL() {
        ExecutorService executorService = Executors.newFixedThreadPool(THREAD_COUNT);
        List<Long> responseTimes = new ArrayList<>();
        List<Integer> statusCodes = new ArrayList<>();  
        List<String> failedRequests = new ArrayList<>();

       for (int i = 0; i < REQUEST_COUNT; i++) {
            executorService.submit(() -> {
                long startTime = System.currentTimeMillis();
                try {
                    ResponseEntity<String> response = restTemplate.getForEntity(API_URL, String.class);
                    long responseTime = System.currentTimeMillis() - startTime;

                    System.out.println("Request " + " took " + responseTime + " ms");

                    synchronized (responseTimes) {
                        responseTimes.add(responseTime);
                    }
                    synchronized (statusCodes) {
                        statusCodes.add(response.getStatusCodeValue());
                    }

                } catch (Exception e) {
                    System.err.println("Request failed: " + e.getMessage());
                }
            });
        }

        // Shut down the executor and await termination
        executorService.shutdown();
        try {
            if (!executorService.awaitTermination(1, TimeUnit.MINUTES)) {
                executorService.shutdownNow();
            }
        } catch (InterruptedException e) {
            executorService.shutdownNow();
        }

        // Calculate average response time
        double averageResponseTime = responseTimes.stream().mapToLong(Long::longValue).average().orElse(0);
        System.out.println("Average response time: " + averageResponseTime + " ms");
        System.out.println("Total requests: " + REQUEST_COUNT);
        System.out.println("Successful responses: " + responseTimes.size());

        System.out.println("Status Codes Returned:");
        statusCodes.stream().distinct().forEach(code ->
            System.out.println("Status Code " + code + ": " + 
                statusCodes.stream().filter(status -> status == code).count() + " times")
        );
    }
}
