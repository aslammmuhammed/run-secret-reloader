package utils

// Google Cloud Secret Manager Event Types
// Reference: https://cloud.google.com/secret-manager/docs/event-notifications
const (
	// Core secret lifecycle events
	SecretCreate = "SECRET_CREATE" // Sent when a new secret is successfully created
	SecretUpdate = "SECRET_UPDATE" // Sent when a new secret is successfully updated
	SecretDelete = "SECRET_DELETE" // Sent when a secret is deleted

	// Secret version lifecycle events
	SecretVersionAdd              = "SECRET_VERSION_ADD"               // Sent when a new secret version is successfully added
	SecretVersionEnable           = "SECRET_VERSION_ENABLE"            // Sent when a secret version is enabled
	SecretVersionDisable          = "SECRET_VERSION_DISABLE"           // Sent when a secret version is disabled
	SecretVersionDestroy          = "SECRET_VERSION_DESTROY"           // Sent when a secret version is destroyed
	SecretVersionDestroyScheduled = "SECRET_VERSION_DESTROY_SCHEDULED" // Sent when destruction is scheduled

	// Rotation and configuration events
	SecretRotate    = "SECRET_ROTATE"    // Sent when it is time to rotate a secret
	TopicConfigured = "TOPIC_CONFIGURED" // Test message sent when topics are configured
)

// ShouldProcessSecretEvent determines if the given event type requires service reloading
// Currently focused on SECRET_VERSION_ADD as this indicates a new secret value is available
func ShouldProcessSecretEvent(eventType string) bool {
	switch eventType {
	case SecretVersionAdd:
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
		SecretCreate:                  "New secret successfully created",
		SecretUpdate:                  "Secret successfully updated",
		SecretDelete:                  "Secret deleted",
		SecretVersionAdd:              "New secret version successfully added",
		SecretVersionEnable:           "Secret version enabled",
		SecretVersionDisable:          "Secret version disabled",
		SecretVersionDestroy:          "Secret version destroyed",
		SecretVersionDestroyScheduled: "Secret version destruction scheduled",
		SecretRotate:                  "Time to rotate secret",
		TopicConfigured:               "Pub/Sub topics configured (test message)",
	}

	if description, exists := descriptions[eventType]; exists {
		return description
	}
	return "Unknown event type"
}
