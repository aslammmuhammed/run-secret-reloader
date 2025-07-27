package repo

import (
	"context"

	run "google.golang.org/api/run/v1"
)

type CloudRunRepository interface {
	GetServicesByLabel(ctx context.Context, labelKey string) ([]*run.Service, error)
	UpdateCloudRunAnnotationsAndLabels(ctx context.Context, serviceV1 *run.Service, newAnnotations map[string]string, newLabels map[string]string) error
}
