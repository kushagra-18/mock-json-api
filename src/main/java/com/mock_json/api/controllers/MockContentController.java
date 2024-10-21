package com.mock_json.api.controllers;

import java.util.HashMap;
import java.util.Map;
import java.util.Optional;

import javax.validation.Valid;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.mock_json.api.annotations.HeaderIntercepted;
import com.mock_json.api.contexts.HeaderContext;
import com.mock_json.api.dtos.MockContentUrlDto;
import com.mock_json.api.exceptions.NotFoundException;
import com.mock_json.api.exceptions.responses.RateLimitException;
import com.mock_json.api.models.MockContent;
import com.mock_json.api.models.Url;
import com.mock_json.api.services.MockContentService;
import com.mock_json.api.services.ProjectService;
import com.mock_json.api.services.RequestLogService;
import com.mock_json.api.services.UrlService;

import jakarta.servlet.http.HttpServletRequest;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.http.HttpStatus;


@RestController
@ResponseBody
@Validated
public class MockContentController {

    private final Logger logger = LoggerFactory.getLogger(HomeController.class);

    @Autowired
    private ProjectService projectService;

    @Autowired
    private MockContentService mockContentService;

    @Autowired
    private UrlService urlService;

    @Autowired
    private RequestLogService requestLogService;

    @Value("${BASE_URL}")
    private String baseUrl;

    @PostMapping("/api/v1/mock")
    @Transactional
    public ResponseEntity<Map<String, Object>> saveMockContentData(
            @Valid @RequestBody MockContentUrlDto mockContentUrlDto) {

        Map<String, Object> response = new HashMap<>();

        String urlString = mockContentUrlDto.getUrlData().getUrl();

        Url urlData = urlService.findUrlDataByUrl(urlString)
                .orElseGet(
                        () -> urlService.saveData(mockContentUrlDto.getUrlData(), projectService.findProjectById(1L)));

        MockContent mockContent = mockContentUrlDto.getMockContentList().get(0);
        MockContent savedMockedData = mockContentService.saveMockContentData(mockContent, urlData);

        String mockedUrl = urlString + ".free." + baseUrl;

        response.put("url", mockedUrl);
        response.put("data", savedMockedData);
        response.put("status_code", HttpStatus.CREATED.value());

        return ResponseEntity.status(HttpStatus.CREATED).body(response);
    }

    
    @GetMapping("/**")
    @HeaderIntercepted
    public ResponseEntity<?> getMockedJSON(HttpServletRequest request) {

        String url = mockContentService.getUrl(request);

        String method = request.getMethod();

        String ip = request.getRemoteAddr();

        String teamSlug = HeaderContext.getTeamSlug();

        String projectSlug = HeaderContext.getProjectSlug();

        Optional<Url> urlData = urlService.findUrlDataByUrlAndTeamAndProject(teamSlug, projectSlug, url);

        if (!urlData.isPresent()) {
            throw new NotFoundException("Url not found");
        }

        Integer allowedRequests = urlData.get().getRequests();
        Long timeWindow = urlData.get().getTime();

        if (urlService.isRateLimited(ip, url, allowedRequests, timeWindow)) {
            throw new RateLimitException("Rate limit exceeded, Please try again later.");
        }

        MockContent mockContentData = mockContentService.selectRandomJson(urlData.get().getMockContentList());

        mockContentService.simulateLatency(mockContentData);

        ObjectMapper objectMapper = new ObjectMapper();

        Object jsonObject;

        String mockContentDataString = mockContentData.getData();

        try {
            jsonObject = objectMapper.readValue(mockContentDataString, Object.class);

        } catch (JsonProcessingException e) {
            return ResponseEntity.badRequest().body("Error parsing JSON data");
        }

        requestLogService.saveRequestLogAsync(url, method, ip, HttpStatus.OK.value(), null);

        requestLogService.emitPusherEvent(method, url, null, HttpStatus.OK.value());

        return ResponseEntity.ok(jsonObject);
    }
}
