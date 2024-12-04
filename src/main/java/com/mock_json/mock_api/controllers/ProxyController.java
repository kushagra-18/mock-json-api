package com.mock_json.mock_api.controllers;

import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.servlet.support.ServletUriComponentsBuilder;

import com.mock_json.mock_api.dtos.UrlDataDto;
import com.mock_json.mock_api.exceptions.ConflictException;
import com.mock_json.mock_api.exceptions.NotFoundException;
import com.mock_json.mock_api.models.ForwardProxy;
import com.mock_json.mock_api.models.MockContent;
import com.mock_json.mock_api.models.Project;
import com.mock_json.mock_api.models.Url;
import com.mock_json.mock_api.services.ProjectService;
import com.mock_json.mock_api.services.ProxyService;
import com.mock_json.mock_api.services.UrlService;

import jakarta.persistence.EntityNotFoundException;
import jakarta.transaction.Transactional;
import jakarta.validation.Valid;

import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.web.bind.annotation.PatchMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;

import java.net.URI;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.Map;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;

@RestController
@ResponseBody
@RequestMapping("proxy")
@Validated
public class ProxyController {

    @Autowired
    private ProxyService proxyService;

    @Autowired
    private ProjectService projectService;

    @PostMapping("forward")
    @Transactional
    public ResponseEntity<?> saveForwardProxy(@Valid @RequestBody ForwardProxy forwardProxy) {

        ForwardProxy existingForwardProxy = proxyService.getForwardProxyByProjectId(forwardProxy.getProjectId());

        if (existingForwardProxy != null) {
            throw new ConflictException("A Proxy already exists for this project");
        }

        Project project = projectService.getDataById(forwardProxy.getProjectId());

        if (project == null) {
            throw new EntityNotFoundException("Project not found");
        }

        // save forward proxy
        ForwardProxy savedForwardProxy = proxyService.saveForwardProxy(forwardProxy, project);

        return ResponseEntity.status(HttpStatus.CREATED).body(savedForwardProxy);
    }

    @PatchMapping("forward/active/{projectId}")
    @Transactional
    public ResponseEntity<?> updateForwardProxyActiveStatus(@PathVariable Long projectId) {

        Project project = projectService.getDataById(projectId);

        if (project == null) {
            throw new NotFoundException("Project not found");
        }

        Boolean status = project.getIsForwardProxyActive();

        Boolean newStatus = !status;

        Project updatedProject = projectService.updateForwardFroxyActiveStatus(project, newStatus);

        return ResponseEntity.ok(updatedProject);

    }
}
