# SIA - AGENTS.md
## Project Overview
SIA - REST API para gestión de instituciones educativas en Colombia.
**Tech Stack:** Go 1.24.1, Echo v4, JWT, PostgreSQL, Argon2id

---

## Build, Test & Commands
```bash
go build -o sia.exe ./cmd/server    # Build binary
go run ./cmd/server                  # Run server
go test ./...                        # Run all tests
go test -v -run TestName ./path      # Run single test
go fmt ./...; go vet ./...            # Format & Lint
go mod tidy                          # Tidy dependencies
```

---

## Environment Variables
```bash
SERVER_PORT=8080; JWT_SECRET=secret; JWT_EXPIRY=24h
DB_HOST=localhost; DB_PORT=5433; DB_USER=postgres
DB_PASSWORD=postgres; DB_NAME=sia
```
**Default Login:** admin / admin123

---

## Code Style

### File Organization
```
cmd/server/main.go
internal/
├── config/       # Configuration
├── database/     # DB connection, migrations
├── handlers/     # HTTP handlers (one file per entity)
├── middleware/   # Auth, logging
├── models/       # All structs in models.go
├── repository/   # Data access layer
└── utils/        # Helpers (JWT, password hashing)
web/             # Static HTML/CSS/JS
```

### Imports (order - REQUIRED)
```go
import (
    "database/sql"; "net/http"; "strconv"; "strings"
    "github.com/edalmava/sia/internal/middleware"
    "github.com/edalmava/sia/internal/models"
    "github.com/edalmava/sia/internal/repository"
    "github.com/labstack/echo/v4"
)
```

### Naming Conventions
| Element | Convention | Example |
|---------|------------|---------|
| Packages | snake_case | `repository`, `middleware` |
| Files | snake_case.go | `usuario.go` |
| Exported types | PascalCase | `UsuarioHandler`, `Rol` |
| Variables/fields | camelCase | `nombreUsuario`, `idRol` |
| DB/JSON fields | snake_case | `id_usuario` |
| Errors | errXXX | `errNotFound`, `errInvalidInput` |

### Struct Tags
```go
type Usuario struct {
    IDUsuario     int    `json:"id_usuario" db:"id_usuario"`
    NombreUsuario string `json:"nombre_usuario" db:"nombre_usuario"`
    Clave         string `json:"-" db:"clave"` // Never return in JSON
}
```

### Error Handling
```go
func (h *Handler) GetByID(c echo.Context) error {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil { return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: "ID inválido"}) }
    entity, err := h.repo.GetByID(id)
    if err == sql.ErrNoRows { return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "not_found", Message: "Entidad no encontrada"}) }
    if err != nil { c.Logger().Errorf("Error: %v", err); return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal_error", Message: "Error al obtener la entidad"}) }
    return c.JSON(http.StatusOK, entity)
}
```

### Handler Pattern
```go
type EntityHandler struct { repo *repository.EntityRepository }
func NewEntityHandler(repo *repository.EntityRepository) *EntityHandler { return &EntityHandler{repo: repo} }
func (h *EntityHandler) GetAll(c echo.Context) error { ... }
func (h *EntityHandler) GetByID(c echo.Context) error { ... }
func (h *EntityHandler) Create(c echo.Context) error { ... }
func (h *EntityHandler) Update(c echo.Context) error { ... }
func (h *EntityHandler) Delete(c echo.Context) error { ... }
```

### Repository Pattern
```go
type EntityRepository struct { db *sql.DB }
func NewEntityRepository(db *sql.DB) *EntityRepository { return &EntityRepository{db: db} }
func (r *EntityRepository) GetAll(offset, limit int) ([]Entity, int, error) { ... }
func (r *EntityRepository) GetByID(id int) (*Entity, error) { ... }
func (r *EntityRepository) Create(e *Entity) error { ... }
func (r *EntityRepository) Update(e *Entity) error { ... }
func (r *EntityRepository) Delete(id int) error { ... }
```

### Middleware Pattern
```go
func JWTAuth(cfg *config.Config) echo.MiddlewareFunc { return func(next echo.HandlerFunc) echo.HandlerFunc { return func(c echo.Context) error { return next(c) } } }
func RequirePermission(permisos ...string) echo.MiddlewareFunc { return func(next echo.HandlerFunc) echo.HandlerFunc { return func(c echo.Context) error { return next(c) } } }
```

### Pagination
```GET /entity?offset=0&limit=20``` → Response: `{ "data": [...], "pagination": { "total": 150, "offset": 0, "limit": 20 } }`
Default limit: 20, Max: 100

---

## Auth & Authorization
- Passwords MUST use Argon2id (`utils.HashPassword`)
- JWT tokens include user ID, role, and permissions
- Use `middleware.RequirePermission("codigo_permiso")` for protection
- Permission codes: `recurso_accion` (e.g., `usuarios_crear`)

---

## Database Conventions
- Use `github.com/lib/pq` for PostgreSQL
- Always use parameterized queries (NO string concatenation)
- Use `defer rows.Close()` for query results
- Handle `sql.ErrNoRows` explicitly for GetByID

---

## Quick Reference
| Task | Command |
|------|---------|
| Build | `go build -o sia.exe ./cmd/server` |
| Run | `./sia.exe` |
| Test | `go test ./...` |
| Single test | `go test -v -run TestName ./path` |

---

## Architecture
- **Repository**: DB access via `internal/repository/`
- **Handlers**: Business logic, request/response
- **Middleware**: Auth, logging
- **Models**: Single `models.go` with all structs
- **Database**: PostgreSQL with auto-migrations
- **Auth**: JWT + Argon2id
- **Static**: Web UI in `web/`
