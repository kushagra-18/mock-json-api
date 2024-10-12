package com.mock_json.api.controllers;

import java.time.LocalDateTime;
import java.util.HashMap;
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
import com.mock_json.api.models.Json;
import com.mock_json.api.models.Project;
import com.mock_json.api.models.Team;
import com.mock_json.api.repositories.JsonRepository;
import com.mock_json.api.repositories.ProjectRepository;
import com.mock_json.api.services.JsonService;
import com.mock_json.api.services.ProjectService;

import jakarta.servlet.http.HttpServletRequest;

import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;

@RestController
@ResponseBody
public class JsonController {

    private static final Logger logger = LoggerFactory.getLogger(HomeController.class);

    @Autowired
    private ProjectService projectService;

    @Autowired
    private JsonService jsonService;

    @PostMapping("/api/v1/json")
    public ResponseEntity<?> saveJsonData(@Valid @RequestBody Json json) {

        Project project = projectService.findProjectById(1L);

        Json savedJson = jsonService.saveJsonData(json, project);

        return ResponseEntity.ok(savedJson);
    }

    @GetMapping("/**")
    public ResponseEntity<?> home(HttpServletRequest request) {

        String url = jsonService.getUrl(request);

        Json json = jsonService.findJsonByUrl(url);

        ObjectMapper objectMapper = new ObjectMapper();

        Object jsonObject;

        String jsonData = json.getJsonData();

        try {
            jsonObject = objectMapper.readValue(jsonData, Object.class);
            return ResponseEntity.ok(jsonObject);
        } catch (JsonProcessingException e) {
            return ResponseEntity.badRequest().body("Error parsing JSON data: " + e.getMessage());
        }
    }

}
