package alert

import (
	"fmt"
)

const (
	// AlertTypeSlack represents the Slack alert provider
	AlertTypeSlack = "slack"
	// Add other alert types as needed
)

// Factory creates alert clients based on the provider type
type Factory struct {
	providers map[string]Provider
}

// newFactory creates a new alert factory
func newFactory() *Factory {
	return &Factory{
		providers: make(map[string]Provider),
	}
}

// Register registers an alert provider with the factory
func (f *Factory) Register(providerType string, provider Provider) {
	f.providers[providerType] = provider
}

// create creates a new alert client of the specified type
func (f *Factory) create(providerType string) (Client, error) {
	provider, exists := f.providers[providerType]
	if !exists {
		return nil, fmt.Errorf("alert provider %s not registered", providerType)
	}

	return provider()
}
