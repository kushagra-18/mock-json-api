package com.mock_json.mock_api.controllers;

import org.springframework.web.bind.annotation.RestController;

import com.mock_json.mock_api.dtos.UrlDataDto;
import com.mock_json.mock_api.exceptions.NotFoundException;
import com.mock_json.mock_api.models.MockContent;
import com.mock_json.mock_api.models.Url;
import com.mock_json.mock_api.services.UrlService;

import jakarta.persistence.EntityNotFoundException;
import jakarta.transaction.Transactional;
import jakarta.validation.Valid;

import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.web.bind.annotation.PatchMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.Map;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestParam;

@RestController
@ResponseBody
@RequestMapping("url")
@Validated
public class UrlController {

    @Autowired
    private UrlService urlService;

    @PatchMapping("{urlId}")
    @Transactional
    public ResponseEntity<Map<String, Object>> updateUrlInfo(
            @PathVariable Long urlId,
            @Valid @RequestBody UrlDataDto urlDto) {

        Map<String, Object> response = new HashMap<>();

        // Fetch the existing URL from the database
        Url existingUrl = urlService.findById(urlId)
                .orElseThrow(() -> new NotFoundException("URL not found with ID: " + urlId));

        Url updatedUrl = urlService.save(existingUrl, urlDto);

        updatedUrl.setMockContentList(null);

        response.put("data", updatedUrl);
        response.put("status_code", HttpStatus.OK.value());

        return ResponseEntity.ok(response);
    }

    @GetMapping("{urlId}")
    public ResponseEntity<Map<String, Object>> getURLDetailsWithMockContent(@PathVariable Long urlId) {

        Url url = urlService.findById(urlId)
                .orElseThrow(() -> new NotFoundException("URL not found with ID: " + urlId));

        Map<String, Object> response = new HashMap<>();

        response.put("url", url);
        response.put("status_code", HttpStatus.OK.value());

        return ResponseEntity.ok(response);

    }

}
