package com.mock_json.mock_api.services;

import java.util.HashMap;
import java.util.Map;
import java.util.Optional;

import org.springframework.scheduling.annotation.Async;
import org.springframework.stereotype.Service;

import com.mock_json.mock_api.constants.PusherChannels;
import com.mock_json.mock_api.dtos.PusherRequestEventDto;
import com.mock_json.mock_api.models.Project;
import com.mock_json.mock_api.models.RequestLog;
import com.mock_json.mock_api.models.Url;
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
    public void saveRequestLogAsync(String ip, Url url,String method, String urlString, int status) {

        Long urlId = (url != null) ? url.getId() : -1L;

        RequestLog requestLog = RequestLog.builder()
                .ip(ip)
                .urlId(urlId)
                .url(urlString)
                .method(method)
                .status(status)
                .build();

        requestLogRepository.save(requestLog);
    }
    
    @Async
    public void emitPusherEvent(String ip, Url url,String method, String urlString, int status,String channelId) {
       
        if (ip == null || url == null  || channelId == null) {
            return;
        }
    
        Long urlId = Optional.ofNullable(url)
                             .map(Url::getId)
                             .orElse(-1L);  
    
        PusherRequestEventDto eventData = new PusherRequestEventDto(method, urlString, urlId, ip, status);

        String channelName = PusherChannels.REQUEST_CHANNEL + channelId;

        pusherService.trigger(channelName, "mock-url-created", eventData);
    }

}
