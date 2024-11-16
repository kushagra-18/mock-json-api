package com.mock_json.mock_api.controllers;

import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Optional;

import javax.validation.Valid;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RestController;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.mock_json.mock_api.annotations.HeaderIntercepted;
import com.mock_json.mock_api.constants.ResponseMessages;
import com.mock_json.mock_api.contexts.HeaderContext;
import com.mock_json.mock_api.dtos.MockContentUrlDto;
import com.mock_json.mock_api.exceptions.NotFoundException;
import com.mock_json.mock_api.exceptions.responses.RateLimitException;
import com.mock_json.mock_api.models.MockContent;
import com.mock_json.mock_api.models.Url;
import com.mock_json.mock_api.services.MockContentService;
import com.mock_json.mock_api.services.ProjectService;
import com.mock_json.mock_api.services.RequestLogService;
import com.mock_json.mock_api.services.UrlService;

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

    @PostMapping("/api/v1/mock/{projectSlug}")
    @Transactional
    public ResponseEntity<Map<String, Object>> saveMockContentData(
            @Valid @RequestBody MockContentUrlDto mockContentUrlDto, @PathVariable String projectSlug) {

        Map<String, Object> response = new HashMap<>();

        String urlString = mockContentUrlDto.getUrlData().getUrl();

        Url urlData = urlService.findUrlDataByUrlAndProjectSlug(projectSlug, urlString)
                .orElseGet(
                        () -> urlService.saveData(mockContentUrlDto.getUrlData(),
                                projectService.findProjectBySlug(projectSlug)));

        List<MockContent> mockContentList = mockContentUrlDto.getMockContentList();
        List<MockContent> savedMockedDataList = mockContentService.saveMockContentData(mockContentList, urlData);

        String mockedUrl = projectSlug + ".free." + baseUrl + "/" + urlString;

        response.put("url", mockedUrl);
        response.put("data", savedMockedDataList);
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

            requestLogService.saveRequestLogAsync(ip, null, method, url, HttpStatus.OK.value());

            requestLogService.emitPusherEvent(ip, null, method, url, HttpStatus.OK.value());

            return ResponseEntity.status(HttpStatus.OK)
                    .body(ResponseMessages.NO_CONTENT_URL);
        }

        Integer allowedRequests = urlData.get().getRequests();

        Long timeWindow = urlData.get().getTime();

        if (urlService.isRateLimited(ip, url, allowedRequests, timeWindow)) {
            throw new RateLimitException(ResponseMessages.RATE_LIMIT_EXCEEDED);
        }

        MockContent mockContentData = mockContentService.selectRandomJson(urlData.get().getMockContentList());

        mockContentService.simulateLatency(mockContentData);

        ObjectMapper objectMapper = new ObjectMapper();

        Object jsonObject;

        String mockContentDataString = mockContentData.getData();

        try {
            jsonObject = objectMapper.readValue(mockContentDataString, Object.class);

        } catch (JsonProcessingException e) {
            return ResponseEntity.badRequest().body(ResponseMessages.JSON_PARSE_ERROR);
        }

        requestLogService.saveRequestLogAsync(ip, urlData.get(), method, url, HttpStatus.OK.value());

        requestLogService.emitPusherEvent(ip, urlData.get(), method, url, HttpStatus.OK.value());

        return ResponseEntity.ok(jsonObject);
    }
}
