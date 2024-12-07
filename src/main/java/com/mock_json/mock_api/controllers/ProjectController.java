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
import com.mock_json.mock_api.helpers.StringHelpers;
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

    private static final Long DEFAULT_TEAM_ID = 1L;

    @Autowired
    private ProjectRepository projectRepository;
    @Autowired
    private TeamRepository teamRepository;

    @Autowired
    private RequestLogService requestLogService;

    /**
     * Create a free project. Validates slug and assigns it to a default team.
     * @param project the project to create
     * @return ResponseEntity with the created project or an existing project if
     *         slug is duplicate
     */

    @PostMapping("free")
    public ResponseEntity<?> createFreeProject(@Valid @RequestBody Project project) {
        
        System.out.println("Creating a free project");

        validateSlug(project.getSlug());

        Optional<Project> existingProject = projectRepository.findBySlug(project.getSlug());
        
        if (existingProject.isPresent()) {
            return ResponseEntity.ok(existingProject.get());
        }

        Project savedProject = saveProjectWithTeam(project, project.getSlug(),project.getName());
       
        return ResponseEntity.status(HttpStatus.CREATED).body(savedProject);
    }

    /**
     * Create a free project with a random slug and channel ID.
     * @return ResponseEntity with the created project or an existing project if
     *         slug is duplicate
     */

    @PostMapping("free/fast-forward")
    public ResponseEntity<?> createFreeFastForwardProject() {
       
        String randomSlug = StringHelpers.generateRandomString(10);

        Optional<Project> existingProject = projectRepository.findBySlug(randomSlug);
       
        if (existingProject.isPresent()) {
            return ResponseEntity.ok(existingProject.get());
        }

        Project project = new Project();

        String projectName = StringHelpers.unslug(randomSlug);
        
        Project savedProject = saveProjectWithTeam(project, randomSlug,projectName);
        
        return ResponseEntity.status(HttpStatus.CREATED).body(savedProject);
    }


    /**
     * Validates the project slug, ensuring it is not in the disallowed list.
     * @param slug the project slug to validate
     */
    private void validateSlug(String slug) {
        if (slug == null || slug.isBlank()) {
            throw new BadRequestException("Slug cannot be null or empty.");
        }

        List<String> disallowedSlugs = Arrays.asList(DisallowedProjectSlugs.DISALLOWED_SLUGS);
        if (disallowedSlugs.contains(slug.toLowerCase())) {
            throw new BadRequestException(ResponseMessages.RESTRICT_PROJECT_SLUG);
        }
    }

    /**
     * Saves the project with a given slug and associates it with the default team.
     *
     * @param project the project to save
     * @param slug    the slug for the project
     * @return the saved project
     */
    private Project saveProjectWithTeam(Project project, String slug,String name) {
        Team team = teamRepository.findById(DEFAULT_TEAM_ID)
                .orElseThrow(() -> new RuntimeException("Default team not found"));

        String channelId = UUID.randomUUID().toString();

        project.setSlug(slug);
        project.setChannelId(channelId);
        project.setTeam(team);
        project.setName(name);
        project.setIsForwardProxyActive(false);
        project.setCreatedAt(LocalDateTime.now());
        project.setUpdatedAt(LocalDateTime.now());

        return projectRepository.save(project);
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
