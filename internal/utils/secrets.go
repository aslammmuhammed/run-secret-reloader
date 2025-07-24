package utils

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

// ExtractSecretName extracts the secret name from the full secret ID path
// Example: "projects/1234567890/secrets/user-service-env" -> "user-service-env"
func ExtractSecretName(secretID string) (string, error) {
	if secretID == "" {
		return "", fmt.Errorf("secretID cannot be empty")
	}

	parts := strings.Split(secretID, "/")
	if len(parts) < 4 {
		return "", fmt.Errorf("invalid secretID format: expected 'projects/PROJECT_ID/secrets/SECRET_NAME', got '%s'", secretID)
	}

	if parts[0] != "projects" {
		return "", fmt.Errorf("invalid secretID format: must start with 'projects', got '%s'", parts[0])
	}

	if parts[2] != "secrets" {
		return "", fmt.Errorf("invalid secretID format: expected 'secrets' at position 2, got '%s'", parts[2])
	}

	secretName := parts[3]
	if secretName == "" {
		return "", fmt.Errorf("secret name cannot be empty in secretID: '%s'", secretID)
	}

	return secretName, nil
}

// ExtractProjectID extracts the project ID from a GCP resource path
// Example: "projects/1234567890/secrets/user-service-env" -> "1234567890"
func ExtractProjectID(resourcePath string) (string, error) {
	if resourcePath == "" {
		return "", fmt.Errorf("resourcePath cannot be empty")
	}

	parts := strings.Split(resourcePath, "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid resource path format: expected 'projects/PROJECT_ID/...', got '%s'", resourcePath)
	}

	if parts[0] != "projects" {
		return "", fmt.Errorf("invalid resource path format: must start with 'projects', got '%s'", parts[0])
	}

	projectID := parts[1]
	if projectID == "" {
		return "", fmt.Errorf("project ID cannot be empty in resource path: '%s'", resourcePath)
	}

	return projectID, nil
}

// HashSecretName returns the first 32 hex characters of the SHA256 hash of the secret name
func HashSecretName(secretName string) string {
	hash := sha256.Sum256([]byte(secretName))
	return fmt.Sprintf("%x", hash[:16]) // 32 hex chars (16 bytes)
}
