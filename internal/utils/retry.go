package utils

import (
	"slices"

	"github.com/aslammmuhammed/run-secret-reloader/internal/constants"
	"google.golang.org/grpc/status"
)

// IsRetryableError checks if the given error is one that should trigger a retry
func IsRetryableError(err error) bool {
	if st, ok := status.FromError(err); ok {
		return slices.Contains(constants.GetRetryableGrpcCodes(), st.Code())
	}
	return false
}
