package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aslammmuhammed/run-secret-reloader/config"
	"github.com/aslammmuhammed/run-secret-reloader/internal/controller"
	"github.com/aslammmuhammed/run-secret-reloader/internal/dependencies"
	"github.com/aslammmuhammed/run-secret-reloader/pkg/logger"
)

// Run starts the application and returns an exit code upon shutdown.
func Run() int {
	config, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	appLogger := logger.NewRootLogger(config.Logger.Format, config.Logger.Level)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	svc, err := dependencies.NewDependencies(ctx, config)
	if err != nil {
		appLogger.Fatal("Failed to initialize service:", err)
	}
	defer func() {
		if err := svc.Close(); err != nil {
			appLogger.Error("Failed to close dependencies:", err)
		}
	}()

	controller.AddRouter(ctx, svc, config)

	// Channel to catch server errors
	serverError := make(chan error, 1)

	// Graceful shutdown
	go func() {
		appLogger.Info("Server started listening on port " + config.HTTP.Port)
		if err := svc.Server.Run(config.HTTP.Port); err != nil && err != http.ErrServerClosed {
			serverError <- err
		}
	}()

	// Wait for either server error or shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	exitCode := 0
	select {
	case err := <-serverError:
		appLogger.Error("Server error, shutting down:", err)
		exitCode = 1
	case <-quit:
		appLogger.Info("Shutdown signal received...")
	}

	appLogger.Info("Shutting down server...")

	// The context is used to cancel the server shutdown if it takes too long
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := svc.Server.Shutdown(shutdownCtx); err != nil {
		appLogger.Error("Server forced to shutdown:", err)
	}

	appLogger.Info("Server exiting")
	return exitCode
}
