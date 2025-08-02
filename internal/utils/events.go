package utils

import "github.com/aslammmuhammed/run-secret-reloader/internal/constants"

// ShouldProcessSecretEvent determines if the given event type requires service reloading
// Currently focused on SECRET_VERSION_ADD as this indicates a new secret value is available
func ShouldProcessSecretEvent(eventType string) bool {
	switch eventType {
	case constants.SecretVersionAdd:
		// New secret version added - services need to reload to get new value
		return true

	// Future: Add other events that might require service reload
	// case SecretVersionEnable:
	//     return true
	// case SecretCreate:
	//     return true

	default:
		// All other events don't require service reload
		return false
	}
}

// GetEventDescription returns a human-readable description of the event type
func GetEventDescription(eventType string) string {
	descriptions := map[string]string{
		constants.SecretCreate:                  "New secret successfully created",
		constants.SecretUpdate:                  "Secret successfully updated",
		constants.SecretDelete:                  "Secret deleted",
		constants.SecretVersionAdd:              "New secret version successfully added",
		constants.SecretVersionEnable:           "Secret version enabled",
		constants.SecretVersionDisable:          "Secret version disabled",
		constants.SecretVersionDestroy:          "Secret version destroyed",
		constants.SecretVersionDestroyScheduled: "Secret version destruction scheduled",
		constants.SecretRotate:                  "Time to rotate secret",
		constants.TopicConfigured:               "Pub/Sub topics configured (test message)",
	}

	if description, exists := descriptions[eventType]; exists {
		return description
	}
	return "Unknown event type"
}
