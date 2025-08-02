package constants

import "google.golang.org/grpc/codes"

// Retryable gRPC error codes
const (
	// RetryableErrorAborted is for concurrency issues
	RetryableErrorAborted = codes.Aborted
	// RetryableErrorOutOfRange is for fixable range issues
	RetryableErrorOutOfRange = codes.OutOfRange
	// RetryableErrorFailedPrecondition is for temporary system state issues
	RetryableErrorFailedPrecondition = codes.FailedPrecondition
)

// GetRetryableGrpcCodes returns the slice of error codes that should trigger a retry
func GetRetryableGrpcCodes() []codes.Code {
	return []codes.Code{
		RetryableErrorAborted,
		RetryableErrorOutOfRange,
		RetryableErrorFailedPrecondition,
	}
}
