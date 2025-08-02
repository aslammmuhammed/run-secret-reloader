package constants

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
