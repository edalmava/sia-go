package middleware

import (
	"net/http"
	"strings"

	"github.com/edalmava/sia/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Claims struct {
	IDUsuario     int      `json:"id_usuario"`
	NombreUsuario string   `json:"nombre_usuario"`
	IDRol         int      `json:"id_rol"`
	Rol           string   `json:"rol"`
	Permisos      []string `json:"permisos"`
	jwt.RegisteredClaims
}

func JWTAuth(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error":   "auth_error",
					"message": "Token de autorización requerido",
				})
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error":   "auth_error",
					"message": "Formato de token inválido",
				})
			}

			tokenString := parts[1]
			claims := &Claims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(cfg.JWT.Secret), nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error":   "auth_error",
					"message": "Token inválido o expirado",
				})
			}

			c.Set("user", claims)
			return next(c)
		}
	}
}

func GetClaims(c echo.Context) *Claims {
	user := c.Get("user")
	if user != nil {
		return user.(*Claims)
	}
	return nil
}

func RequireRole(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims := GetClaims(c)
			if claims == nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error":   "auth_error",
					"message": "Usuario no autenticado",
				})
			}

			for _, role := range roles {
				if claims.Rol == role {
					return next(c)
				}
			}

			return c.JSON(http.StatusForbidden, map[string]string{
				"error":   "forbidden",
				"message": "No tienes permisos para acceder a este recurso",
			})
		}
	}
}

func RequirePermission(permisos ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims := GetClaims(c)
			if claims == nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error":   "auth_error",
					"message": "Usuario no autenticado",
				})
			}

			for _, permiso := range permisos {
				if hasPermiso(claims.Permisos, permiso) {
					return next(c)
				}
			}

			return c.JSON(http.StatusForbidden, map[string]string{
				"error":   "forbidden",
				"message": "No tienes el permiso necesario: " + strings.Join(permisos, " o "),
			})
		}
	}
}

func hasPermiso(permisosDelUsuario []string, permisoRequerido string) bool {
	for _, p := range permisosDelUsuario {
		if p == permisoRequerido {
			return true
		}
	}
	return false
}
