package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/edalmava/sia/internal/config"
	"github.com/edalmava/sia/internal/middleware"
	"github.com/edalmava/sia/internal/models"
	"github.com/edalmava/sia/internal/repository"
	"github.com/edalmava/sia/internal/utils"
	"github.com/labstack/echo/v4"
)

type AsignaturaHandler struct {
	repo *repository.AsignaturaRepository
}

func NewAsignaturaHandler(repo *repository.AsignaturaRepository) *AsignaturaHandler {
	return &AsignaturaHandler{repo: repo}
}

func (h *AsignaturaHandler) GetAll(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	if limit == 0 {
		limit = 20
	}

	asignaturas, total, err := h.repo.GetAll(offset, limit)
	if err != nil {
		c.Logger().Errorf("Error fetching asignaturas: %v", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al obtener las asignaturas",
		})
	}

	return c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       asignaturas,
		Pagination: models.Pagination{Total: total, Offset: offset, Limit: limit},
	})
}

func (h *AsignaturaHandler) GetByID(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "ID inválido",
		})
	}

	asignatura, err := h.repo.GetByID(id)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "not_found",
			Message: "Asignatura no encontrada",
		})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al obtener la asignatura",
		})
	}

	return c.JSON(http.StatusOK, asignatura)
}

func (h *AsignaturaHandler) Create(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}
	var a models.Asignatura
	if err := c.Bind(&a); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Error en los datos de entrada",
		})
	}

	if err := h.repo.Create(&a); err != nil {
		c.Logger().Errorf("Error creating asignatura: %v", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al crear la asignatura",
		})
	}

	return c.JSON(http.StatusCreated, a)
}

func (h *AsignaturaHandler) Update(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "ID inválido",
		})
	}

	var a models.Asignatura
	if err := c.Bind(&a); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Error en los datos de entrada",
		})
	}

	a.IDAsignatura = id

	if err := h.repo.Update(&a); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al actualizar la asignatura",
		})
	}

	return c.JSON(http.StatusOK, a)
}

func (h *AsignaturaHandler) Delete(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "ID inválido",
		})
	}

	if err := h.repo.Delete(id); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al eliminar la asignatura",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

type EstudianteHandler struct {
	repo *repository.EstudianteRepository
}

func NewEstudianteHandler(repo *repository.EstudianteRepository) *EstudianteHandler {
	return &EstudianteHandler{repo: repo}
}

func (h *EstudianteHandler) GetAll(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	if limit == 0 {
		limit = 20
	}

	estudiantes, total, err := h.repo.GetAll(offset, limit)
	if err != nil {
		c.Logger().Errorf("Error fetching estudiantes: %v", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al obtener los estudiantes",
		})
	}

	return c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       estudiantes,
		Pagination: models.Pagination{Total: total, Offset: offset, Limit: limit},
	})
}

func (h *EstudianteHandler) GetByID(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "ID inválido",
		})
	}

	estudiante, err := h.repo.GetByID(id)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "not_found",
			Message: "Estudiante no encontrado",
		})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al obtener el estudiante",
		})
	}

	return c.JSON(http.StatusOK, estudiante)
}

func (h *EstudianteHandler) Create(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}
	var e models.Estudiante
	if err := c.Bind(&e); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Error en los datos de entrada",
		})
	}

	if err := h.repo.Create(&e); err != nil {
		c.Logger().Errorf("Error creating estudiante: %v", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al crear el estudiante",
		})
	}

	return c.JSON(http.StatusCreated, e)
}

func (h *EstudianteHandler) Update(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "ID inválido",
		})
	}

	var e models.Estudiante
	if err := c.Bind(&e); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Error en los datos de entrada",
		})
	}

	e.IDEstudiante = id

	if err := h.repo.Update(&e); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al actualizar el estudiante",
		})
	}

	return c.JSON(http.StatusOK, e)
}

func (h *EstudianteHandler) Delete(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "ID inválido",
		})
	}

	if err := h.repo.Delete(id); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al eliminar el estudiante",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

type DocenteHandler struct {
	repo *repository.DocenteRepository
}

func NewDocenteHandler(repo *repository.DocenteRepository) *DocenteHandler {
	return &DocenteHandler{repo: repo}
}

