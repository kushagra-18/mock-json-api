package com.mock_json.api.controllers;

import com.mock_json.api.models.Project;
import com.mock_json.api.models.Team;
import com.mock_json.api.repositories.ProjectRepository; 
import com.mock_json.api.repositories.TeamRepository;

import javax.validation.Valid;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.time.LocalDateTime;

@RestController
@RequestMapping("/api/v1/project")
public class ProjectController {

    private static final Logger logger = LoggerFactory.getLogger(ProjectController.class);

    @Autowired
    private ProjectRepository projectRepository;
    @Autowired
    private TeamRepository teamRepository;

    @PostMapping("/create-free")
    public ResponseEntity<?> createFreeProject(@Valid @RequestBody Project project) {
       
        Team team = teamRepository.findById(1L)
        .orElseThrow(() -> new RuntimeException("Team not found"));

        project.setCreatedAt(LocalDateTime.now());
        project.setUpdatedAt(LocalDateTime.now());
        project.setTeam(team);

        Project savedProject = projectRepository.save(project);

        logger.info("Free project created: {}", savedProject);

        return ResponseEntity.ok(savedProject);
    }
}
