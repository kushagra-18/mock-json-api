package com.mock_json.mock_api.filters;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Component;

import jakarta.servlet.Filter;
import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.ServletRequest;
import jakarta.servlet.ServletResponse;
import jakarta.servlet.http.HttpServletRequest;
import java.io.IOException;


@Component
public class RequestLoggingFilter implements Filter {

    private static final Logger logger = LoggerFactory.getLogger(RequestLoggingFilter.class);

    @Override
    public void doFilter(ServletRequest request, ServletResponse response, FilterChain chain)
            throws IOException, ServletException {
        if (request instanceof HttpServletRequest) {
            HttpServletRequest httpServletRequest = (HttpServletRequest) request;
            logger.info("Incoming request: {} {}", httpServletRequest.getMethod(), httpServletRequest.getRequestURI());
            logger.info("Headers:");
            httpServletRequest.getHeaderNames().asIterator().forEachRemaining(header ->
                logger.info("{}: {}", header, httpServletRequest.getHeader(header))
            );
        }
        chain.doFilter(request, response);
    }
}