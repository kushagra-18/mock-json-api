package com.mock_json.api.controllers;

import java.util.Optional;

import javax.validation.Valid;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.mock_json.api.annotations.HeaderIntercepted;
import com.mock_json.api.contexts.HeaderContext;
import com.mock_json.api.dtos.MockContentUrlDto;
import com.mock_json.api.exceptions.NotFoundException;
import com.mock_json.api.models.MockContent;
import com.mock_json.api.models.Project;
import com.mock_json.api.models.Url;
import com.mock_json.api.requests.cockContentUrlDto;
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

    @PostMapping("/api/v1/mock")
    @Transactional
    public ResponseEntity<?> saveMockContentData(@Valid @RequestBody MockContentUrlDto cockContentUrlDto) {

        String urlString = cockContentUrlDto.getUrlData().getUrl();

        Optional<Url> existingUrl = urlService.findUrlDataByUrl(urlString);

        Project project = projectService.findProjectById(1L);

        Url urlData;

        if (existingUrl.isPresent()) {
            urlData = existingUrl.get();
        } else {
            urlData = urlService.saveData(cockContentUrlDto.getUrlData(), project);
        }

        MockContent savedMockedData = mockContentService.saveMockContentData(cockContentUrlDto.getMockContentList().get(0), urlData);

        return ResponseEntity.ok(savedMockedData);
    }

    @GetMapping("/**")
    @HeaderIntercepted
    public ResponseEntity<?> getMockedJSON(HttpServletRequest request) {

        String url = mockContentService.getUrl(request);

        String method = request.getMethod();

        String ip = request.getRemoteAddr();

        int status = 200;

        String teamSlug = HeaderContext.getTeamSlug();
        
        String projectSlug = HeaderContext.getProjectSlug();

        Optional<Url> urlData = urlService.findUrlDataByUrlAndTeamAndProject(teamSlug, projectSlug, url);

        if(!urlData.isPresent()) {
             throw new NotFoundException("Url not found");
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

        requestLogService.saveRequestLogAsync(url, method, ip, status, null);

        requestLogService.emitPusherEvent(method, url, null, status);

        return ResponseEntity.ok(jsonObject);
    }

}
