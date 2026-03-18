package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/edalmava/sia/internal/models"
	"github.com/edalmava/sia/internal/repository"
	"github.com/labstack/echo/v4"
)

type JornadaHandler struct {
	repo *repository.JornadaRepository
}

func NewJornadaHandler(repo *repository.JornadaRepository) *JornadaHandler {
	return &JornadaHandler{repo: repo}
}

func (h *JornadaHandler) GetAll(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	nombre := c.QueryParam("nombre")

	if limit == 0 {
		limit = 20
	}

	jornadas, total, err := h.repo.GetAll(offset, limit, nombre)
	if err != nil {
		c.Logger().Errorf("Error fetching jornadas: %v", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al obtener las jornadas",
		})
	}

	return c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       jornadas,
		Pagination: models.Pagination{Total: total, Offset: offset, Limit: limit},
	})
}

func (h *JornadaHandler) GetByID(c echo.Context) error {
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

	jornada, err := h.repo.GetByID(id)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "not_found",
			Message: "Jornada no encontrada",
		})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al obtener la jornada",
		})
	}

	return c.JSON(http.StatusOK, jornada)
}

func (h *JornadaHandler) Create(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}
	var jornada models.Jornada
	if err := c.Bind(&jornada); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Error en los datos de entrada",
		})
	}

	if err := h.repo.Create(&jornada); err != nil {
		c.Logger().Errorf("Error creating jornada: %v", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al crear la jornada",
		})
	}

	return c.JSON(http.StatusCreated, jornada)
}

func (h *JornadaHandler) Update(c echo.Context) error {
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

	var jornada models.Jornada
	if err := c.Bind(&jornada); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Error en los datos de entrada",
		})
	}

	jornada.IDJornada = id

	if err := h.repo.Update(&jornada); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al actualizar la jornada",
		})
	}

	return c.JSON(http.StatusOK, jornada)
}

func (h *JornadaHandler) Delete(c echo.Context) error {
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
			Message: "Error al eliminar la jornada",
		})
	}

	return c.NoContent(http.StatusNoContent)
}
