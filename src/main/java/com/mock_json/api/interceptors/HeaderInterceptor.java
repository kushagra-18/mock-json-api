package com.mock_json.api.interceptors;

import org.springframework.stereotype.Component;
import org.springframework.web.servlet.HandlerInterceptor;

import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import com.mock_json.api.contexts.HeaderContext;
import org.springframework.lang.NonNull;
import com.mock_json.api.annotations.HeaderIntercepted;
import org.springframework.web.method.HandlerMethod;

@Component
public class HeaderInterceptor implements HandlerInterceptor {

    @Override
    public boolean preHandle(@NonNull HttpServletRequest request, @NonNull HttpServletResponse response,
                             @NonNull Object handler) throws Exception {
        if (handler instanceof HandlerMethod) {
            HandlerMethod handlerMethod = (HandlerMethod) handler;

            if (handlerMethod.getMethod().isAnnotationPresent(HeaderIntercepted.class) ||
                handlerMethod.getBeanType().isAnnotationPresent(HeaderIntercepted.class)) {

                String teamHeader = request.getHeader("X-header-team");
                String projectHeader = request.getHeader("X-header-project");

                // String teamHeader = "free";
                // String projectHeader = "kush";

                if (teamHeader == null || projectHeader == null) {
                    response.setStatus(HttpServletResponse.SC_BAD_REQUEST);
                    response.getWriter().write("Missing required headers: X-header-team or X-header-project");
                    return false;
                }

                HeaderContext.setTeamSlug(teamHeader);
                HeaderContext.setProjectSlug(projectHeader);
            }
        }
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
