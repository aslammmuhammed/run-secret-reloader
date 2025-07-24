package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler is a struct that contains the health handler
type HealthHandler struct {
}

// Get the health status of the service
func GetHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Service is UP"})
}
