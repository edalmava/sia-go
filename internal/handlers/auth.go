package handlers

import (
	"net/http"
	"time"

	"github.com/edalmava/sia/internal/config"
	"github.com/edalmava/sia/internal/middleware"
	"github.com/edalmava/sia/internal/models"
	"github.com/edalmava/sia/internal/repository"
	"github.com/edalmava/sia/internal/utils"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	cfg              *config.Config
	usuarioRepo      *repository.UsuarioRepository
	permisoRepo      *repository.PermisoRepository
	rolRepo          *repository.RolRepository
	refreshTokenRepo *repository.RefreshTokenRepository
	revokedTokenRepo *repository.RevokedTokenRepository
}

func NewAuthHandler(cfg *config.Config, usuarioRepo *repository.UsuarioRepository, permisoRepo *repository.PermisoRepository, rolRepo *repository.RolRepository, refreshTokenRepo *repository.RefreshTokenRepository, revokedTokenRepo *repository.RevokedTokenRepository) *AuthHandler {
	return &AuthHandler{cfg: cfg, usuarioRepo: usuarioRepo, permisoRepo: permisoRepo, rolRepo: rolRepo, refreshTokenRepo: refreshTokenRepo, revokedTokenRepo: revokedTokenRepo}
}

func (h *AuthHandler) setRefreshTokenCookie(c echo.Context, token string, expiration time.Time) {
	cookie := new(http.Cookie)
	cookie.Name = "refresh_token"
	cookie.Value = token
	cookie.Expires = expiration
	cookie.Path = "/auth" // Solo se envía a los endpoints de autenticación
	cookie.HttpOnly = true
	cookie.Secure = h.cfg.Env == "production" // Solo seguro en producción (HTTPS)
	cookie.SameSite = http.SameSiteStrictMode
	c.SetCookie(cookie)
}

func (h *AuthHandler) setAccessTokenCookie(c echo.Context, token string, expirationMinutes int) {
	cookie := new(http.Cookie)
	cookie.Name = "auth_token"
	cookie.Value = token
	cookie.Expires = time.Now().Add(time.Duration(expirationMinutes) * time.Minute)
	cookie.Path = "/" // Se envía a todas las rutas de la API
	cookie.HttpOnly = true
	cookie.Secure = h.cfg.Env == "production"
	cookie.SameSite = http.SameSiteStrictMode
	c.SetCookie(cookie)
}

func (h *AuthHandler) clearRefreshTokenCookie(c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = "refresh_token"
	cookie.Value = ""
	cookie.Expires = time.Now().Add(-1 * time.Hour)
	cookie.Path = "/auth"
	cookie.HttpOnly = true
	cookie.Secure = h.cfg.Env == "production"
	cookie.SameSite = http.SameSiteStrictMode
	c.SetCookie(cookie)
}

func (h *AuthHandler) clearAccessTokenCookie(c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = "auth_token"
	cookie.Value = ""
	cookie.Expires = time.Now().Add(-1 * time.Hour)
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = h.cfg.Env == "production"
	cookie.SameSite = http.SameSiteStrictMode
	c.SetCookie(cookie)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req models.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Credenciales inválidas",
		})
	}

	if h.usuarioRepo == nil {
		return dbUnavailable(c)
	}

	user, err := h.usuarioRepo.GetByUsername(req.Username)
	if err != nil {
		c.Logger().Errorf("Error fetching user: %v", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al iniciar sesión",
		})
	}

	if user == nil {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "auth_error",
			Message: "Credenciales inválidas",
		})
	}

	if !utils.VerifyPassword(req.Password, user.Clave) {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "auth_error",
			Message: "Credenciales inválidas",
		})
	}

	roleName := "USER"
	if h.rolRepo != nil {
		rol, err := h.rolRepo.GetByID(user.IDRol)
		if err == nil && rol != nil {
			roleName = rol.Nombre
		}
	}

	permisos, err := h.permisoRepo.GetPermissionsByUserID(user.IDUsuario)
	if err != nil {
		c.Logger().Errorf("Error fetching permisos: %v", err)
		permisos = []string{}
	}

	accessToken, accessJTI, expiresIn, err := utils.GenerateAccessToken(h.cfg, user.IDUsuario, user.NombreUsuario, user.IDRol, roleName, permisos)
	if err != nil {
		c.Logger().Errorf("Error generating access token: %v", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al generar token",
		})
	}

	refreshToken, tokenHash, expiration, err := utils.GenerateRefreshToken()
	if err != nil {
		c.Logger().Errorf("Error generating refresh token: %v", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al generar refresh token",
		})
	}

	if err := h.refreshTokenRepo.Create(tokenHash, accessJTI, user.IDUsuario, expiration, ""); err != nil {
		c.Logger().Errorf("Error saving refresh token: %v", err)
	}

	// Establecer cookies HttpOnly
	h.setRefreshTokenCookie(c, refreshToken, expiration)
	h.setAccessTokenCookie(c, accessToken, h.cfg.JWT.AccessTTLMinutes)

	return c.JSON(http.StatusOK, models.LoginResponse{
		AccessToken:   "", // Ya no se envía en el cuerpo por seguridad
		TokenType:     "Bearer",
		ExpiresIn:     expiresIn,
		NombreUsuario: user.NombreUsuario,
		Role:          roleName,
		IDRol:         user.IDRol,
		Permisos:      permisos,
	})
}

