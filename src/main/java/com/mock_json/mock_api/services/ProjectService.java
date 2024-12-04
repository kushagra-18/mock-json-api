package com.mock_json.mock_api.services;

import org.springframework.stereotype.Service;

import com.mock_json.mock_api.exceptions.NotFoundException;
import com.mock_json.mock_api.models.Project;
import com.mock_json.mock_api.repositories.ProjectRepository;

@Service
public class ProjectService {

    private final ProjectRepository projectRepository;

    public ProjectService(ProjectRepository projectRepository) {
        this.projectRepository = projectRepository;
    }

    public Project getDataById(Long projectId) {
        return projectRepository.findById(projectId)
                .orElseThrow(() -> new NotFoundException("Project with ID " + projectId + " not found"));
    }

    public Project getDataBySlugAndTeamSlug(String teamSlug, String projectSlug) {
        return projectRepository.findByTeamSlugAndSlug(teamSlug, projectSlug)
                .orElseThrow(() -> new NotFoundException("Project with slug " + projectSlug + " not found"));
    }

    public Project findProjectById(Long projectId) {
        return projectRepository.findById(projectId)
                .orElseThrow(() -> new NotFoundException("Project with ID " + projectId + " not found"));
    }

    public Project findProjectBySlug(String projectSlug) {
        return projectRepository.findBySlug(projectSlug)
                .orElseThrow(() -> new NotFoundException("Project with slug " + projectSlug + " not found"));
    }

    public Project updateForwardFroxyActiveStatus(Project project, boolean status) {
       
        project.setIsForwardProxyActive(status);
        
        return projectRepository.save(project);
    }
}
