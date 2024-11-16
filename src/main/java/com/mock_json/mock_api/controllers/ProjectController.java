package com.mock_json.mock_api.controllers;

import javax.validation.Valid;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import com.mock_json.mock_api.models.Project;
import com.mock_json.mock_api.models.Team;
import com.mock_json.mock_api.repositories.ProjectRepository;
import com.mock_json.mock_api.repositories.TeamRepository;
import com.mock_json.mock_api.exceptions.NotFoundException;
import com.mock_json.mock_api.constants.PusherChannels;
import com.mock_json.mock_api.constants.ResponseMessages;

import java.time.LocalDateTime;
import java.util.UUID;

@RestController
@RequestMapping("/api/v1/project")
public class ProjectController {

    private static final Logger logger = LoggerFactory.getLogger(ProjectController.class);

    @Autowired
    private ProjectRepository projectRepository;
    @Autowired
    private TeamRepository teamRepository;

    @PostMapping("create-free")
    public ResponseEntity<?> createFreeProject(@Valid @RequestBody Project project) {

        Team team = teamRepository.findById(1L)
                .orElseThrow(() -> new RuntimeException("Team not found"));

        String channelId = UUID.randomUUID().toString();

        project.setCreatedAt(LocalDateTime.now());
        project.setUpdatedAt(LocalDateTime.now());
        project.setChannelId(channelId);
        project.setTeam(team);

        Project savedProject = projectRepository.save(project);

        return ResponseEntity.ok(savedProject);
    }

    @GetMapping("{projectSlug}")
    public ResponseEntity<?> getProjectBySlug(@PathVariable String projectSlug) {
        
        if (projectSlug == null || projectSlug.trim().isEmpty()) {
            return ResponseEntity.badRequest().body("Project slug is required");
        }

        Project project = projectRepository.findBySlug(projectSlug)
                .orElseThrow(() -> new NotFoundException(ResponseMessages.NO_PROJECT));

        String channelName = PusherChannels.REQUEST_CHANNEL + project.getChannelId();

        project.setChannelId(channelName);

        return ResponseEntity.ok(project);
    }

}
