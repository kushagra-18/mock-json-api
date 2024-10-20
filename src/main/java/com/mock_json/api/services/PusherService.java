package com.mock_json.api.services;

import com.pusher.rest.Pusher;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

import java.util.Collections;

@Service
public class PusherService {

    private final Pusher pusher;

    private static final Logger logger = LoggerFactory.getLogger(PusherService.class);

    public PusherService(
            @Value("${PUSHER_APP_ID}") String appId,
            @Value("${PUSHER_KEY}") String appKey,
            @Value("${PUSHER_SECRET}") String appSecret,
            @Value("${PUSHER_CLUSTER}") String appCluster) {

        pusher = new Pusher(appId, appKey, appSecret);
        pusher.setCluster(appCluster);
        pusher.setEncrypted(true);
    }

    public void trigger(String channel, String event, Object data) {
        pusher.trigger(channel, event, Collections.singletonMap("message", data));
    }
}
