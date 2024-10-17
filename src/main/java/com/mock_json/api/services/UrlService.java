package com.mock_json.api.services;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.data.jpa.domain.Specification;
import org.springframework.stereotype.Service;

import com.mock_json.api.exceptions.NotFoundException;
import com.mock_json.api.helpers.StringHelpers;
import com.mock_json.api.models.MockContent;
import com.mock_json.api.models.Project;
import com.mock_json.api.models.Url;
import com.mock_json.api.repositories.MockContentRepository;
import com.mock_json.api.repositories.UrlRepository;
import com.mock_json.api.specfications.UrlSpecifications;

import jakarta.servlet.http.HttpServletRequest;

import java.net.URI;
import java.net.URISyntaxException;
import java.time.LocalDateTime;
import java.util.List;
import java.util.Optional;
import java.util.concurrent.TimeUnit;

@Service
public class UrlService {

    private final UrlRepository urlRepository;

    private static final Logger logger = LoggerFactory.getLogger(MockContentService.class);

    public UrlService(UrlRepository urlRepository) {
        this.urlRepository = urlRepository;
    }

    public boolean checkURLExists(String url) {
        // TODO: Implement this method to check if the URL exists
        return true;
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
}
