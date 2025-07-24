package middlewares

import (
	"fmt"
	"time"

	"github.com/aslammmuhammed/run-secret-reloader/pkg/logger"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Debug log for request details
		log.Debug(c.Request.Context(), fmt.Sprintf("Request: %s %s",
			method, path))

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)
		status := c.Writer.Status()

		// Debug log for response details
		log.Debug(c.Request.Context(), fmt.Sprintf("Response: status=%d, duration=%v",
			status, duration))

		// Standard info log with key metrics
		log.Info(c.Request.Context(), fmt.Sprintf("%s %s [%d] %v", method, path, status, duration))
	}
}
