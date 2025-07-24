package dependencies

import (
	"context"

	"github.com/aslammmuhammed/run-secret-reloader/config"
	"github.com/aslammmuhammed/run-secret-reloader/pkg/cloudrun"
	"github.com/aslammmuhammed/run-secret-reloader/pkg/logger"
	"github.com/aslammmuhammed/run-secret-reloader/pkg/server"
	"github.com/gin-gonic/gin"
)

// Service contains all dependencies for the application
type Dependencies struct {
	// Repositories

	// Use cases

	// Core dependencies
	Config         *config.Config
	Logger         logger.Logger
	Server         *gin.Engine
	CloudRunClient *cloudrun.RunClient
}

// NewService initializes and returns a new Service instance
func NewDependencies(ctx context.Context, cfg *config.Config) (*Dependencies, error) {
	// Initialize logger with proper log level
	log := logger.New(cfg.Logger.Format, cfg.Logger.Level)

	// Initialize server
	server := server.New()

	// Initialize cloudrun client
	cloudrunClient, err := cloudrun.NewClient(ctx, cfg)
	if err != nil {
		return nil, err
	}

	// TO DO Initialize repositories

	// TO DO Initialize use cases

	return &Dependencies{
		// Assign all dependencies
		Config:         cfg,
		Logger:         *log,
		Server:         server,
		CloudRunClient: cloudrunClient,
	}, nil
}
