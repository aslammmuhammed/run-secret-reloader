package server

import (
	"github.com/gin-gonic/gin"
)

// New returns a configured gin engine
func New() *gin.Engine {
	// Disable gin's default logger
	gin.SetMode(gin.ReleaseMode)

	// Create a new engine without any middleware
	engine := gin.New()

	// Add recovery middleware
	engine.Use(gin.Recovery())

	return engine
}
