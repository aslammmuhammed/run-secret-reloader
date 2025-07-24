package repo

import (
	"context"

	run "google.golang.org/api/run/v1"
)

type CloudRunRepository interface {
	GetServicesByLabel(ctx context.Context, labelKey string) ([]*run.Service, error)
}
