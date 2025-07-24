package controller

import (
	"net/http"

	"github.com/aslammmuhammed/run-secret-reloader/config"
	"github.com/aslammmuhammed/run-secret-reloader/internal/dto"
	appError "github.com/aslammmuhammed/run-secret-reloader/internal/errors"
	"github.com/aslammmuhammed/run-secret-reloader/internal/utils"
	"github.com/aslammmuhammed/run-secret-reloader/pkg/logger"
	"github.com/gin-gonic/gin"
)

// HashController provides an endpoint to get the hash for a given secret name.
// This is useful for users/tools to determine the label/annotation hash for a secret.
type HashController struct {
	logger logger.Logger
	config *config.Config
}

// NewHashController registers the /v1/hash endpoint for hashing secret names.
// Usage: POST /v1/hash {"secretName": "..."} -> {"secretNameHash": "..."}
func NewHashController(router *gin.RouterGroup, logger logger.Logger, config *config.Config) *HashController {
	controller := &HashController{
		logger: logger,
		config: config,
	}

	router.POST("", controller.GetHash)
	return controller
}

// GetHash handles POST requests to /v1/hash and returns the hash for the provided secret name.
// Request:  {"secretName": "..."}
// Response: {"secretNameHash": "..."}
func (h *HashController) GetHash(c *gin.Context) {
	var hashRequest dto.HashRequest
	if err := c.ShouldBindJSON(&hashRequest); err != nil {
		h.logger.Error(c.Request.Context(), "Failed to bind request body error: "+err.Error())
		HandleError(c, h.logger, appError.NewBadRequestError("invalid request body"))
		return
	}
	hash := utils.HashSecretName(hashRequest.SecretName)
	hashResponse := dto.HashResponse{
		SecretNameHash: hash,
	}
	c.JSON(http.StatusOK, hashResponse)
}
