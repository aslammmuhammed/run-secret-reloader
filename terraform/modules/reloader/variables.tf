variable "project_id" {
  type        = string
  description = "GCP project ID"
}

variable "region" {
  type        = string
  description = "GCP region for Cloud Run"
  default     = "us-central1"
}


variable "topic_name" {
  type        = string
  description = "Pub/Sub topic name for secret events"
  default     = "run-secret-reloader"
}

variable "subscription_name" {
  type        = string
  description = "Pub/Sub push subscription name"
  default     = "run-secret-reloader-push"
}

variable "cloud_run_service_name" {
  type        = string
  description = "Name of the Cloud Run service"
  default     = "run-secret-reloader"
}

variable "cloud_run_image" {
  type        = string
  description = "Container image for Cloud Run (Docker Hub reference)"
  default     = "keyzersoze/run-secret-reloader:v1.0.0"
}

variable "cloud_run_memory" {
  type        = string
  description = "Cloud Run memory limit"
  default     = "256Mi"
}

variable "cloud_run_cpu" {
  type        = string
  description = "Cloud Run CPU limit"
  default     = "1"
}

variable "push_endpoint_path" {
  type        = string
  description = "Relative path on the Cloud Run service for push endpoint"
  default     = "/v1/webhook"
}

variable "service_account_name" {
  type        = string
  description = "Service account name for Cloud Run and Pub/Sub push"
  default     = "run-secret-reloader-sa"
}

variable "slack_webhook_url_secret_id" {
  type        = string
  description = "Secret ID for Slack webhook URL"
  default     = ""
}

variable "slack_webhook_url_secret_version" {
  type        = string
  description = "Secret version for Slack webhook URL"
  default     = "latest"
}

variable "message_retention_duration" {
  type        = string
  description = "Message retention duration for Pub/Sub"
  default     = "1200s"
}