package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/aslammmuhammed/run-secret-reloader/config"
	"github.com/aslammmuhammed/run-secret-reloader/internal/constants"
	"github.com/aslammmuhammed/run-secret-reloader/internal/dto"
	appError "github.com/aslammmuhammed/run-secret-reloader/internal/errors"
	"github.com/aslammmuhammed/run-secret-reloader/internal/usecase"
	"github.com/aslammmuhammed/run-secret-reloader/internal/utils"
	"github.com/aslammmuhammed/run-secret-reloader/pkg/logger"
	"github.com/gin-gonic/gin"
)

type WebhookController struct {
	webhookUsecase *usecase.WebhookUsecase
	logger         logger.Logger
	config         *config.Config
}

func NewWebhookController(router *gin.RouterGroup, logger logger.Logger, cfg *config.Config, webhookUsecase *usecase.WebhookUsecase) *WebhookController {
	controller := &WebhookController{
		webhookUsecase: webhookUsecase,
		logger:         logger,
		config:         cfg,
	}

	// Register routes
	router.POST("", controller.ProcessWebhook)

	return controller
}

func (w *WebhookController) ProcessWebhook(c *gin.Context) {

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

	// Decode base64 data if needed (for future use)
	// if pubsubReq.Message.Data != "" {
	// 	if decodedData, err := base64.StdEncoding.DecodeString(pubsubReq.Message.Data); err == nil {
	// 		w.logger.Info(c.Request.Context(), "Decoded message data | dataLength: "+strconv.Itoa(len(decodedData)))
	// 	}
	// }

	// Extract secret name from secretId
	secretName, version, err := utils.ExtractSecretNameAndVersionFromVersionID(secretAttrs.VersionID)
	if err != nil {
		HandleError(c, w.logger, err)
		return
	}
	hashSecretName := utils.HashSecretName(secretName)
	labelKey := constants.CloudRunLabelPrefix + hashSecretName // Create label to search for services using this secret

	w.logger.Info(c.Request.Context(), "Searching for services | secretName: "+secretName+" | labelKey: "+labelKey)

	// Find services that use this secret
	services, err := w.webhookUsecase.GetServicesByLabel(c.Request.Context(), labelKey)
	if err != nil {
		HandleError(c, w.logger, err)
		return
	}

	// Log service names found and collect for response
	var serviceNames []string
	if len(services) > 0 {
		for i, service := range services {
			serviceName := "unknown"
			if service.Metadata != nil && service.Metadata.Name != "" {
				serviceName = service.Metadata.Name
				serviceNames = append(serviceNames, serviceName)
			}
			w.logger.Info(c.Request.Context(), "Service found | index: "+strconv.Itoa(i+1)+" | serviceName: "+serviceName+" | labelKey: "+labelKey)
		}
	}

	// Log successful HTTP response
	w.logger.Info(c.Request.Context(), "HTTP request completed | endpoint: /webhook | status: 200 | servicesFound: "+strconv.Itoa(len(services)))
	// Create new annotations and labels
	newAnnotations := map[string]string{
		constants.CloudRunLabelPrefix + hashSecretName + "/name":    secretName,
		constants.CloudRunLabelPrefix + hashSecretName + "/version": version,
	}
	newLabels := map[string]string{
		constants.CloudRunLabelPrefix + hashSecretName + "_version": version,
	}

	// Update services
	updatedServices, failedUpdates := w.webhookUsecase.UpdateCloudRunServices(c.Request.Context(), services, newAnnotations, newLabels)

	// Log HTTP response
	w.logger.Info(c.Request.Context(), "HTTP request completed | endpoint: /webhook | status: 200 | servicesFound: "+strconv.Itoa(len(services))+" | servicesUpdated: "+strconv.Itoa(len(updatedServices))+" | servicesFailed: "+strconv.Itoa(len(failedUpdates)))

	var failedUpdateNames []string
	for _, service := range failedUpdates {
		failedUpdateNames = append(failedUpdateNames, service.Metadata.Name)
	}
	var message string
	if len(updatedServices) > 0 {
		message = "Webhook processed successfully. Services updated: " + strconv.Itoa(len(updatedServices))
	} else {
		message = "Webhook processing failed. Services failed to update: " + strconv.Itoa(len(failedUpdates))
	}

	// Return structured response
	response := dto.WebhookResponse{
		Message:       message,
		TraceID:       logger.GetTraceIDFromContext(c.Request.Context()),
		ProcessedAt:   time.Now().UTC().Format(time.RFC3339),
		ServicesFound: len(serviceNames),
		ServiceNames:  serviceNames,
		FailedUpdates: failedUpdateNames,
	}

	c.JSON(http.StatusOK, response)
}
