package com.mock_json.api.services;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.data.redis.core.RedisTemplate;
import org.springframework.stereotype.Service;
import java.time.Duration;

@Service
public class RedisService {

    @Autowired
    private RedisTemplate<String, Object> redisTemplate;

    @Value("${REDIS_PREFIX}")
    private String redisPrefix;

    public void saveData(String key, Object value) {
        redisTemplate.opsForValue().set(key, value);
    }

    public Object getData(String key) {
        return redisTemplate.opsForValue().get(key);
    }

    public void deleteData(String key) {
        redisTemplate.delete(key);
    }

    public Long incrementValue(String key, long incrementBy) {
        return redisTemplate.opsForValue().increment(key, incrementBy);
    }

    /**
     * Sets an expiration time for a Redis key.
     * @param key The Redis key
     * @param duration The expiration duration
     */
    public void expire(String key, Duration duration) {
        redisTemplate.expire(key, duration);
    }

    /**
     * Creates a Redis key from the given arguments.
     * @param args
     * @return
     */
    public String createRedisKey(String... args) {
        
        StringBuilder result = new StringBuilder();

        result.append(redisPrefix);
        result.append("#");

        for (int i = 0; i < args.length; i++) {
            result.append(args[i]);
            if (i < args.length - 1) {
                result.append("#");
            }
        }

        return result.toString().replaceAll(":", "_");
    }
}
