package com.mock_json.mock_api.services;

import org.springframework.stereotype.Service;

import com.mock_json.mock_api.models.ForwardProxy;
import com.mock_json.mock_api.models.Project;
import com.mock_json.mock_api.repositories.ForwardProxyRepository;

@Service
public class ProxyService {

    private final ForwardProxyRepository forwardProxyRepository;

    public ProxyService(ForwardProxyRepository forwardProxyRepository) {
        this.forwardProxyRepository = forwardProxyRepository;
    }

    public ForwardProxy saveForwardProxy(ForwardProxy forwardProxy,Project project) {

        forwardProxy.setProject(project);

        return forwardProxyRepository.save(forwardProxy);
    }

    public ForwardProxy getForwardProxyByProjectId(Long projectId) {
        return forwardProxyRepository.findByProject_Id(projectId);
    }
}
