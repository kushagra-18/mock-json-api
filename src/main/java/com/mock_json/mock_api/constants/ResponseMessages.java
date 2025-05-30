package com.mock_json.mock_api.constants;


public class ResponseMessages {

    private ResponseMessages() {
    }

    public static final String SUCCESS_MESSAGE = "Operation completed successfully.";
    
    // Error Messages
    public static final String RESOURCE_NOT_FOUND = "The requested resource was not found.";
    public static final String NO_URL_PRESENT = "Hey, nothing in here!";
    public static final String RATE_LIMIT_EXCEEDED = "Rate limit exceeded, Please try again later.";
    public static final String JSON_PARSE_ERROR = "Error parsing JSON data";
    public static final String GLOBAL_RATE_LIMIT_EXCEEDED = "You have exhaused the number of requests allowed for this endpoint. Please contact us if you need more.";
    
    // Validation Messages
    public static final String INVALID_INPUT = "The input provided is invalid.";
    public static final String UNAUTHORIZED_ACCESS = "You are not authorized to perform this action.";
    public static final String RESTRICT_PROJECT_SLUG = "This is a restricted project, please use a different slug, If you are working in one of these projects, please reach out to us.";

    // 
    public static final String NO_CONTENT_URL = "Oopsie! 🐾 Looks like you stumbled upon an uncharted path. No mock API here yet—let’s create one and bring this URL to life";
    public static final String NO_PROJECT = "Oops! We couldn't find that project. Are you sure it exists? 🧐";





   
}
