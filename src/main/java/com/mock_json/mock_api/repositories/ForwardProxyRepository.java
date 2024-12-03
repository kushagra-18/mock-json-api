package com.mock_json.mock_api.repositories;


import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.JpaSpecificationExecutor;
import org.springframework.data.rest.core.annotation.RepositoryRestResource;

import com.mock_json.mock_api.models.ForwardProxy;
import com.mock_json.mock_api.models.Url;

@RepositoryRestResource

public interface ForwardProxyRepository extends JpaRepository<ForwardProxy, Long>, JpaSpecificationExecutor<Url> {
    
    ForwardProxy findByProject_Id(Long projectId);

}
    