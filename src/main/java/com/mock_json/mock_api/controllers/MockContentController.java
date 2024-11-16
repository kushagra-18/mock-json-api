package com.mock_json.mock_api.controllers;

import java.util.Base64;
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
import com.mock_json.mock_api.constants.ResponseMessages;
import com.mock_json.mock_api.dtos.MockContentUrlDto;
import com.mock_json.mock_api.exceptions.responses.RateLimitException;
import com.mock_json.mock_api.models.MockContent;
import com.mock_json.mock_api.models.Project;
import com.mock_json.mock_api.models.Url;
import com.mock_json.mock_api.services.MockContentService;
import com.mock_json.mock_api.services.ProjectService;
import com.mock_json.mock_api.services.RedisService;
import com.mock_json.mock_api.services.RequestLogService;
import com.mock_json.mock_api.services.UrlService;

import jakarta.servlet.http.HttpServletRequest;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.http.HttpStatus;

@RestController
@ResponseBody
@RequestMapping("mock")
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

    @Autowired
    private RedisService redisService;

    @Value("${BASE_URL}")
    private String baseUrl;

    @Value("${GLOBAL_MAX_ALLOWED_REQUESTS}")
    private Integer maxAllowedRequests;

    public Integer getMaxAllowedRequests() {
        return maxAllowedRequests;
    }

    @Value("${GLOBAL_TIME_WINDOW}")
    private Long timeWindow;

    public Long getTimeWindow() {
        return timeWindow;
    }

    @PostMapping("{projectSlug}")
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

    @GetMapping("{teamSlug}/{projectSlug}")
    public ResponseEntity<?> getMockedJSON(
            @RequestParam(required = true) String url,
            @PathVariable String teamSlug,
            @PathVariable String projectSlug,
            HttpServletRequest request) {

        byte[] decodedBytes = Base64.getDecoder().decode(url);
        String decodedUrl = new String(decodedBytes);

        if (this.globalRateLimit(teamSlug, projectSlug)) {
            return ResponseEntity.status(HttpStatus.TOO_MANY_REQUESTS)
                    .body(ResponseMessages.GLOBAL_RATE_LIMIT_EXCEEDED);
        }

        String method = request.getMethod();

        String ip = request.getRemoteAddr();

        Optional<Url> urlDataOpt = urlService.findUrlDataByUrlAndTeamAndProject(teamSlug, projectSlug, decodedUrl);

        // If URL data is not found, handle the response directly
        if (urlDataOpt.isEmpty()) {
            Project project = projectService.getDataBySlugAndTeamSlug(teamSlug, projectSlug);

            Long projectId = project.getId();
            
            String channelId = project.getChannelId();

            requestLogService.saveRequestLogAsync(ip, null, method, decodedUrl, HttpStatus.OK.value(), projectId);
            
            requestLogService.emitPusherEvent(ip, null, method, url, HttpStatus.OK.value(), channelId);

            return ResponseEntity.status(HttpStatus.OK).body(ResponseMessages.NO_CONTENT_URL);
        }

        Url urlData = urlDataOpt.get();
       
        String channelId = urlData.getProject().getChannelId();
        
        Long projectId = urlData.getProject().getId();
        
        Integer allowedRequests = urlData.getRequests();
        
        Long timeWindow = urlData.getTime();

        if (urlService.isRateLimited(ip, url, allowedRequests, timeWindow)) {
            throw new RateLimitException(ResponseMessages.RATE_LIMIT_EXCEEDED);
        }

        MockContent mockContentData = mockContentService.selectRandomJson(urlData.getMockContentList());
        
        mockContentService.simulateLatency(mockContentData);

        ObjectMapper objectMapper = new ObjectMapper();
        
        Object jsonObject;
        
        String mockContentDataString = mockContentData.getData();

        try {
            jsonObject = objectMapper.readValue(mockContentDataString, Object.class);
        } catch (JsonProcessingException e) {
            return ResponseEntity.badRequest().body(ResponseMessages.JSON_PARSE_ERROR);
        }

        requestLogService.saveRequestLogAsync(ip, urlData, method, decodedUrl, HttpStatus.OK.value(), projectId);
        
        requestLogService.emitPusherEvent(ip, urlData, method, decodedUrl, HttpStatus.OK.value(), channelId);

        return ResponseEntity.ok(jsonObject);
    }

    private boolean globalRateLimit(String teamSlug, String projectSlug) {

        String redisKey = redisService.createRedisKey("rate_limit_global", teamSlug, projectSlug);

        return redisService.rateLimit(redisKey, this.getMaxAllowedRequests(), this.getTimeWindow());

    }
}
