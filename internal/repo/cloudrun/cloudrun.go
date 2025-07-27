package cloudrun

import (
	"context"
	"fmt"

	"github.com/aslammmuhammed/run-secret-reloader/pkg/cloudrun"
	run "google.golang.org/api/run/v1"
)

type cloudRunRepository struct {
	cloudrunClient *cloudrun.RunClient
	projectID      string
}

func NewCloudRunRepository(cloudrunClient *cloudrun.RunClient, ctx context.Context, projectID string) *cloudRunRepository {
	return &cloudRunRepository{
		cloudrunClient: cloudrunClient,
		projectID:      projectID,
	}
}

func (r *cloudRunRepository) GetServicesByLabel(ctx context.Context, labelKey string) ([]*run.Service, error) {
	if labelKey == "" {
		return nil, fmt.Errorf("label key is required")
	}

	services, err := r.cloudrunClient.GetServicesByLabel(ctx, labelKey)
	if err != nil {
		return nil, fmt.Errorf("failed to list Cloud Run services: %w", err)
	}

	return services, nil
}

func (r *cloudRunRepository) UpdateCloudRunAnnotationsAndLabels(ctx context.Context, serviceV1 *run.Service, newAnnotations map[string]string, newLabels map[string]string) error {
	return r.cloudrunClient.UpdateCloudRunAnnotationsAndLabels(ctx, serviceV1, newAnnotations, newLabels)
}