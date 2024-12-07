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
import org.springframework.data.rest.webmvc.ResourceNotFoundException;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PatchMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.client.HttpClientErrorException;
import org.springframework.web.client.HttpServerErrorException;
import org.springframework.web.client.RestTemplate;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.mock_json.mock_api.constants.ResponseMessages;
import com.mock_json.mock_api.dtos.MockContentUrlDto;
import com.mock_json.mock_api.dtos.UpdateMockContentUrlDto;
import com.mock_json.mock_api.exceptions.BadRequestException;
import com.mock_json.mock_api.exceptions.responses.RateLimitException;
import com.mock_json.mock_api.models.ForwardProxy;
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

import org.springframework.http.HttpHeaders;
import com.mock_json.mock_api.helpers.StringHelpers;

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

    @PatchMapping("{projectSlug}/{urlId}")
    @Transactional
    public ResponseEntity<Map<String, Object>> updateMockContentData(
            @PathVariable Long urlId,
            @PathVariable String projectSlug,
            @Valid @RequestBody UpdateMockContentUrlDto mockContentUrlDto) {

        Map<String, Object> response = new HashMap<>();

        Url urlData = urlService.findById(urlId)
                .orElseThrow(() -> new ResourceNotFoundException("URL with ID " + urlId + " not found"));

        if (!urlData.getProject().getSlug().equals(projectSlug)) {
            throw new BadRequestException("Project slug does not match the URL's associated project");
        }

        List<MockContent> updatedMockContentList = mockContentService.updateMockContentData(
                mockContentUrlDto.getMockContentList(), urlData);

        String mockedUrl = urlData.getProject().getSlug() + ".free." + baseUrl + "/" + urlData.getUrl();

        response.put("url", mockedUrl);
        response.put("data", updatedMockContentList);
        response.put("status_code", HttpStatus.OK.value());

        return ResponseEntity.ok(response);
    }

    /**
     * Get mocked JSON data: This is the main endpoint that returns the mocked JSON
     * data
     * based on team and project slugs, it also checks for rate limits and logs the
     * request.
     * if proxy is enabled for the project, it will forward the request to the proxy
     * server.
     * 
     * @param url
     * @param ip
     * @param teamSlug
     * @param projectSlug
     * @param request
     * @return
     */
    @GetMapping("/{teamSlug}/{projectSlug}")
    public ResponseEntity<?> getMockedJSON(
            @RequestParam(required = true) String url,
            @RequestParam(required = true) String ip,
            @PathVariable String teamSlug,
            @PathVariable String projectSlug,
            HttpServletRequest request) {

        // Decode URL and IP
        String decodedUrl = StringHelpers.decodeBase64(url);
        String decodedIpString = StringHelpers.decodeBase64(ip);

        // Rate limiting check
        if (this.globalRateLimit(teamSlug, projectSlug)) {
            return ResponseEntity.status(HttpStatus.TOO_MANY_REQUESTS)
                    .body(ResponseMessages.GLOBAL_RATE_LIMIT_EXCEEDED);
        }

        String method = request.getMethod();
        Optional<Url> urlDataOpt = urlService.findUrlDataByUrlAndTeamAndProject(teamSlug, projectSlug, decodedUrl);

        // If URL data is not found, handle directly
        if (urlDataOpt.isEmpty()) {
            return handleUrlNotFound(teamSlug, projectSlug, decodedUrl, decodedIpString, method, request);
        }

        // Handle the case when URL data is found
        return handleUrlDataFound(urlDataOpt.get(), decodedIpString, decodedUrl, method, request);
    }

    private ResponseEntity<?> handleUrlNotFound(String teamSlug, String projectSlug, String decodedUrl,
            String decodedIpString, String method, HttpServletRequest request) {

        // Query the project and check for forward proxy
        Project project = projectService.getDataBySlugAndTeamSlug(teamSlug, projectSlug);

        if (project.getIsForwardProxyActive()) {
            return handleForwardProxyRequest(decodedIpString, method, decodedUrl, project);
        }

        // Log the request and return default response
        Long projectId = project.getId();

        String channelId = project.getChannelId();
        requestLogService.saveRequestLogAsync(decodedIpString, null, method, decodedUrl, HttpStatus.OK.value(),
                projectId, false);

        requestLogService.emitPusherEvent(decodedIpString, null, method, decodedUrl, HttpStatus.OK.value(), channelId,
                false);

        return ResponseEntity.status(HttpStatus.OK).body(ResponseMessages.NO_CONTENT_URL);
    }

    private ResponseEntity<?> handleForwardProxyRequest(String decodedIpString, String method, String decodedUrl,
            Project project) {
        ForwardProxy forwardProxy = project.getForwardProxy();

        ResponseEntity<?> forwardProxyResponse = forwardRequestToProxyServer(forwardProxy, decodedUrl);

        // Log the forward proxy request
        Long projectId = project.getId();

        String channelId = project.getChannelId();

        requestLogService.saveRequestLogAsync(decodedIpString, null, method, decodedUrl,
                forwardProxyResponse.getStatusCodeValue(), projectId, true);

        requestLogService.emitPusherEvent(decodedIpString, null, method, decodedUrl,
                forwardProxyResponse.getStatusCodeValue(), channelId, true);

        return forwardProxyResponse;
    }

    private ResponseEntity<?> handleUrlDataFound(Url urlData, String decodedIpString, String decodedUrl, String method,
            HttpServletRequest request) {

        // Check if rate limit is exceeded
        String channelId = urlData.getProject().getChannelId();
        Long projectId = urlData.getProject().getId();
        Integer allowedRequests = urlData.getRequests();
        Long timeWindow = urlData.getTime();

        if (urlService.isRateLimited(decodedIpString, decodedUrl, allowedRequests, timeWindow)) {
            throw new RateLimitException(ResponseMessages.RATE_LIMIT_EXCEEDED);
        }

        // Simulate mock response
        MockContent mockContentData = mockContentService.selectRandomJson(urlData.getMockContentList());
        mockContentService.simulateLatency(mockContentData);

        // Parse the mock content into JSON
        Object jsonObject = parseJson(mockContentData.getData());

        // Log the request
        requestLogService.saveRequestLogAsync(decodedIpString, urlData, method, decodedUrl, HttpStatus.OK.value(),
                projectId, false);
        requestLogService.emitPusherEvent(decodedIpString, urlData, method, decodedUrl, HttpStatus.OK.value(),
                channelId,
                false);

        return ResponseEntity.ok(createResponse(jsonObject, urlData.getStatus().getCode()));
    }

    private Object parseJson(String jsonContent) {
        try {
            ObjectMapper objectMapper = new ObjectMapper();
            return objectMapper.readValue(jsonContent, Object.class);
        } catch (JsonProcessingException e) {
            return ResponseEntity.badRequest().body(ResponseMessages.JSON_PARSE_ERROR);
        }
    }

    /**
     * Forward the request to the proxy server
     * 
     * @param forwardProxy
     * @param decodedUrl
     * @return
     */
    private ResponseEntity<?> forwardRequestToProxyServer(ForwardProxy forwardProxy, String decodedUrl) {
        String domain = forwardProxy.getDomain();

        if (!domain.endsWith("/")) {
            domain += "/";
        }

        // Construct the full URL
        String url = domain + decodedUrl;

        try {
            RestTemplate restTemplate = new RestTemplate();

            // Make the HTTP GET request
            var response = restTemplate.getForEntity(url, String.class);

            HttpHeaders headers = response.getHeaders();

            @SuppressWarnings("null")
            String contentType = headers.getContentType() != null ? headers.getContentType().toString() : "";

            if (contentType.contains("application/json")) {
                return ResponseEntity.status(response.getStatusCode())
                        .headers(headers)
                        .body(response.getBody());
            } else if (contentType.contains("text/html")) {
                return ResponseEntity.status(response.getStatusCode())
                        .headers(headers)
                        .body(response.getBody());
            } else {
                return ResponseEntity.status(response.getStatusCode())
                        .headers(headers)
                        .body(response.getBody());
            }
        } catch (HttpClientErrorException | HttpServerErrorException e) {
            return ResponseEntity.status(e.getStatusCode())
                    .headers(e.getResponseHeaders())
                    .body(e.getResponseBodyAsString());
        } catch (Exception e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
                    .body("Request failed: ");
        }
    }

    private Map<String, Object> createResponse(Object jsonObject, Integer statusCode) {
        Map<String, Object> response = new HashMap<>();
        response.put("json_data", jsonObject);
        response.put("status_code", statusCode);
        return response;
    }

    private boolean globalRateLimit(String teamSlug, String projectSlug) {

        String redisKey = redisService.createRedisKey("rate_limit_global", teamSlug, projectSlug);

        return redisService.rateLimit(redisKey, this.getMaxAllowedRequests(), this.getTimeWindow());

    }
}
