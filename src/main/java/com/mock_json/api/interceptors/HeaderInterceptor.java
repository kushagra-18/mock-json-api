package com.mock_json.api.interceptors;

import org.springframework.stereotype.Component;
import org.springframework.web.servlet.HandlerInterceptor;

import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import com.mock_json.api.contexts.HeaderContext;
import org.springframework.lang.NonNull;

@Component
public class HeaderInterceptor implements HandlerInterceptor {

    @Override
    public boolean preHandle(@NonNull HttpServletRequest request, @NonNull HttpServletResponse response,
            @NonNull Object handler)
            throws Exception {

        String teamHeader = request.getHeader("X-header-team");

        String projectHeader = request.getHeader("X-header-project");

        if (teamHeader == null || projectHeader == null) {
            response.setStatus(HttpServletResponse.SC_BAD_REQUEST);
            response.getWriter().write("Missing required headers: X-header-team or X-header-project");
            return false;
        }

        HeaderContext.setTeamId(teamHeader);
        HeaderContext.setProjectId(projectHeader);

        return true;
    }

    @Override
    public void afterCompletion(@NonNull HttpServletRequest request,
            @NonNull HttpServletResponse response,
            @NonNull Object handler,
            Exception ex) {
        HeaderContext.clear();
    }
}
