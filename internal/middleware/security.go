package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

// SecurityHeaders añade cabeceras de seguridad estándar a todas las respuestas
func SecurityHeaders() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			res := c.Response().Header()

			// Previene que el navegador intente adivinar el tipo de contenido (MIME sniffing)
			res.Set("X-Content-Type-Options", "nosniff")

			// Previene ataques de Clickjacking
			res.Set("X-Frame-Options", "DENY")

			// Habilita el filtro XSS del navegador
			res.Set("X-XSS-Protection", "1; mode=block")

			// Controla cuánta información de referencia se envía
			res.Set("Referrer-Policy", "strict-origin-when-cross-origin")

			// Política de transporte estricta (solo si es HTTPS)
			// res.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

			return next(c)
		}
	}
}

// RateLimiter crea un limitador de peticiones para proteger endpoints sensibles
// rate.Limit(2) significa 2 peticiones por segundo, con una ráfaga (burst) de 5
func RateLimiter() echo.MiddlewareFunc {
	config := echoMiddleware.RateLimiterConfig{
		Skipper: echoMiddleware.DefaultSkipper,
		Store: echoMiddleware.NewRateLimiterMemoryStoreWithConfig(
			echoMiddleware.RateLimiterMemoryStoreConfig{
				Rate:      rate.Limit(2),
				Burst:     5,
				ExpiresIn: 0,
			},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusTooManyRequests, map[string]string{
				"message": "Demasiadas peticiones. Por favor, intente más tarde.",
			})
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, map[string]string{
				"message": "Límite de peticiones excedido.",
			})
		},
	}
	return echoMiddleware.RateLimiterWithConfig(config)
}
