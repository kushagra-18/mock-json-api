package com.mock_json.mock_api.controllers;

import javax.validation.Valid;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import com.mock_json.mock_api.models.Project;
import com.mock_json.mock_api.models.RequestLog;
import com.mock_json.mock_api.models.Team;
import com.mock_json.mock_api.repositories.ProjectRepository;
import com.mock_json.mock_api.repositories.TeamRepository;
import com.mock_json.mock_api.services.RequestLogService;
import com.mock_json.mock_api.exceptions.BadRequestException;
import com.mock_json.mock_api.exceptions.NotFoundException;
import com.mock_json.mock_api.constants.DisallowedProjectSlugs;
import com.mock_json.mock_api.constants.PusherChannels;
import com.mock_json.mock_api.constants.ResponseMessages;

import java.time.LocalDateTime;
import java.util.Arrays;
import java.util.HashMap;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

@RestController
@RequestMapping("project")
public class ProjectController {

    private static final Logger logger = LoggerFactory.getLogger(ProjectController.class);

    @Autowired
    private ProjectRepository projectRepository;
    @Autowired
    private TeamRepository teamRepository;

    @Autowired
    private RequestLogService requestLogService;

    @PostMapping("create-free")
    public ResponseEntity<?> createFreeProject(@Valid @RequestBody Project project) {

        List<String> disallowedSlugs = Arrays.asList(DisallowedProjectSlugs.DISALLOWED_SLUGS);

        if (disallowedSlugs.contains(project.getSlug().toLowerCase())) {
            throw new BadRequestException(ResponseMessages.RESTRICT_PROJECT_SLUG);
        }

        Optional<Project> existingProject = projectRepository.findBySlug(project.getSlug());

        if (existingProject.isPresent()) {
            return ResponseEntity.ok(existingProject.get());
        }

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

    @GetMapping("/{projectSlug}")
    public ResponseEntity<?> getProjectBySlug(
            @PathVariable String projectSlug,
            @RequestParam(defaultValue = "10") Integer limit,
            @RequestParam(defaultValue = "0") Integer offset) {

        if (projectSlug == null || projectSlug.trim().isEmpty()) {
            return ResponseEntity.badRequest().body("Project slug is required");
        }

        Project project = projectRepository.findBySlug(projectSlug)
                .orElseThrow(() -> new NotFoundException(ResponseMessages.NO_PROJECT));

        List<RequestLog> requestLogs = requestLogService.getLogsByProjectId(project.getId(), limit, offset);

        project.setChannelId(PusherChannels.REQUEST_CHANNEL + project.getChannelId());

        HashMap<String, Object> response = buildResponse(project, requestLogs);

        return ResponseEntity.ok(response);
    }

    private HashMap<String, Object> buildResponse(Project project, List<RequestLog> requestLogs) {
        HashMap<String, Object> response = new HashMap<>();
        response.put("project", project);
        response.put("request_logs", requestLogs);
        response.put("status_code", HttpStatus.OK.value());
        return response;
    }
}
