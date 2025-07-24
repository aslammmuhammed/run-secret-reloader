package middlewares

import (
	"github.com/aslammmuhammed/run-secret-reloader/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CorsMiddleware creates and returns a CORS middleware with default configuration
func CorsMiddleware(cfg *config.Config) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowAllOrigins:  true, // For webhook endpoints
		AllowMethods:     cfg.HTTP.AllowedMethods,
		AllowHeaders:     cfg.HTTP.AllowedHeaders,
		ExposeHeaders:    cfg.HTTP.ExposeHeaders,
		AllowCredentials: cfg.HTTP.AllowCredentials,
	})
}
