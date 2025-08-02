package cloudrun

import (
	"context"
	"fmt"
	"maps"

	run_v2 "cloud.google.com/go/run/apiv2"
	runpb_v2 "cloud.google.com/go/run/apiv2/runpb"
	"github.com/aslammmuhammed/run-secret-reloader/config"
	"github.com/aslammmuhammed/run-secret-reloader/pkg/logger"
	"google.golang.org/api/option"
	run "google.golang.org/api/run/v1"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// RunClient wraps the Google Cloud Run v1 client
type RunClient struct {
	serviceClient    *run.APIService
	projectID        string
	serviceClient_v2 *run_v2.ServicesClient
	logger           *logger.Logger
}

// NewClient creates a new CloudRun client with v1 service client only
func NewClient(ctx context.Context, config *config.Config, logger *logger.Logger) (*RunClient, error) {
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

	// Create v2 service client
	serviceClient_v2, err := run_v2.NewServicesClient(ctx, opts...)
	if err != nil {
		// cancel() // Clean up if creation fails
		return nil, fmt.Errorf("failed to create v2 run client: %w", err)
	}

	return &RunClient{
		serviceClient:    serviceClient,
		projectID:        config.ProjectID,
		serviceClient_v2: serviceClient_v2,
		logger:           logger,
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

	labelSelector := fmt.Sprintf("%s=true", labelKey)

	c.logger.Debug(ctx, "Using label selector: "+labelSelector)

	var services []*run.Service
	var pageToken string
	for {
		// Create the list request with pagination and label selector
		listCall := c.serviceClient.Namespaces.Services.List(parent).
			Context(ctx).
			LabelSelector(labelSelector)
		// Add page token from a previous response
		if pageToken != "" {
			listCall = listCall.Continue(pageToken)
		}
		resp, err := listCall.Do()
		if err != nil {
			return nil, fmt.Errorf("failed to list Cloud Run services: %w", err)
		}
		services = append(services, resp.Items...)

		// continuation token for the next page
		if resp.Metadata == nil || resp.Metadata.Continue == "" {
			// No more pages
			break
		}
		pageToken = resp.Metadata.Continue
		c.logger.Debug(ctx, fmt.Sprintf("Fetching next page with token: %s", pageToken))
	}

	c.logger.Info(ctx, fmt.Sprintf("Found %d services with label %s", len(services), labelKey))
	return services, nil
}

func (c *RunClient) UpdateCloudRunAnnotationsAndLabels(ctx context.Context, serviceV1 *run.Service, newAnnotations map[string]string, newLabels map[string]string) error {

	serviceName := serviceV1.Metadata.Name
	region := serviceV1.Metadata.Labels["cloud.googleapis.com/location"]
	fullName := fmt.Sprintf("projects/%s/locations/%s/services/%s", c.projectID, region, serviceName)
	// Fetch current service
	c.logger.Debug(ctx, "Fetching service: "+fullName)
	service, err := c.serviceClient_v2.GetService(ctx, &runpb_v2.GetServiceRequest{Name: fullName})
	if err != nil {
		return fmt.Errorf("failed to get service: %w", err)
	}

	// Initialize and update template annotations if provided
	if len(newAnnotations) > 0 {
		if service.Template == nil {
			service.Template = &runpb_v2.RevisionTemplate{}
		}
		if service.Template.Annotations == nil {
			service.Template.Annotations = make(map[string]string)
		}

		c.logger.Debug(ctx, "Template annotations before update:"+fmt.Sprintf("%+v", service.Template.Annotations))

		// Apply new annotations to spec.template.annotations
		maps.Copy(service.Template.Annotations, newAnnotations)
	}

	// Initialize and update template labels if provided
	if len(newLabels) > 0 {
		if service.Template == nil {
			service.Template = &runpb_v2.RevisionTemplate{}
		}
		if service.Template.Labels == nil {
			service.Template.Labels = make(map[string]string)
		}

		c.logger.Debug(ctx, "Template labels before update:"+fmt.Sprintf("%+v", service.Template.Labels))

		// Apply new labels to spec.template.metadata.labels
		maps.Copy(service.Template.Labels, newLabels)
	}

	fieldPaths := []string{"template.annotations", "template.labels"}

	c.logger.Debug(ctx, "Using UpdateMask:"+fmt.Sprintf("%v", fieldPaths))

	// Create the update request with specific field mask for both annotations and labels
	req := &runpb_v2.UpdateServiceRequest{
		Service: service,
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: fieldPaths, // Update specific annotation and label keys
		},
	}
	// Trigger the update
	op, err := c.serviceClient_v2.UpdateService(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to update service: %w", err)
	}

	// Wait for operation to complete with better error handling
	result, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("failed to update service: %w", err)
	}

	c.logger.Debug(ctx, "Update operation completed successfully:"+result.Name)
	if len(newAnnotations) > 0 {
		c.logger.Debug(ctx, "Template annotations after update:"+fmt.Sprintf("%+v", service.Template.Annotations))
	}
	if len(newLabels) > 0 {
		c.logger.Debug(ctx, "Template labels after update:"+fmt.Sprintf("%+v", service.Template.Labels))
	}

	return nil
}
