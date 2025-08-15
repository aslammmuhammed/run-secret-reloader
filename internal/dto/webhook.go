package dto

// WebhookResponse represents the standard webhook response
type WebhookResponse struct {
	Message        string   `json:"message"`
	TraceID        string   `json:"traceId,omitempty"`
	ProcessedAt    string   `json:"processedAt,omitempty"`
	ServicesFound  int      `json:"servicesFound,omitempty"`
	ServiceNames   []string `json:"serviceNames,omitempty"`
	FailedUpdates  []string `json:"failedUpdates,omitempty"`
	SkippedUpdates []string `json:"skippedUpdates,omitempty"`
}
