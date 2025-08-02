package constants

import (
	"time"
)

const (
	UpdateRetries       int           = 3                      // Number of retries for updating Cloud Run services
	UpdateRetryDelay    time.Duration = 10 * time.Second       // Delay between retries
	CloudRunLabelPrefix string        = "run_secret_reloader-" // Prefix for Cloud Run labels
)
