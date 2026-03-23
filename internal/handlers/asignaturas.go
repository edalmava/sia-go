package handlers

import (
	"database/sql"
	"net/http"

	"github.com/edalmava/sia/internal/models"
	"github.com/edalmava/sia/internal/repository"
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

	offset, limit := parsePagination(c)

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

	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}

	asignatura, err := h.repo.GetByID(id)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "not_found",
			Message: "Asignatura no encontrada",
		})
	}
	if err != nil {
		c.Logger().Errorf("Error fetching asignatura: %v", err)
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
			Message: "Datos inválidos",
		})
	}

	if err := c.Validate(&a); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
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

	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}

	var a models.Asignatura
	if err := c.Bind(&a); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Datos inválidos",
		})
	}

	if err := c.Validate(&a); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
	}

	a.IDAsignatura = id
	if err := h.repo.Update(&a); err != nil {
		c.Logger().Errorf("Error updating asignatura: %v", err)
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

	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}

	if err := h.repo.Delete(id); err != nil {
		c.Logger().Errorf("Error deleting asignatura: %v", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al eliminar la asignatura",
		})
	}

	return c.NoContent(http.StatusNoContent)
}
