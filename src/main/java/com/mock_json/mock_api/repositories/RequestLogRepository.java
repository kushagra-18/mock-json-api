package com.mock_json.mock_api.repositories;


import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import com.mock_json.mock_api.models.RequestLog;

@Repository
public interface RequestLogRepository extends JpaRepository<RequestLog, Long> {


}
    