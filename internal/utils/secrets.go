package utils

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

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

// ExtractSecretNameAndVersionFromVersionID extracts the secret name and version from the secret version ID
// Example: "projects/1234567890/secrets/user-service-env/versions/3" -> "user-service-env", "3"
func ExtractSecretNameAndVersionFromVersionID(versionID string) (string, string, error) {
	if versionID == "" {
		return "", "", fmt.Errorf("versionID cannot be empty")
	}

	parts := strings.Split(versionID, "/")
	if len(parts) < 6 {
		return "", "", fmt.Errorf("invalid versionID format: expected 'projects/PROJECT_ID/secrets/SECRET_NAME/versions/VERSION', got '%s'", versionID)
	}

	if parts[0] != "projects" {
		return "", "", fmt.Errorf("invalid versionID format: must start with 'projects', got '%s'", parts[0])
	}

	if parts[2] != "secrets" {
		return "", "", fmt.Errorf("invalid versionID format: expected 'secrets' at position 2, got '%s'", parts[2])
	}

	if parts[4] != "versions" {
		return "", "", fmt.Errorf("invalid versionID format: expected 'versions' at position 4, got '%s'", parts[4])
	}

	secretName := parts[3]
	if secretName == "" {
		return "", "", fmt.Errorf("secret name cannot be empty in versionID: '%s'", versionID)
	}

	version := parts[5]
	if version == "" {
		return "", "", fmt.Errorf("version cannot be empty in versionID: '%s'", versionID)
	}

	return secretName, version, nil
}
