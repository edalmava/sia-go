# SIA - Sistema Integrado Académico

SIA is a comprehensive school management system designed for educational institutions in Colombia. It manages everything from basic institutional data to complex academic structures like grades, groups, students, teachers, enrollments, and qualifications.

## Project Overview

- **Backend:** Go 1.24.1 using the [Echo](https://echo.labstack.com/) framework.
- **Frontend:** Vanilla HTML5, CSS3, and JavaScript (located in `/web`).
- **Database:** PostgreSQL (using `lib/pq`).
- **Authentication:** JWT (JSON Web Tokens) with passwords hashed using Argon2id.
- **API:** RESTful API following the OpenAPI 3.0 specification (`api/openapi-definition.yaml`).

## Architecture

The project follows a layered architecture:
- `cmd/server/main.go`: Application entry point and router configuration.
- `internal/handlers/`: HTTP request handlers (controllers).
- `internal/repository/`: Data access layer (SQL queries).
- `internal/models/`: Data structures and database entities.
- `internal/database/`: Database connection and automatic migrations.
- `internal/middleware/`: Auth, logging, and CORS middleware.
- `internal/config/`: Configuration management via environment variables.
- `internal/utils/`: Security and JWT helpers.
- `web/`: Static frontend assets.

## Getting Started

### Prerequisites

- Go 1.24+
- PostgreSQL database
- `.env` file (see `.env.example`)

### Building and Running

1.  **Environment Setup:**
    ```bash
    cp .env.example .env
    # Edit .env with your database credentials and a secure JWT_SECRET
    ```

2.  **Run the Server:**
    ```bash
    go run cmd/server/main.go
    ```
    The server starts by default on port `8080`.

3.  **Database Migrations:**
    Migrations run automatically when the server starts. See `internal/database/migrations.go` for the schema definition.

4.  **Testing:**
    The project currently has an empty `test` directory. Standard Go testing can be performed using:
    ```bash
    go test ./...
    ```

### Accessing the Application

- **Web UI:** [http://localhost:8080/](http://localhost:8080/) (Redirects to `/web/login.html`)
- **API Health Check:** `GET /health`
- **API Base Path:** `/api/v1`

## Development Conventions

- **Handlers:** Should be initialized with their respective repositories in `main.go`.
- **Repositories:** Use standard SQL with the `database/sql` package.
- **Frontend:** Avoid heavy frameworks; use vanilla JS and CSS. Components are split into logical files in `web/js/` and `web/css/`.
- **API Documentation:** Keep `api/openapi-definition.yaml` updated with any endpoint changes.
- **Security:** Always use `middleware.JWTAuth` for protected routes and `middleware.RequireRole` for admin-only operations.

## Key Files

- `cmd/server/main.go`: Main routing logic.
- `internal/database/migrations.go`: Source of truth for the DB schema.
- `api/openapi-definition.yaml`: Comprehensive API documentation.
- `.env.example`: Template for environment configuration.
- `web/index.html`: Main entry point for the frontend dashboard.
