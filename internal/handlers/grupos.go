package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/edalmava/sia/internal/models"
	"github.com/edalmava/sia/internal/repository"
	"github.com/labstack/echo/v4"
)

type GrupoHandler struct {
	repo *repository.GrupoRepository
}

func NewGrupoHandler(repo *repository.GrupoRepository) *GrupoHandler {
	return &GrupoHandler{repo: repo}
}

func (h *GrupoHandler) GetAll(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	idGrado, _ := strconv.Atoi(c.QueryParam("id_grado"))
	idSede, _ := strconv.Atoi(c.QueryParam("id_sede"))
	idJornada, _ := strconv.Atoi(c.QueryParam("id_jornada"))

	if limit == 0 {
		limit = 20
	}

	grupos, total, err := h.repo.GetAll(offset, limit, idGrado, idSede, idJornada)
	if err != nil {
		c.Logger().Errorf("Error fetching grupos: %v", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al obtener los grupos",
		})
	}

	return c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       grupos,
		Pagination: models.Pagination{Total: total, Offset: offset, Limit: limit},
	})
}

func (h *GrupoHandler) GetByID(c echo.Context) error {
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

	grupo, err := h.repo.GetByID(id)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "not_found",
			Message: "Grupo no encontrado",
		})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al obtener el grupo",
		})
	}

	return c.JSON(http.StatusOK, grupo)
}

func (h *GrupoHandler) Create(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}
	var grupo models.Grupo
	if err := c.Bind(&grupo); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Error en los datos de entrada",
		})
	}

	if err := h.repo.Create(&grupo); err != nil {
		c.Logger().Errorf("Error creating grupo: %v", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al crear el grupo",
		})
	}

	return c.JSON(http.StatusCreated, grupo)
}

func (h *GrupoHandler) Update(c echo.Context) error {
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

	var grupo models.Grupo
	if err := c.Bind(&grupo); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Error en los datos de entrada",
		})
	}

	grupo.IDGrupo = id

	if err := h.repo.Update(&grupo); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al actualizar el grupo",
		})
	}

	return c.JSON(http.StatusOK, grupo)
}

func (h *GrupoHandler) Delete(c echo.Context) error {
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
			Message: "Error al eliminar el grupo",
		})
	}

	return c.NoContent(http.StatusNoContent)
}
