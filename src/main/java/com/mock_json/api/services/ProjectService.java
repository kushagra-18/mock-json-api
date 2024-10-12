package com.mock_json.api.services;

import org.springframework.stereotype.Service;

import com.mock_json.api.exceptions.NotFoundException;
import com.mock_json.api.models.Project;
import com.mock_json.api.repositories.ProjectRepository;

@Service
public class ProjectService {

    private final ProjectRepository projectRepository;

    public ProjectService(ProjectRepository projectRepository) {
        this.projectRepository = projectRepository;
    }

    public Project findProjectById(Long projectId) {
        return projectRepository.findById(projectId)
                .orElseThrow(() -> new NotFoundException("Project with ID " + projectId + " not found"));
    }
}
