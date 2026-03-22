package utils

import (
	"time"

	"github.com/edalmava/sia/internal/config"
	"github.com/edalmava/sia/internal/middleware"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(cfg *config.Config, idUsuario int, nombreUsuario string, idRol int, rol string, permisos []string) (string, error) {
	claims := &middleware.Claims{
		IDUsuario:     idUsuario,
		NombreUsuario: nombreUsuario,
		IDRol:         idRol,
		Rol:           rol,
		Permisos:      permisos,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.JWT.AccessTTLMinutes) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "sia",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWT.Secret))
}

func ValidateToken(cfg *config.Config, tokenString string) (*middleware.Claims, error) {
	claims := &middleware.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}
