package utils

import (
	"context"

	"github.com/aslammmuhammed/run-secret-reloader/pkg/logger"
)

// GetTraceID extracts trace ID from context, returns "unknown" if not found
func GetTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(logger.TraceIDKey).(string); ok && traceID != "" {
		return traceID
	}
	return "unknown"
}
