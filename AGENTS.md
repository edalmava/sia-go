# SIA - AGENTS.md

REST API para gestión de instituciones educativas en Colombia.
**Tech Stack:** Go 1.24.1, Echo v4, JWT, PostgreSQL, Argon2id

---

## Build, Test & Commands

```bash
go build -o sia.exe ./cmd/server    # Build binary
go run ./cmd/server                  # Run server
go test ./...                        # Run all tests
go test -v -run TestName ./path      # Run single test (e.g., go test -v -run TestUsuario ./internal/repository)
go fmt ./... && go vet ./...          # Format & Lint
go mod tidy                          # Tidy dependencies
```

---

## Environment Variables

```bash
SERVER_PORT=8080
JWT_SECRET=<REQUIRED - panics if missing>
JWT_EXPIRY=24h
DB_HOST=localhost; DB_PORT=5433
DB_USER=postgres; DB_PASSWORD=postgres; DB_NAME=sia
```

**Default Login:** admin / admin123

---

## File Organization

```
cmd/server/main.go
internal/
├── config/       # Configuration (panics if JWT_SECRET missing)
├── database/     # DB connection, auto-migrations
├── handlers/     # HTTP handlers (one file per entity/domain)
├── middleware/   # Auth (JWT), logging, CORS
├── models/       # All structs in models.go + request/response types
├── repository/   # Data access layer (SQL queries)
└── utils/        # JWT, password hashing (Argon2id)
web/
├── index.html, login.html
├── pages/        # dashboard.html, usuarios.html, configuracion.html
├── js/           # ES6 modules (api.js, auth.js, ui.js, config.js)
└── css/          # Modular CSS (01-variables.css, 02-base.css, etc.)
```

---

## Go Code Style

### Imports (order - REQUIRED)
```go
import (
    "database/sql"; "net/http"; "strconv"
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
| Files | snake_case.go | `usuario.go`, `academicos.go` |
| Exported types | PascalCase | `UsuarioHandler`, `Rol` |
| Variables/fields | PascalCase Go, snake_case JSON/DB | `NombreUsuario` / `nombre_usuario` |
| DB/JSON tags | snake_case | `json:"nombre_usuario" db:"nombre_usuario"` |
| Errors | err prefix | `errNotFound`, `errInvalidInput` |

### Struct Tags
```go
type Usuario struct {
    IDUsuario     int    `json:"id_usuario" db:"id_usuario"`
    NombreUsuario string `json:"nombre_usuario" db:"nombre_usuario"`
    Clave         string `json:"-" db:"clave"` // SECURITY: Never return in JSON
}
```

### Request/Response Pattern (SECURITY)
Use separate request structs to avoid accidentally exposing sensitive fields:

```go
// models.go
type UsuarioResponse struct {      // Safe - no password field
    IDUsuario     int    `json:"id_usuario"`
    NombreUsuario string `json:"nombre_usuario"`
    // ... no Clave field
}

type UsuarioCreateRequest struct { // For binding CREATE requests
    NombreUsuario string `json:"nombre_usuario"`
    Clave         string `json:"clave"` // Include for creation
    IDRol         int    `json:"id_rol"`
}