func (h *AuthHandler) Refresh(c echo.Context) error {
	if h.refreshTokenRepo == nil {
		return dbUnavailable(c)
	}

	// Obtener el refresh token desde la cookie
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "auth_error",
			Message: "Refresh token no encontrado",
		})
	}
	refreshTokenStr := cookie.Value

	tokenHash := utils.HashRefreshToken(refreshTokenStr)
	storedToken, err := h.refreshTokenRepo.GetByTokenHash(tokenHash)
	if err != nil {
		c.Logger().Errorf("Error fetching refresh token: %v", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al validar refresh token",
		})
	}

	if storedToken == nil || !storedToken.Activo {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "auth_error",
			Message: "Refresh token inválido o expirado",
		})
	}

	if time.Now().After(storedToken.FechaExpiracion) {
		h.refreshTokenRepo.Revoke(tokenHash)
		h.clearRefreshTokenCookie(c)
		h.clearAccessTokenCookie(c)
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "auth_error",
			Message: "Refresh token expirado",
		})
	}

	user, err := h.usuarioRepo.GetByID(storedToken.IDUsuario)
	if err != nil || user == nil {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "auth_error",
			Message: "Usuario no encontrado",
		})
	}

	permisos, err := h.permisoRepo.GetPermissionsByUserID(storedToken.IDUsuario)
	if err != nil {
		c.Logger().Errorf("Error fetching permisos: %v", err)
		permisos = []string{}
	}

	roleName := "USER"
	if h.rolRepo != nil {
		rol, err := h.rolRepo.GetByID(user.IDRol)
		if err == nil && rol != nil {
			roleName = rol.Nombre
		}
	}

	// Invalidar el token anterior (rotación de tokens)
	h.refreshTokenRepo.Revoke(tokenHash)

	accessToken, accessJTI, expiresIn, err := utils.GenerateAccessToken(h.cfg, user.IDUsuario, user.NombreUsuario, user.IDRol, roleName, permisos)
	if err != nil {
		c.Logger().Errorf("Error generating access token: %v", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al generar token",
		})
	}

	newRefreshToken, newTokenHash, expiration, err := utils.GenerateRefreshToken()
	if err != nil {
		c.Logger().Errorf("Error generating refresh token: %v", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al generar refresh token",
		})
	}

	if err := h.refreshTokenRepo.Create(newTokenHash, accessJTI, user.IDUsuario, expiration, ""); err != nil {
		c.Logger().Errorf("Error saving refresh token: %v", err)
	}

	// Establecer las nuevas cookies
	h.setRefreshTokenCookie(c, newRefreshToken, expiration)
	h.setAccessTokenCookie(c, accessToken, h.cfg.JWT.AccessTTLMinutes)

	return c.JSON(http.StatusOK, models.RefreshTokenResponse{
		AccessToken:   "", // Ya no se envía en el cuerpo
		TokenType:     "Bearer",
		ExpiresIn:     expiresIn,
		NombreUsuario: user.NombreUsuario,
		Role:          roleName,
		IDRol:         user.IDRol,
		Permisos:      permisos,
	})
}

func (h *AuthHandler) Logout(c echo.Context) error {
	if h.refreshTokenRepo == nil {
		return dbUnavailable(c)
	}

	// Revocar el access token si existe
	jti := middleware.GetJTI(c)
	if jti != "" && h.revokedTokenRepo != nil {
		expiresAt := time.Now().Add(time.Duration(h.cfg.JWT.AccessTTLMinutes) * time.Minute)
		h.revokedTokenRepo.Add(jti, expiresAt)
	}

	// Revocar el refresh token desde la cookie
	cookie, err := c.Cookie("refresh_token")
	if err == nil {
		tokenHash := utils.HashRefreshToken(cookie.Value)
		h.refreshTokenRepo.Revoke(tokenHash)
	}

	// Limpiar las cookies
	h.clearRefreshTokenCookie(c)
	h.clearAccessTokenCookie(c)

	return c.JSON(http.StatusOK, map[string]string{"message": "Sesión cerrada exitosamente"})
}

func (h *AuthHandler) LogoutAll(c echo.Context) error {
	if h.refreshTokenRepo == nil {
		return dbUnavailable(c)
	}

	claims := middleware.GetClaims(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "auth_error",
			Message: "Usuario no autenticado",
		})
	}

	jti := middleware.GetJTI(c)
	if jti != "" && h.revokedTokenRepo != nil {
		expiresAt := time.Now().Add(time.Duration(h.cfg.JWT.AccessTTLMinutes) * time.Minute)
		h.revokedTokenRepo.Add(jti, expiresAt)
	}

	h.refreshTokenRepo.RevokeAllForUser(claims.IDUsuario)
	
	// Limpiar las cookies del cliente actual
	h.clearRefreshTokenCookie(c)
	h.clearAccessTokenCookie(c)

	return c.JSON(http.StatusOK, map[string]string{"message": "Todas las sesiones cerradas exitosamente"})
}
