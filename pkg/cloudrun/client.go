package cloudrun

import (
	"context"
	"fmt"

	"github.com/aslammmuhammed/run-secret-reloader/config"
	"google.golang.org/api/option"
	run "google.golang.org/api/run/v1"
)

// RunClient wraps the Google Cloud Run v1 client
type RunClient struct {
	serviceClient *run.APIService
	projectID     string
}

// NewClient creates a new CloudRun client with v1 service client only
func NewClient(ctx context.Context, config *config.Config) (*RunClient, error) {
	if config.ProjectID == "" {
		return nil, fmt.Errorf("project ID is required")
	}

	var opts []option.ClientOption

	// Add credentials file if provided
	if config.CredentialsFile != nil && *config.CredentialsFile != "" {
		opts = append(opts, option.WithCredentialsFile(*config.CredentialsFile))
	}

	// Add cloud platform scope
	opts = append(opts, option.WithScopes(run.CloudPlatformScope))

	// Create v1 service client
	serviceClient, err := run.NewService(ctx, opts...)
	if err != nil {
		// cancel() // Clean up if creation fails
		return nil, fmt.Errorf("failed to create v1 run client: %w", err)
	}

	return &RunClient{
		serviceClient: serviceClient,
		projectID:     config.ProjectID,
	}, nil
}

// GetServicesByLabel lists all Cloud Run services with a particular label key
// This method uses the v1 API with label selector
func (c *RunClient) GetServicesByLabel(ctx context.Context, labelKey string) ([]*run.Service, error) {
	if labelKey == "" {
		return nil, fmt.Errorf("label key is required")
	}

	// List all services across entire project (all regions)
	parent := fmt.Sprintf("namespaces/%s", c.projectID)

	// Call the API using Namespaces (V1) API with the provided request context
	// This ensures proper request tracing and timeout handling
	resp, err := c.serviceClient.Namespaces.Services.List(parent).Context(ctx).LabelSelector(labelKey).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list Cloud Run services: %w", err)
	}

	return resp.Items, nil
}


