package dto

import "time"

// PubSubRequest represents the incoming Pub/Sub webhook request
type PubSubRequest struct {
	Message      PubSubMessage `json:"message"`
	Subscription string        `json:"subscription"`
}

// PubSubMessage represents the Pub/Sub message structure
type PubSubMessage struct {
	Attributes  SecretAttributes `json:"attributes"`
	Data        string           `json:"data"` // base64 encoded
	MessageID   string           `json:"messageId"`
	PublishTime time.Time        `json:"publishTime"`
}

// SecretAttributes represents the secret manager event attributes
type SecretAttributes struct {
	DataFormat string    `json:"dataFormat"` // "JSON_API_V1"
	EventType  string    `json:"eventType"`  // "SECRET_VERSION_ADD"
	SecretID   string    `json:"secretId"`   // "projects/1234567890/secrets/user-service-env"
	Timestamp  time.Time `json:"timestamp"`
	VersionID  string    `json:"versionId"` // "projects/1234567890/secrets/user-service-env/versions/3"
}

