package com.mock_json.testing_api.controllers;
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
import com.mock_json.mock_api.annotations.HeaderIntercepted;
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
import com.mock_json.testing_api.services.WebTestingService;

import jakarta.servlet.http.HttpServletRequest;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.http.HttpStatus;


@RestController
@ResponseBody
@RequestMapping("api/v1/test")
@Validated
public class WebTestingController {
    
    private static final Logger logger = LoggerFactory.getLogger(WebTestingController.class);

    @Autowired
    private WebTestingService webTestingService;

    @GetMapping("/check")
    public String check() {

        webTestingService.hitURL();

        return "Hello World";
    }
    

}
