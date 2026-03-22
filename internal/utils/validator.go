package utils

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// CustomValidator envuelve el validador de go-playground
type CustomValidator struct {
	Validator *validator.Validate
}

// Validate implementa la interfaz echo.Validator
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		// Opcionalmente, se puede transformar el error en algo más amigable
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

// NewValidator crea una nueva instancia del validador
func NewValidator() *CustomValidator {
	return &CustomValidator{Validator: validator.New()}
}
