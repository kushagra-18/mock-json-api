package com.mock_json.testing_api.controllers;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;


import com.mock_json.testing_api.services.WebTestingService;


import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.web.bind.annotation.RequestMapping;


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

        // webTestingService.hitURL();

        return "Hello World";
    }
    

}
