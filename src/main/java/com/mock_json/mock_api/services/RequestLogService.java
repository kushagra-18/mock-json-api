package com.mock_json.mock_api.services;

import java.util.HashMap;
import java.util.Map;

import org.springframework.scheduling.annotation.Async;
import org.springframework.stereotype.Service;

import com.mock_json.mock_api.models.Project;
import com.mock_json.mock_api.models.RequestLog;
import com.mock_json.mock_api.repositories.RequestLogRepository;


@Service
public class RequestLogService {

    private final RequestLogRepository requestLogRepository;
    private final PusherService pusherService;

    public RequestLogService(RequestLogRepository requestLogRepository, PusherService pusherService) {
        this.requestLogRepository = requestLogRepository;
        this.pusherService = pusherService;
    }

    @Async
    public void saveRequestLogAsync(String url, String method,String ip, int status, Project project) {

        RequestLog requestLog = RequestLog.builder()
                .url(url)
                .method(method)
                .ip(ip)
                .status(status)
                // .project(project)
                .build();

        requestLogRepository.save(requestLog);
    }
    
    @Async
    public void emitPusherEvent(String method, String url, Project project, int status) {
       
        Map<String, Object> data = new HashMap<>();
        data.put("method", method);
        // data.put("url", url);
        // // data.put("projectId", project.getId());
        // // data.put("projectName", project.getName()); 
        data.put("status", status);

        pusherService.trigger("project-events", "project-created", data);
    }

}
