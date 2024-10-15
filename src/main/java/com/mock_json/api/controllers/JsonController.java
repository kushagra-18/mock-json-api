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
import com.mock_json.api.models.Json;
import com.mock_json.api.models.Project;
import com.mock_json.api.models.Url;
import com.mock_json.api.requests.JsonUrlRequest;
import com.mock_json.api.services.JsonService;
import com.mock_json.api.services.ProjectService;
import com.mock_json.api.services.RequestLogService;
import com.mock_json.api.services.UrlService;

import jakarta.servlet.http.HttpServletRequest;
import org.springframework.transaction.annotation.Transactional;

import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;

@RestController
@ResponseBody
public class JsonController {

    private final Logger logger = LoggerFactory.getLogger(HomeController.class);

    @Autowired
    private ProjectService projectService;

    @Autowired
    private JsonService jsonService;

    @Autowired
    private UrlService urlService;

    @Autowired
    private RequestLogService requestLogService;

    @PostMapping("/api/v1/json")
    @Transactional
    public ResponseEntity<?> saveJsonData(@Valid @RequestBody JsonUrlRequest jsonUrlRequest) {

        String urlString = jsonUrlRequest.getUrlData().getUrl();

        Optional<Url> existingUrl = urlService.findUrlDataByUrl(urlString);

        Project project = projectService.findProjectById(1L);

        Url urlData;

        if (existingUrl.isPresent()) {
            urlData = existingUrl.get();
        } else {
            urlData = urlService.saveData(jsonUrlRequest.getUrlData(), project);
        }

        Json savedJson = jsonService.saveJsonData(jsonUrlRequest.getJsonList().get(0), urlData);

        return ResponseEntity.ok(savedJson);
    }

    @GetMapping("/**")
    @HeaderIntercepted
    public ResponseEntity<?> getMockedJSON(HttpServletRequest request) {

        String url = jsonService.getUrl(request);

        String method = request.getMethod();

        String ip = request.getRemoteAddr();

        int status = 200;

        String teamSlug = HeaderContext.getTeamSlug();
        
        String projectSlug = HeaderContext.getProjectSlug();

        Optional<Url> urlData = urlService.findUrlDataByUrl(url);

        Json jsonData = jsonService.selectRandomJson(urlData.get().getJsonList());

        jsonService.simulateLatency(jsonData);

        ObjectMapper objectMapper = new ObjectMapper();

        Object jsonObject;

        String jsonDataString = jsonData.getJsonData();

        try {
            jsonObject = objectMapper.readValue(jsonDataString, Object.class);

        } catch (JsonProcessingException e) {
            return ResponseEntity.badRequest().body("Error parsing JSON data");
        }

        requestLogService.saveRequestLogAsync(url, method, ip, status, null);

        requestLogService.emitPusherEvent(method, url, null, status);

        return ResponseEntity.ok(jsonObject);
    }

}
