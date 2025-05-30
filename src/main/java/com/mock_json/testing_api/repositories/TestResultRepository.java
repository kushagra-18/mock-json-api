package com.mock_json.testing_api.repositories;

import com.mock_json.testing_api.entities.TestResultEntity;
import org.springframework.data.mongodb.repository.MongoRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface TestResultRepository extends MongoRepository<TestResultEntity, String> {
}
