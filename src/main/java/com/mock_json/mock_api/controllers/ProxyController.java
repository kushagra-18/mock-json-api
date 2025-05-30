package com.mock_json.mock_api.controllers;

import org.springframework.web.bind.annotation.RestController;

import com.mock_json.mock_api.dtos.ForwardProxyDto;
import com.mock_json.mock_api.exceptions.ConflictException;
import com.mock_json.mock_api.exceptions.NotFoundException;
import com.mock_json.mock_api.models.ForwardProxy;
import com.mock_json.mock_api.models.Project;
import com.mock_json.mock_api.services.ProjectService;
import com.mock_json.mock_api.services.ProxyService;

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


import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;

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
    public ResponseEntity<?> saveForwardProxy(@Valid @RequestBody ForwardProxyDto forwardProxy) {

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


        if (project.getIsForwardProxyActive() == null) {
            project.setIsForwardProxyActive(false);
            projectService.save(project);
        }else{
            project.setIsForwardProxyActive(true);
            projectService.save(project);
        }

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
