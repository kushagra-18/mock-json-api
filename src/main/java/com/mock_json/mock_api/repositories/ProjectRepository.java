package com.mock_json.mock_api.repositories;

import java.util.Optional;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import com.mock_json.mock_api.models.Project;

@Repository
public interface ProjectRepository extends JpaRepository<Project, Long> {

    Optional<Project> findBySlug(String slug);

    Optional<Project> findByTeamSlugAndSlug(String teamSlug, String projectSlug);


}
    