package com.mock_json.mock_api.repositories;

import java.util.List;

import org.springframework.data.domain.Pageable;
import org.springframework.data.mongodb.repository.MongoRepository;
import org.springframework.stereotype.Repository;

import com.mock_json.mock_api.models.RequestLog;

@Repository
public interface RequestLogRepository extends MongoRepository<RequestLog, String> {

    List<RequestLog> findByProjectId(Long projectId, Pageable pageable);

    // void updateRequestByUrlIdAndUrl(Long urlId, String url);

}
