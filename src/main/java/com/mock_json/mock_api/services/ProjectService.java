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

    public Project findProjectById(Long projectId) {
        return projectRepository.findById(projectId)
                .orElseThrow(() -> new NotFoundException("Project with ID " + projectId + " not found"));
    }
}
