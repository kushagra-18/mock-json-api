package com.mock_json.mock_api.repositories;

import java.util.Optional;

import org.springframework.data.jpa.domain.Specification;
import org.springframework.data.jpa.repository.EntityGraph;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.JpaSpecificationExecutor;
import org.springframework.data.rest.core.annotation.RepositoryRestResource;

import com.mock_json.mock_api.models.Url;

@RepositoryRestResource

public interface UrlRepository extends JpaRepository<Url, Long>, JpaSpecificationExecutor<Url> {
   
    @EntityGraph(attributePaths = {"project", "mockContentList"})
    Optional<Url> findByUrl(Specification<Url> spec);

    // @EntityGraph(attributePaths = {"project", "mockContentList"})
    // Optional<Url> findByUrl(String spec);

    @EntityGraph(attributePaths = {"project", "mockContentList"})
    Optional<Url> findOne(Specification<Url> spec);

    

}
    