terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 5.0"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

data "google_project" "this" {}

# Service account for Cloud Run and Pub/Sub push auth
resource "google_service_account" "reloader" {
  account_id   = var.service_account_name
  display_name = "Run Secret Reloader Service Account"
}

resource "google_project_iam_member" "cloudrun_admin" {
  project = data.google_project.this.project_id
  role    = "roles/run.admin"
  member  = "serviceAccount:${google_service_account.reloader.email}"
}

resource "google_project_iam_member" "service_account_user" {
  project = data.google_project.this.project_id
  role    = "roles/iam.serviceAccountUser"
  member  = "serviceAccount:${google_service_account.reloader.email}"
}



resource "google_project_iam_member" "secretmanager_viewer" {
  project = data.google_project.this.project_id
  role    = "roles/secretmanager.secretAccessor"
  member  = "serviceAccount:${google_service_account.reloader.email}"
}

# Cloud Run service
resource "google_cloud_run_v2_service" "reloader" {
  name     = var.cloud_run_service_name
  location = var.region

  template {
    service_account = google_service_account.reloader.email
    timeout         = "600s"
    containers {
      image = var.cloud_run_image
      resources {
        limits = {
          memory = var.cloud_run_memory
          cpu    = var.cloud_run_cpu
        }
        cpu_idle = true
      }

      liveness_probe {
        initial_delay_seconds = 10
        period_seconds        = 10
        timeout_seconds       = 5
        failure_threshold     = 3
        http_get {
          path = "/health"
          port = 8080
        }
      }
      startup_probe {
        initial_delay_seconds = 10
        period_seconds        = 10
        timeout_seconds       = 5
        failure_threshold     = 3
        http_get {
          path = "/health"
          port = 8080
        }
      }

      env {
        name  = "GOOGLE_CLOUD_PROJECT_ID"
        value = var.project_id
      }

      dynamic "env" {
        for_each = var.slack_webhook_url_secret_id != "" ? [1] : []
        content {
          name = "SLACK_WEBHOOK_URL"
          value_source {
            secret_key_ref {
              secret  = var.slack_webhook_url_secret_id
              version = var.slack_webhook_url_secret_version
            }
          }
        }
      }

      ports {
        container_port = 8080
      }
    }
    scaling {
      min_instance_count = 0
      max_instance_count = 5
    }
    max_instance_request_concurrency = 50
  }

  lifecycle {
    prevent_destroy = false
  }

}

# Pub/Sub topic for secret updates
resource "google_pubsub_topic" "topic" {
  name                       = var.topic_name
  message_retention_duration = var.message_retention_duration
}

# Push subscription with OIDC token using service account
resource "google_pubsub_subscription" "push" {
  name  = var.subscription_name
  topic = google_pubsub_topic.topic.name

  ack_deadline_seconds       = 600
  message_retention_duration = var.message_retention_duration

  push_config {
    push_endpoint = "${google_cloud_run_v2_service.reloader.uri}${var.push_endpoint_path}"

    oidc_token {
      service_account_email = google_service_account.reloader.email
    }
  }

  expiration_policy {
    ttl = "" #never
  }

  retry_policy {
    minimum_backoff = "300s"
    maximum_backoff = "500s"
  }
}

