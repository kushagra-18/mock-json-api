package com.mock_json.mock_api.interceptors;
// package com.mock_json.api.interceptors;

// import org.springframework.beans.factory.annotation.Autowired;
// import org.springframework.web.bind.annotation.ControllerAdvice;
// import org.springframework.web.bind.annotation.ModelAttribute;

// import com.mock_json.api.controllers.MockContentController;

// import jakarta.servlet.http.HttpServletRequest;
// import jakarta.servlet.http.HttpServletResponse;

// @ControllerAdvice
// public class MockRouteInterceptor {
        
//     @Autowired
//     private MockContentController mockApiController;

//     @ModelAttribute
//     public void handleMockApiRequests(HttpServletRequest request, HttpServletResponse response) {
//         // String mockApiHeader = request.getHeader("X-mock-api");
//         String mockApiHeader = "true";

//         if ("true".equalsIgnoreCase(mockApiHeader)) {
//             mockApiController.getMockedJSON(request);
//         }
//     }
// }
