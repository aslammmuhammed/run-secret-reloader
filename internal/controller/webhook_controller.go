package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/aslammmuhammed/run-secret-reloader/config"
	"github.com/aslammmuhammed/run-secret-reloader/internal/dto"
	appError "github.com/aslammmuhammed/run-secret-reloader/internal/errors"
	"github.com/aslammmuhammed/run-secret-reloader/internal/usecase"
	"github.com/aslammmuhammed/run-secret-reloader/internal/utils"
	"github.com/aslammmuhammed/run-secret-reloader/pkg/logger"
	"github.com/gin-gonic/gin"
)

type webhookController struct {
	webhookUsecase *usecase.WebhookUsecase
	logger         logger.Logger
	config         *config.Config
}

func NewWebhookController(router *gin.RouterGroup, logger logger.Logger, cfg *config.Config, webhookUsecase *usecase.WebhookUsecase) *webhookController {
	controller := &webhookController{
		webhookUsecase: webhookUsecase,
		logger:         logger,
		config:         cfg,
	}

	// Register routes
	router.POST("", controller.processWebhook)

	return controller
}

func (w *webhookController) processWebhook(c *gin.Context) {

	w.logger.Info(c.Request.Context(), "HTTP request received on endpoint: /webhook | method: POST")

	// Log raw request body
	// bodyBytes, _ := c.GetRawData()
	// w.logger.Info(c.Request.Context(), "Raw request body: "+string(bodyBytes))

	// Parse Pub/Sub request body
	var pubsubReq dto.PubSubRequest
	if err := c.ShouldBindJSON(&pubsubReq); err != nil {
		w.logger.Error(c.Request.Context(), "Failed to bind request body error: "+err.Error())
		HandleError(c, w.logger, appError.NewBadRequestError("invalid request body"))
		return
	}

	// Extract secret information
	secretAttrs := pubsubReq.Message.Attributes
	eventDescription := utils.GetEventDescription(secretAttrs.EventType)
	w.logger.Info(c.Request.Context(), "Secret event received | eventType: "+secretAttrs.EventType+" | description: "+eventDescription+" | secretId: "+secretAttrs.SecretID)
	// Check if this is an event type we need to process
	if !utils.ShouldProcessSecretEvent(secretAttrs.EventType) {
		w.logger.Info(c.Request.Context(), "Event type ignored | eventType: "+secretAttrs.EventType+" | description: "+eventDescription+" | reason: does not require service reload")

		response := dto.WebhookResponse{
			Message:     "Event processed successfully (no action required)",
			TraceID:     logger.GetTraceIDFromContext(c.Request.Context()),
			ProcessedAt: time.Now().UTC().Format(time.RFC3339),
		}
		c.JSON(http.StatusOK, response)
		return
	}

	w.logger.Info(c.Request.Context(), "Processing event that requires service reload | eventType: "+secretAttrs.EventType+" | description: "+eventDescription)

	updatedNames, failedNames, skippedNames, err := w.webhookUsecase.ProcessSecretEvent(c.Request.Context(), secretAttrs)
	if err != nil {
		HandleError(c, w.logger, err)
		return
	}

	// Log HTTP response
	allNames := append(append(updatedNames, failedNames...), skippedNames...)
	statusCode := http.StatusOK
	// If all services failed to update
	if len(allNames) > 0 && len(allNames) == len(failedNames) {
		statusCode = http.StatusInternalServerError
	}
	w.logger.Info(c.Request.Context(), "HTTP request completed | endpoint: /webhook | status: "+strconv.Itoa(statusCode)+" | servicesFound: "+strconv.Itoa(len(allNames))+" | servicesUpdated: "+strconv.Itoa(len(updatedNames))+" | servicesFailed: "+strconv.Itoa(len(failedNames))+" | servicesSkipped: "+strconv.Itoa(len(skippedNames)))

	message := fmt.Sprintf("sent response: Webhook processing completed. Updated: %d ,skipped: %d, failed: %d", len(updatedNames), len(skippedNames), len(failedNames))
	w.logger.Info(c.Request.Context(), message)

	// Return structured response
	response := dto.WebhookResponse{
		Message:        message,
		TraceID:        logger.GetTraceIDFromContext(c.Request.Context()),
		ProcessedAt:    time.Now().UTC().Format(time.RFC3339),
		ServicesFound:  len(allNames),
		ServiceNames:   allNames,
		FailedUpdates:  failedNames,
		SkippedUpdates: skippedNames,
	}

	c.JSON(statusCode, response)
}
