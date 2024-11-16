package com.mock_json.mock_api.specfications;

import org.springframework.data.jpa.domain.Specification;

import com.mock_json.mock_api.models.Url;

public class UrlSpecifications {

    public static Specification<Url> hasTeamSlug(String teamSlug) {
        return (root, query, criteriaBuilder) -> {
            return criteriaBuilder.equal(root.join("project").join("team").get("slug"), teamSlug);
        };
    }

    public static Specification<Url> hasProjectSlug(String projectSlug) {
        return (root, query, criteriaBuilder) -> {
            return criteriaBuilder.equal(root.join("project").get("slug"), projectSlug);
        };
    }

    public static Specification<Url> hasUrl(String url) {
        return (root, query, criteriaBuilder) -> {
            return criteriaBuilder.equal(root.get("url"), url);
        };
    }

    public static Specification<Url> hasProjectSlugAndUrl(String projectSlug, String url) {
        return hasProjectSlug(projectSlug)
                .and(hasUrl(url));
    }

    public static Specification<Url> hasTeamSlugAndProjectSlugAndUrl(String teamSlug, String projectSlug, String url) {
        return hasTeamSlug(teamSlug)
                .and(hasProjectSlug(projectSlug))
                .and(hasUrl(url));
    }
}