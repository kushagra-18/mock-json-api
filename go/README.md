# Mock API - Go Version

This directory contains a rewrite of the Mock API application in Go, using the Gin framework for routing and GORM for database interaction.

## Prerequisites

Before running this application, ensure you have the following installed:

*   Go (version 1.20 or higher recommended)
*   PostgreSQL (version 12 or higher recommended)
*   Redis (version 5 or higher recommended)

## Configuration

Configuration for the application is managed through environment variables, typically loaded from a `.env` file.

1.  Copy the example environment file:
    ```bash
    cp env.example .env
    ```
2.  Edit the `.env` file with your specific configuration details for:
    *   Database connection (DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)
    *   Redis connection (REDIS_ADDR, REDIS_PASSWORD, REDIS_DB)
    *   JWT settings (JWT_SECRET_KEY, JWT_EXPIRATION_HOURS)
    *   Application settings (BASE_URL, SERVER_PORT)
    *   Global rate limiting parameters (GLOBAL_MAX_ALLOWED_REQUESTS, GLOBAL_TIME_WINDOW_SECONDS)
    # Gemini API Configuration
    GEMINI_API_KEY=your_gemini_api_key_here
    GEMINI_MODEL_NAME=gemini-1.5-flash-latest # Or your preferred default

## Dependencies

Go modules are used for dependency management. To install or verify dependencies:

```bash
go mod tidy
# or
# go get ./...
```

## Database Setup

The application uses GORM's `AutoMigrate` feature to create and update the database schema based on the defined models.
Ensure that:
1.  Your PostgreSQL server is running.
2.  The database specified in your `.env` file (e.g., `mockapi_db`) exists.
3.  The user specified in your `.env` file has privileges to connect to and create tables in this database.

## Running the Application

There are two main ways to run the application:

1.  **For Development (using `go run`):**
    This command will compile and run the application directly from the source code.
    ```bash
    go run main.go
    ```

2.  **Building a Binary (for deployment/production):**
    First, build the executable:
    ```bash
    go build -o mockapi_server main.go
    ```
    Then, run the compiled binary:
    ```bash
    ./mockapi_server
    ```

The server will start on the port specified by `SERVER_PORT` in your `.env` file (default is `8080`).

## API Endpoints

This Go application aims to replicate the functionality and API endpoints of the original Java-based Mock API. Please refer to the existing API documentation or controller implementations for details on available endpoints. Key functionalities include:
*   Project creation and management.
*   Mock content definition and retrieval.
*   URL configuration for mocks.
*   Forward proxy capabilities.
*   Rate limiting.

The main mock serving endpoint is accessible via:
`GET /mock/:teamSlug/:projectSlug/*actualMockPath?token=<your_jwt_token>`

### AI Prompting

A new endpoint is available for interacting with Google's Gemini AI:

*   **Endpoint:** `POST /api/v1/ai/prompt`
*   **Purpose:** Accepts a user-provided text prompt and returns a JSON response from the Gemini AI model.
*   **Request Body (JSON):**
    ```json
    {
        "prompt": "Your text prompt for the AI."
    }
    ```
*   **Response Body (JSON):**
    The endpoint returns a JSON object containing the processed response from the Gemini API.
    ```json
    {
        "text": "The AI's generated text response.",
        "finish_reason": "STOP"
        // Other fields like "token_count" might be included in the future.
    }
    ```

Refer to the `routes/routes.go` file for a complete list of registered routes and their corresponding controller handlers.
