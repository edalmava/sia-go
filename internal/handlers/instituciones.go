package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/edalmava/sia/internal/models"
	"github.com/edalmava/sia/internal/repository"
	"github.com/labstack/echo/v4"
)

type InstitucionHandler struct {
	repo *repository.InstitucionRepository
}

func NewInstitucionHandler(repo *repository.InstitucionRepository) *InstitucionHandler {
	return &InstitucionHandler{repo: repo}
}

func (h *InstitucionHandler) GetAll(c echo.Context) error {
	if h.repo == nil {
		return c.JSON(http.StatusServiceUnavailable, models.ErrorResponse{
			Error:   "database_unavailable",
			Message: "Base de datos no disponible",
		})
	}

	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	nombre := c.QueryParam("nombre")

	if limit == 0 {
		limit = 20
	}

	instituciones, total, err := h.repo.GetAll(offset, limit, nombre)
	if err != nil {
		c.Logger().Errorf("Error fetching instituciones: %v", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al obtener las instituciones",
		})
	}

	return c.JSON(http.StatusOK, models.PaginatedResponse{
		Data: instituciones,
		Pagination: models.Pagination{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	})
}

func (h *InstitucionHandler) GetByID(c echo.Context) error {
	if h.repo == nil {
		return c.JSON(http.StatusServiceUnavailable, models.ErrorResponse{
			Error:   "database_unavailable",
			Message: "Base de datos no disponible",
		})
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "ID inválido",
		})
	}

	inst, err := h.repo.GetByID(id)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "not_found",
			Message: "Institución no encontrada",
		})
	}
	if err != nil {
		c.Logger().Errorf("Error fetching institucion: %v", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al obtener la institución",
		})
	}

	return c.JSON(http.StatusOK, inst)
}

func (h *InstitucionHandler) Create(c echo.Context) error {
	var institucion models.Institucion
	if err := c.Bind(&institucion); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Error en los datos de entrada",
		})
	}

	if err := h.repo.Create(&institucion); err != nil {
		c.Logger().Errorf("Error creating institucion: %v", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al crear la institución",
		})
	}

	return c.JSON(http.StatusCreated, institucion)
}

func (h *InstitucionHandler) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "ID inválido",
		})
	}

	var institucion models.Institucion
	if err := c.Bind(&institucion); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Error en los datos de entrada",
		})
	}

	institucion.IDInstitucion = id

	if err := h.repo.Update(&institucion); err != nil {
		c.Logger().Errorf("Error updating institucion: %v", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al actualizar la institución",
		})
	}

	return c.JSON(http.StatusOK, institucion)
}

func (h *InstitucionHandler) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "ID inválido",
		})
	}

	if err := h.repo.Delete(id); err != nil {
		c.Logger().Errorf("Error deleting institucion: %v", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al eliminar la institución",
		})
	}

	return c.NoContent(http.StatusNoContent)
}
