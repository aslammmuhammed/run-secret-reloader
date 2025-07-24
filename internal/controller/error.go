package controller

import (
	"net/http"

	appError "github.com/aslammmuhammed/run-secret-reloader/internal/errors"
	"github.com/aslammmuhammed/run-secret-reloader/pkg/logger"

	"github.com/gin-gonic/gin"
)

// HandleError is an utility function to handle AppError and respond appropriately
func HandleError(c *gin.Context, logger logger.Logger, err error) {
	logger.Error(c, err.Error())
	if appErr, ok := err.(*appError.AppError); ok {
		c.JSON(appErr.HTTPCode, gin.H{
			"error":   appErr.Code,
			"message": appErr.Message,
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}
