package alert

import (
	"fmt"

	"github.com/aslammmuhammed/run-secret-reloader/config"
	"github.com/aslammmuhammed/run-secret-reloader/pkg/logger"
)

// ProviderRegistrar defines the function signature for registering an alert provider.
type ProviderRegistrar func(factory *Factory, cfg *config.Config, logger *logger.Logger)

// Provider represents a function that creates an alert client.
type Provider func() (Client, error)

// NewClient creates a new alert client based on the configuration.
func NewClient(cfg *config.Config, logger *logger.Logger, registrars ...ProviderRegistrar) (Client, error) {
	factory := newFactory()

	// Register providers using the provided registrar functions
	for _, registrar := range registrars {
		registrar(factory, cfg, logger)
	}

	// Create client based on configured provider
	client, err := factory.create(cfg.Alert.Provider)
	if err != nil {
		return nil, fmt.Errorf("failed to create alert client: %w", err)
	}

	return client, nil
}
