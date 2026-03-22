package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/edalmava/sia/internal/models"
	"github.com/edalmava/sia/internal/repository"
	"github.com/labstack/echo/v4"
)

type SedeHandler struct {
	repo *repository.SedeRepository
}

func NewSedeHandler(repo *repository.SedeRepository) *SedeHandler {
	return &SedeHandler{repo: repo}
}

func (h *SedeHandler) GetAll(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}

	offset, limit := parsePagination(c)
	idInstitucion, _ := strconv.Atoi(c.QueryParam("id_institucion"))
	nombre := c.QueryParam("nombre")

	sedes, total, err := h.repo.GetAll(offset, limit, idInstitucion, nombre)
	if err != nil {
		c.Logger().Errorf("Error fetching sedes: %v", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al obtener las sedes",
		})
	}

	return c.JSON(http.StatusOK, models.PaginatedResponse{
		Data: sedes,
		Pagination: models.Pagination{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	})
}

func (h *SedeHandler) GetByID(c echo.Context) error {
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

	sede, err := h.repo.GetByID(id)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "not_found",
			Message: "Sede no encontrada",
		})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al obtener la sede",
		})
	}

	return c.JSON(http.StatusOK, sede)
}

func (h *SedeHandler) Create(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}

	var sede models.Sede
	if err := c.Bind(&sede); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Error en los datos de entrada",
		})
	}

	if err := h.repo.Create(&sede); err != nil {
		c.Logger().Errorf("Error creating sede: %v", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al crear la sede",
		})
	}

	return c.JSON(http.StatusCreated, sede)
}

func (h *SedeHandler) Update(c echo.Context) error {
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

	var sede models.Sede
	if err := c.Bind(&sede); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Error en los datos de entrada",
		})
	}

	sede.IDSede = id

	if err := h.repo.Update(&sede); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al actualizar la sede",
		})
	}

	return c.JSON(http.StatusOK, sede)
}

func (h *SedeHandler) Delete(c echo.Context) error {
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
			Message: "Error al eliminar la sede",
		})
	}

	return c.NoContent(http.StatusNoContent)
}
