package com.mock_json.api.repositories;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;
import com.mock_json.api.models.Json;

@Repository
public interface JsonRepository extends JpaRepository<Json, Long> {



}
    