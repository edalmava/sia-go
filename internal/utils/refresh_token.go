package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/edalmava/sia/internal/config"
	"github.com/edalmava/sia/internal/middleware"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateAccessToken(cfg *config.Config, idUsuario int, nombreUsuario string, idRol int, rol string, permisos []string) (string, string, int, error) {
	expiresInSeconds := cfg.JWT.AccessTTLMinutes * 60
	jti := uuid.New().String()

	claims := &middleware.Claims{
		IDUsuario:     idUsuario,
		NombreUsuario: nombreUsuario,
		IDRol:         idRol,
		Rol:           rol,
		Permisos:      permisos,
		JTI:           jti,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.JWT.AccessTTLMinutes) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "sia",
			ID:        jti,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JWT.Secret))
	return tokenString, jti, expiresInSeconds, err
}

func GenerateRefreshToken() (string, string, time.Time, error) {
	refreshBytes := make([]byte, 32)
	if _, err := rand.Read(refreshBytes); err != nil {
		return "", "", time.Time{}, err
	}

	refreshToken := hex.EncodeToString(refreshBytes)
	hash := HashRefreshToken(refreshToken)
	expiration := time.Now().Add(168 * time.Hour)

	return refreshToken, hash, expiration, nil
}

func HashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func ValidateRefreshToken(cfg *config.Config, tokenString string) (*middleware.Claims, error) {
	claims := &middleware.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}
