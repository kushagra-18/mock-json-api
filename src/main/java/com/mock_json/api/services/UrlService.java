package com.mock_json.api.services;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.jpa.domain.Specification;
import org.springframework.stereotype.Service;

import com.mock_json.api.helpers.StringHelpers;
import com.mock_json.api.models.Project;
import com.mock_json.api.models.Url;
import com.mock_json.api.repositories.UrlRepository;
import com.mock_json.api.specfications.UrlSpecifications;
import jakarta.servlet.http.HttpServletRequest;

import java.time.LocalDateTime;
import java.util.Optional;
import java.time.Duration;

@Service
public class UrlService {

    private final UrlRepository urlRepository;

    @Autowired
    private RedisService redisService;

    private static final Logger logger = LoggerFactory.getLogger(MockContentService.class);

    public UrlService(UrlRepository urlRepository) {
        this.urlRepository = urlRepository;
    }

    public Optional<Url> findUrlDataByUrlAndTeamAndProject(String teamSlug, String projectSlug, String url) {
        Specification<Url> spec = UrlSpecifications.hasTeamSlugAndProjectSlugAndUrl(teamSlug, projectSlug, url);
        return urlRepository.findOne(spec);
    }

    public Optional<Url> findUrlDataByUrl(String url) {
        return urlRepository.findByUrl(url);
    }

    /**
     * Returns the full URL with query string.
     * which was recieved in the request
     * 
     * @param request
     */
    public String getUrl(HttpServletRequest request) {

        StringBuilder fullPathWithQuery = new StringBuilder(request.getRequestURI());

        String queryString = request.getQueryString();

        if (queryString != null) {
            fullPathWithQuery.append("?").append(queryString);
        }

        return StringHelpers.ltrim(fullPathWithQuery.toString(), '/');
    }

    public Url saveData(Url url, Project project) {
        LocalDateTime currTime = LocalDateTime.now();

        url.setCreatedAt(currTime);
        url.setUpdatedAt(currTime);
        url.setProject(project);

        return urlRepository.save(url);
    }

    /**
     * Checks if the IP has exceeded the rate limit for the given URL.
     * 
     * @param ip  The IP address of the request.
     * @param url The URL object containing rate limit information.
     * @return true if the rate limit is exceeded, false otherwise.
     */
    public boolean isRateLimited(String ip, String url, Integer allowedRequests, Long timeWindow) {

        if (allowedRequests != null && timeWindow != null) {

            String sanitizedIp = ip.replaceAll(":", ".");

            String redisKey = redisService.createRedisKey("rate_limit",sanitizedIp, url);

            Long requestCount = redisService.incrementValue(redisKey, 1);

            long count = (requestCount != null) ? requestCount : 0L; 

            if (count == 1) {
                redisService.expire(redisKey, Duration.ofSeconds(timeWindow));
            }

            return count > allowedRequests;
        }

        return false;
    }

}
