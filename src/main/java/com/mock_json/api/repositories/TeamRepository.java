package com.mock_json.api.repositories;

import java.util.Optional;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;
import com.mock_json.api.models.Project;
import com.mock_json.api.models.Team;

@Repository
public interface TeamRepository extends JpaRepository<Team, Long> {

    Optional<Project> findBySlug(String slug);


}
    