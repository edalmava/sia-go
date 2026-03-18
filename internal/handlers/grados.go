package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/edalmava/sia/internal/models"
	"github.com/edalmava/sia/internal/repository"
	"github.com/labstack/echo/v4"
)

type GradoHandler struct {
	repo *repository.GradoRepository
}

func NewGradoHandler(repo *repository.GradoRepository) *GradoHandler {
	return &GradoHandler{repo: repo}
}

func (h *GradoHandler) GetAll(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	nombre := c.QueryParam("nombre")

	if limit == 0 {
		limit = 20
	}

	grados, total, err := h.repo.GetAll(offset, limit, nombre)
	if err != nil {
		c.Logger().Errorf("Error fetching grados: %v", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al obtener los grados",
		})
	}

	return c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       grados,
		Pagination: models.Pagination{Total: total, Offset: offset, Limit: limit},
	})
}

func (h *GradoHandler) GetByID(c echo.Context) error {
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

	grado, err := h.repo.GetByID(id)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "not_found",
			Message: "Grado no encontrado",
		})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al obtener el grado",
		})
	}

	return c.JSON(http.StatusOK, grado)
}

func (h *GradoHandler) Create(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}
	var grado models.Grado
	if err := c.Bind(&grado); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Error en los datos de entrada",
		})
	}

	if err := h.repo.Create(&grado); err != nil {
		c.Logger().Errorf("Error creating grado: %v", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al crear el grado",
		})
	}

	return c.JSON(http.StatusCreated, grado)
}

func (h *GradoHandler) Update(c echo.Context) error {
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

	var grado models.Grado
	if err := c.Bind(&grado); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Error en los datos de entrada",
		})
	}

	grado.IDGrado = id

	if err := h.repo.Update(&grado); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al actualizar el grado",
		})
	}

	return c.JSON(http.StatusOK, grado)
}

func (h *GradoHandler) Delete(c echo.Context) error {
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
			Message: "Error al eliminar el grado",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *GradoHandler) GetAsignaturas(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	_ = id
	return c.JSON(http.StatusOK, []models.Asignatura{})
}

func (h *GradoHandler) AddAsignatura(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	_ = id
	var ga models.GradoAsignatura
	c.Bind(&ga)
	return c.JSON(http.StatusCreated, ga)
}

func (h *GradoHandler) RemoveAsignatura(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}
