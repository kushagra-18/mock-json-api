package com.mock_json.mock_api.services;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.jpa.domain.Specification;
import org.springframework.stereotype.Service;

import com.mock_json.mock_api.dtos.UrlDataDto;
import com.mock_json.mock_api.helpers.StringHelpers;
import com.mock_json.mock_api.models.Project;
import com.mock_json.mock_api.models.Url;
import com.mock_json.mock_api.repositories.UrlRepository;
import com.mock_json.mock_api.specfications.UrlSpecifications;

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


    public Optional<Url> findById(Long id) {
        return urlRepository.findById(id);
    }

    public Optional<Url> findUrlDataByUrlAndTeamAndProject(String teamSlug, String projectSlug, String url) {
        Specification<Url> spec = UrlSpecifications.hasTeamSlugAndProjectSlugAndUrl(teamSlug, projectSlug, url);
        return urlRepository.findOne(spec);
    }

    public Optional<Url> findUrlDataByUrlAndProjectSlug(String projectSlug, String url) {
        Specification<Url> spec = UrlSpecifications.hasProjectSlugAndUrl(projectSlug, url);
        return urlRepository.findOne(spec);
    }

    public Url save(Url existingUrl, UrlDataDto urlDto) {

        if (urlDto.getDescription() != null) {
            existingUrl.setDescription(urlDto.getDescription());
        }
        if (urlDto.getName() != null) {
            existingUrl.setName(urlDto.getName());
        }
        if (urlDto.getRequests() != null) {
            existingUrl.setRequests(urlDto.getRequests());
        }
        if (urlDto.getTime() != null) {
            existingUrl.setTime(urlDto.getTime());
        }

        return urlRepository.save(existingUrl);
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

        String sanitizedIp = ip.replaceAll(":", ".");

        String redisKey = redisService.createRedisKey("rate_limit_custom", sanitizedIp, url);

        return redisService.rateLimit(redisKey, allowedRequests, timeWindow);
    }

}
