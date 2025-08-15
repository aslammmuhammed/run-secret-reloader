package alert

import (
	"context"
)

// Message represents a structured alert message for the secret reloader service.
type Message struct {
	SecretName string
	NewVersion string
	TraceID    string
	Succeeded  []string
	Failed     []string
	Skipped    []string
}

// Client represents an alert client interface
type Client interface {
	// SendMessage sends a message to an alert provider.
	SendMessage(ctx context.Context, message Message) error
}
