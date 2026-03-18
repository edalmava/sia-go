package handlers

import (
	"net/http"

	"github.com/edalmava/sia/internal/models"
	"github.com/labstack/echo/v4"
)

func dbUnavailable(c echo.Context) error {
	return c.JSON(http.StatusServiceUnavailable, models.ErrorResponse{
		Error:   "database_unavailable",
		Message: "Base de datos no disponible",
	})
}