func (h *DocenteHandler) GetAll(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	if limit == 0 {
		limit = 20
	}

	docentes, total, err := h.repo.GetAll(offset, limit)
	if err != nil {
		c.Logger().Errorf("Error fetching docentes: %v", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al obtener los docentes",
		})
	}

	return c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       docentes,
		Pagination: models.Pagination{Total: total, Offset: offset, Limit: limit},
	})
}

func (h *DocenteHandler) GetByID(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "ID inválido",
		})
	}

	docente, err := h.repo.GetByID(id)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "not_found",
			Message: "Docente no encontrado",
		})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al obtener el docente",
		})
	}

	return c.JSON(http.StatusOK, docente)
}

func (h *DocenteHandler) Create(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}
	var d models.Docente
	if err := c.Bind(&d); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Error en los datos de entrada",
		})
	}

	if err := h.repo.Create(&d); err != nil {
		c.Logger().Errorf("Error creating docente: %v", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al crear el docente",
		})
	}

	return c.JSON(http.StatusCreated, d)
}

func (h *DocenteHandler) Update(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "ID inválido",
		})
	}

	var d models.Docente
	if err := c.Bind(&d); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Error en los datos de entrada",
		})
	}

	d.IDDocente = id

	if err := h.repo.Update(&d); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al actualizar el docente",
		})
	}

	return c.JSON(http.StatusOK, d)
}

func (h *DocenteHandler) Delete(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "ID inválido",
		})
	}

	if err := h.repo.Delete(id); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al eliminar el docente",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

type AuthHandler struct {
	cfg         *config.Config
	usuarioRepo *repository.UsuarioRepository
	permisoRepo *repository.PermisoRepository
}

func NewAuthHandler(cfg *config.Config, usuarioRepo *repository.UsuarioRepository, permisoRepo *repository.PermisoRepository) *AuthHandler {
	return &AuthHandler{cfg: cfg, usuarioRepo: usuarioRepo, permisoRepo: permisoRepo}
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req models.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Credenciales inválidas",
		})
	}

	if h.usuarioRepo == nil {
		return dbUnavailable(c)
	}

	user, err := h.usuarioRepo.GetByUsername(req.Username)
	if err != nil {
		c.Logger().Errorf("Error fetching user: %v", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al iniciar sesión",
		})
	}

	if user == nil {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "auth_error",
			Message: "Credenciales inválidas",
		})
	}

	if !utils.VerifyPassword(req.Password, user.Clave) {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "auth_error",
			Message: "Credenciales inválidas",
		})
	}

	roleName := "USER"
	if user.IDRol == 1 {
		roleName = "ADMIN"
	}

	permisos, err := h.permisoRepo.GetPermissionsByUserID(user.IDUsuario)
	if err != nil {
		c.Logger().Errorf("Error fetching permisos: %v", err)
		permisos = []string{}
	}

	token, err := utils.GenerateToken(h.cfg, user.IDUsuario, user.NombreUsuario, user.IDRol, roleName, permisos)
	if err != nil {
		c.Logger().Errorf("Error generating token: %v", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al generar token",
		})
	}

	return c.JSON(http.StatusOK, models.LoginResponse{
		Token: token,
		Type:  "Bearer",
		Role:  roleName,
		IDRol: user.IDRol,
	})
}

func (h *AuthHandler) Refresh(c echo.Context) error {
	claims := middleware.GetClaims(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "auth_error",
			Message: "Token inválido",
		})
	}

	roleName := "USER"
	if claims.IDRol == 1 {
		roleName = "ADMIN"
	}

	permisos, err := h.permisoRepo.GetPermissionsByUserID(claims.IDUsuario)
	if err != nil {
		c.Logger().Errorf("Error fetching permisos: %v", err)
		permisos = claims.Permisos
	}

	token, err := utils.GenerateToken(h.cfg, claims.IDUsuario, claims.NombreUsuario, claims.IDRol, roleName, permisos)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al generar token",
		})
	}

	return c.JSON(http.StatusOK, models.LoginResponse{
		Token: token,
		Type:  "Bearer",
		Role:  roleName,
		IDRol: claims.IDRol,
	})
}

func (h *AuthHandler) Logout(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "Logout exitoso"})
}
