package controller

import (
	"context"

	"github.com/aslammmuhammed/run-secret-reloader/config"
	"github.com/aslammmuhammed/run-secret-reloader/internal/dependencies"
	"github.com/aslammmuhammed/run-secret-reloader/internal/middlewares"
	"github.com/aslammmuhammed/run-secret-reloader/internal/repo/cloudrun"
	"github.com/aslammmuhammed/run-secret-reloader/internal/usecase"
	"github.com/gin-gonic/gin"
)

func AddRouter(ctx context.Context, svc *dependencies.Dependencies, config *config.Config) {
	svc.Server.Use(middlewares.CorsMiddleware(config))
	svc.Server.Use(middlewares.TraceMiddleware())
	svc.Server.Use(middlewares.LoggerMiddleware(svc.Logger))

	// Health check endpoint
	svc.Server.GET("/health", GetHealth)

	v1ApiGroup := svc.Server.Group("/v1")
	initializeV1Routers(ctx, v1ApiGroup, svc)
}

func initializeV1Routers(ctx context.Context, server *gin.RouterGroup, svc *dependencies.Dependencies) {
	hashRouterGroup := server.Group("/hash")
	NewHashController(hashRouterGroup, svc.Logger, svc.Config)
	webhookRouterGroup := server.Group("/webhook")
	webhookRepo := cloudrun.NewCloudRunRepository(svc.CloudRunClient, ctx, svc.Config.ProjectID)
	webhookUsecase := usecase.NewWebhookUsecase(webhookRepo, svc.Logger)
	NewWebhookController(webhookRouterGroup, svc.Logger, svc.Config, webhookUsecase)
}
