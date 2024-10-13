package com.mock_json.api.repositories;

import java.util.Optional;

import org.springframework.data.jpa.repository.EntityGraph;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.rest.core.annotation.RepositoryRestResource;
import com.mock_json.api.models.Url;

@RepositoryRestResource
public interface UrlRepository extends JpaRepository<Url, Long> {

    @EntityGraph(attributePaths = {"project", "jsonList"})
    Optional<Url> findByUrl(String url);

}
    