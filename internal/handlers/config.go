package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/edalmava/sia/internal/middleware"
	"github.com/edalmava/sia/internal/models"
	"github.com/edalmava/sia/internal/repository"
	"github.com/labstack/echo/v4"
)

type ConfigHandler struct {
	rolRepo     *repository.RolRepository
	permisoRepo *repository.PermisoRepository
	moduloRepo  *repository.ModuloRepository
}

func NewConfigHandler(rolRepo *repository.RolRepository, permisoRepo *repository.PermisoRepository, moduloRepo *repository.ModuloRepository) *ConfigHandler {
	return &ConfigHandler{
		rolRepo:     rolRepo,
		permisoRepo: permisoRepo,
		moduloRepo:  moduloRepo,
	}
}

func (h *ConfigHandler) GetRoles(c echo.Context) error {
	if h.rolRepo == nil {
		return dbUnavailable(c)
	}

	roles, err := h.rolRepo.GetAll()
	if err != nil {
		c.Logger().Errorf("Error fetching roles: %v", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al obtener los roles",
		})
	}

	return c.JSON(http.StatusOK, models.PaginatedResponse{
		Data: roles,
	})
}

func (h *ConfigHandler) GetRoleByID(c echo.Context) error {
	if h.rolRepo == nil {
		return dbUnavailable(c)
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "ID inválido",
		})
	}

	rol, err := h.rolRepo.GetByIDWithPermisos(id)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "not_found",
			Message: "Rol no encontrado",
		})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al obtener el rol",
		})
	}

	return c.JSON(http.StatusOK, rol)
}

func (h *ConfigHandler) CreateRole(c echo.Context) error {
	if h.rolRepo == nil {
		return dbUnavailable(c)
	}

	var req struct {
		Nombre       string `json:"nombre" validate:"required,max=50"`
		Descripcion  string `json:"descripcion" validate:"max=255"`
		EsRolSistema bool   `json:"es_rol_sistema"`
		Permisos     []int  `json:"permisos"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Error en los datos de entrada",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
	}

	rol := &models.Rol{
		Nombre:       req.Nombre,
		Descripcion:  req.Descripcion,
		EsRolSistema: req.EsRolSistema,
	}

	if err := h.rolRepo.Create(rol); err != nil {
		c.Logger().Errorf("Error creating rol: %v", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al crear el rol",
		})
	}

	if len(req.Permisos) > 0 {
		h.rolRepo.SetPermisos(rol.IDRol, req.Permisos)
	}

	if claims := middleware.GetClaims(c); claims != nil {
		c.Logger().Infof("AUDIT: Usuario '%s' (ID: %d) creó rol '%s' (ID: %d)",
			claims.NombreUsuario, claims.IDUsuario, rol.Nombre, rol.IDRol)
	}

	return c.JSON(http.StatusCreated, rol)
}

func (h *ConfigHandler) UpdateRole(c echo.Context) error {
	if h.rolRepo == nil {
		return dbUnavailable(c)
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "ID inválido",
		})
	}

	existingRol, err := h.rolRepo.GetByID(id)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "not_found",
			Message: "Rol no encontrado",
		})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al obtener el rol",
		})
	}

	if existingRol.EsRolSistema {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "No se puede modificar un rol de sistema",
		})
	}

	var req struct {
		Nombre       string `json:"nombre" validate:"required,max=50"`
		Descripcion  string `json:"descripcion" validate:"max=255"`
		EsRolSistema bool   `json:"es_rol_sistema"`
		Permisos     []int  `json:"permisos"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Error en los datos de entrada",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
	}

	rol := &models.Rol{
		IDRol:        id,
		Nombre:       req.Nombre,
		Descripcion:  req.Descripcion,
		EsRolSistema: req.EsRolSistema,
	}

	if err := h.rolRepo.Update(rol); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al actualizar el rol",
		})
	}

	if req.Permisos != nil {
		h.rolRepo.SetPermisos(id, req.Permisos)
	}

	if claims := middleware.GetClaims(c); claims != nil {
		c.Logger().Infof("AUDIT: Usuario '%s' (ID: %d) actualizó rol ID: %d",
			claims.NombreUsuario, claims.IDUsuario, id)
	}

	return c.JSON(http.StatusOK, rol)
}

func (h *ConfigHandler) DeleteRole(c echo.Context) error {
	if h.rolRepo == nil {
		return dbUnavailable(c)
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "ID inválido",
		})
	}

	existingRol, err := h.rolRepo.GetByID(id)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "not_found",
			Message: "Rol no encontrado",
		})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al obtener el rol",
		})
	}

	if existingRol.EsRolSistema {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "No se puede eliminar un rol de sistema",
		})
	}

	if err := h.rolRepo.Delete(id); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al eliminar el rol",
		})
	}

	if claims := middleware.GetClaims(c); claims != nil {
		c.Logger().Infof("AUDIT: Usuario '%s' (ID: %d) eliminó rol ID: %d",
			claims.NombreUsuario, claims.IDUsuario, id)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *ConfigHandler) GetPermisos(c echo.Context) error {
	if h.permisoRepo == nil {
		return dbUnavailable(c)
	}

	permisos, err := h.permisoRepo.GetAll()
	if err != nil {
		c.Logger().Errorf("Error fetching permisos: %v", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al obtener los permisos",
		})
	}

	return c.JSON(http.StatusOK, models.PaginatedResponse{
		Data: permisos,
	})
}

func (h *ConfigHandler) GetModulos(c echo.Context) error {
	if h.moduloRepo == nil {
		return dbUnavailable(c)
	}

	modulos, err := h.moduloRepo.GetAll()
	if err != nil {
		c.Logger().Errorf("Error fetching modulos: %v", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al obtener los módulos",
		})
	}

	return c.JSON(http.StatusOK, models.PaginatedResponse{
		Data: modulos,
	})
}
