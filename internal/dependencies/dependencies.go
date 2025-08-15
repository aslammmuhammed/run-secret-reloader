package dependencies

import (
	"context"

	"github.com/aslammmuhammed/run-secret-reloader/config"
	"github.com/aslammmuhammed/run-secret-reloader/pkg/alert"
	"github.com/aslammmuhammed/run-secret-reloader/pkg/alert/slack"
	"github.com/aslammmuhammed/run-secret-reloader/pkg/cloudrun"
	"github.com/aslammmuhammed/run-secret-reloader/pkg/logger"
	"github.com/aslammmuhammed/run-secret-reloader/pkg/server"
)

// Service contains all dependencies for the application
type Dependencies struct {
	// Repositories

	// Use cases

	// Core dependencies
	Config         *config.Config
	Logger         logger.Logger
	Server         *server.Server
	CloudRunClient *cloudrun.RunClient
	Alert          alert.Client
}

// NewService initializes and returns a new Service instance
func NewDependencies(ctx context.Context, cfg *config.Config) (*Dependencies, error) {
	// Initialize logger with proper log level
	log := logger.New(cfg.Logger.Format, cfg.Logger.Level)

	// Initialize server
	server := server.New()

	// Initialize cloudrun client
	cloudrunClient, err := cloudrun.NewClient(ctx, cfg, log)
	if err != nil {
		return nil, err
	}

	// TO DO Initialize repositories

	// TO DO Initialize use cases

	// Initialize alert client
	var alertClient alert.Client
	if cfg.Alert.Provider != "" {
		log.Info(ctx, "Initializing alert client | provider: "+cfg.Alert.Provider)
		switch cfg.Alert.Provider {
		case alert.AlertTypeSlack:
			log.Info(ctx, "Initializing slack alert client")
			alertClient, err = alert.NewClient(cfg, log, slack.RegisterProvider)
			if err != nil {
				return nil, err
			}
		}
	} else {
		log.Info(ctx, "No alert provider configured")
	}

	return &Dependencies{
		// Assign all dependencies
		Config:         cfg,
		Logger:         *log,
		Server:         server,
		CloudRunClient: cloudrunClient,
		Alert:          alertClient,
	}, nil
}

func (d *Dependencies) Close() error {
	if d.CloudRunClient != nil {
		return d.CloudRunClient.Close()
	}
	return nil
}