type UsuarioUpdateRequest struct { // For binding UPDATE requests
    NombreUsuario string `json:"nombre_usuario"` // No Clave field
    IDRol         int    `json:"id_rol"`
}
```

### Error Handling Pattern
```go
func (h *Handler) GetByID(c echo.Context) error {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return c.JSON(http.StatusBadRequest, models.ErrorResponse{
            Error:   "validation_error",
            Message: "ID inválido",
        })
    }
    entity, err := h.repo.GetByID(id)
    if err == sql.ErrNoRows {
        return c.JSON(http.StatusNotFound, models.ErrorResponse{
            Error:   "not_found",
            Message: "Entidad no encontrada",
        })
    }
    if err != nil {
        c.Logger().Errorf("Error: %v", err)
        return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
            Error:   "internal_error",
            Message: "Error al obtener la entidad",
        })
    }
    return c.JSON(http.StatusOK, entity)
}
```

### Handler Pattern
```go
type EntityHandler struct { repo *repository.EntityRepository }
func NewEntityHandler(repo *repository.EntityRepository) *EntityHandler { return &EntityHandler{repo} }
func (h *EntityHandler) GetAll(c echo.Context) error { ... }
func (h *EntityHandler) GetByID(c echo.Context) error { ... }
func (h *EntityHandler) Create(c echo.Context) error { ... }
func (h *EntityHandler) Update(c echo.Context) error { ... }
func (h *EntityHandler) Delete(c echo.Context) error { ... }
```

### Repository Pattern
```go
type EntityRepository struct { db *sql.DB }
func NewEntityRepository(db *sql.DB) *EntityRepository { return &EntityRepository{db} }
func (r *EntityRepository) GetAll(offset, limit int) ([]Model, int, error) { ... }
func (r *EntityRepository) GetByID(id int) (*Model, error) { ... }
func (r *EntityRepository) Create(m *Model) error { ... }
func (r *EntityRepository) Update(m *Model) error { ... }
func (r *EntityRepository) Delete(id int) error { ... }
```

### Middleware Pattern
```go
func JWTAuth(cfg *config.Config) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error { return next(c) }
    }
}
func RequirePermission(permisos ...string) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error { return next(c) }
    }
}
```

### Pagination
`GET /entity?offset=0&limit=20` → Response: `{ "data": [...], "pagination": { "total": 150, "offset": 0, "limit": 20 } }`
Default limit: 20, Max: 100

---

## Frontend Code Style (ES6 Modules)

### File Structure
```
web/js/
├── config.js      # Constants (API_URL, endpoints, role names)
├── auth.js        # Auth helpers (checkAuth, hasPermission, logout)
├── api.js         # Generic API client + specific modules (usuarioApi, rolApi)
├── ui.js          # UI helpers (showToast, modals, handleApiError)
├── login.js       # Login page logic
├── dashboard.js   # Dashboard logic
├── users.js       # User management (CRUD + password change)
└── configuracion.js # Roles/permissions configuration
```

### API Pattern
```javascript
// api.js - Generic client with 401 handling
export const api = {
    async get(endpoint, params = {}) { ... },
    async post(endpoint, data) { ... },
    async put(endpoint, data) { ... },
    async delete(endpoint) { ... }
};

// Specific modules
export const usuarioApi = {
    getAll: (offset, limit, search) => api.get('/usuarios', { offset, limit, search }),
    create: (data) => api.post('/usuarios', data),
    update: (id, data) => api.put(`/usuarios/${id}`, data),
    delete: (id) => api.delete(`/usuarios/${id}`),
    changePassword: (id, newPassword) => api.post(`/usuarios/${id}/change-password`, { password: newPassword })
};
```

### Form Submission Protection
Always disable submit buttons during async operations to prevent double-submission:
```javascript
async function handleSubmit(e) {
    e.preventDefault();
    const btn = e.target.querySelector('button[type="submit"]');
    if (btn.disabled) return;
    btn.disabled = true;
    
    try {
        await api.post('/endpoint', data);
        showToast('Éxito', 'success');
    } catch (error) {
        handleApiError(error, 'Error al guardar');
    } finally {
        btn.disabled = false;
    }
}
```

### Auth Helpers (auth.js)
```javascript
export function checkAuth() { ... }        // Verify token exists
export function hasPermission(permiso) { }  // Check single permission
export function hasAnyPermissionPrefix(perms) { } // Check prefix matching
export function getUserData() { }           // Get { username, role, permissions }
export function logout() { }                // Clear token and redirect
```

---

## Database Conventions
- Use `github.com/lib/pq` for PostgreSQL
- Always use parameterized queries (NO string concatenation)
- Use `defer rows.Close()` for query results
- Handle `sql.ErrNoRows` explicitly for GetByID
- Repository methods for reads return safe types (e.g., `UsuarioResponse`) without password hashes

---

## Auth & Authorization
- Passwords MUST use Argon2id (`utils.HashPassword`)
- JWT tokens include user ID, role, and permissions array
- Use `middleware.RequirePermission("codigo_permiso")` for route protection
- Permission codes: `recurso_accion` (e.g., `usuarios_crear`, `usuarios_editar`)
- `GET /usuarios` returns `UsuarioResponse` (no password), not `Usuario`

---

## Quick Reference
| Task | Command |
|------|---------|
| Build | `go build -o sia.exe ./cmd/server` |
| Run | `./sia.exe` |
| Test | `go test ./...` |
| Single test | `go test -v -run TestName ./path` |
| Lint | `go fmt ./... && go vet ./...` |
