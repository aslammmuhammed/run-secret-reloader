package app

import (
	"github.com/aslammmuhammed/run-secret-reloader/config"
	"github.com/aslammmuhammed/run-secret-reloader/internal/controller"
	"github.com/aslammmuhammed/run-secret-reloader/internal/dependencies"
	"github.com/aslammmuhammed/run-secret-reloader/pkg/logger"
	"golang.org/x/net/context"
)

func Run() {
	config, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	logger := logger.NewRootLogger(config.Logger.Format, config.Logger.Level)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	svc, err := dependencies.NewDependencies(ctx, config)
	if err != nil {
		logger.Fatal("Failed to initialize service:", err)
	}
	// TODO: Add graceful shutdown and closing of dependencies
	controller.AddRouter(ctx, svc, config)
	logger.Info("Server started listening on port " + config.HTTP.Port)
	if err := svc.Server.Run(config.HTTP.Port); err != nil {
		logger.Fatal("Failed to start server:", err)
	}
}
