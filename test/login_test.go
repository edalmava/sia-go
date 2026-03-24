package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/edalmava/sia/internal/config"
	"github.com/edalmava/sia/internal/handlers"
	"github.com/edalmava/sia/internal/models"
	"github.com/edalmava/sia/internal/repository"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestLoginProduction(t *testing.T) {
	// Mock config
	cfg := &config.Config{
		Env: "production",
		JWT: config.JWTConfig{
			Secret:           "test_secret",
			AccessTTLMinutes: 15,
		},
	}

	e := echo.New()
	
	// We need a real-ish repository or a mock. 
    // Since we want to test the handler logic, a nil repo will return 503.
    // Let's just see if the handler itself has some logic that returns 401.

	reqJSON := `{"username":"admin","password":"admin"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.HeaderApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := handlers.NewAuthHandler(cfg, nil, nil, nil, nil, nil)

	// This should return 503 because repo is nil
	err := h.Login(c)
	if err != nil {
		t.Errorf("Login returned error: %v", err)
	}

	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
}
