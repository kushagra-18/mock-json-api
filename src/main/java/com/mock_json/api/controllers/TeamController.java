package com.mock_json.api.controllers;

import java.util.HashMap;
import java.util.Optional;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.ResponseBody;

@RestController
@RequestMapping("/api/v1/team")
@ResponseBody
public class TeamController {

    private static final Logger logger = LoggerFactory.getLogger(HomeController.class);

    // @Autowired
    // private UserService userService;

    @GetMapping()
    public String home() {

        HashMap<String, String> response = new HashMap<>();

        String message = "Hello Wddorld, id is kushagra";

        logger.info("hello world");
        logger.error("this is error");

        response.put("message", message);

        // return userService.findById("668aa7acd1ce0a25e3aa194f");

        return "Hello World";
    }

}
