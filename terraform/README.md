# 🚀 Run Secret Reloader - Terraform Module

This is the core Terraform module for deploying Run Secret Reloader to Google Cloud Platform. It automatically redeploys Cloud Run services when secrets change in Google Secret Manager.

## 🏗️ What Gets Created

The Terraform module provisions:

- **Cloud Run v2 service** with health checks and auto-scaling (0-5 instances)
- **Service Account** with required permissions:
  - `roles/run.admin`
  - `roles/iam.serviceAccountUser`
  - `roles/secretmanager.secretAccessor`
- **Pub/Sub topic** for receiving Secret Manager events
- **Push subscription** with OIDC authentication to Cloud Run service
- **IAM bindings** for secure service-to-service communication

## 🚀 Quick Start

### 1. Basic Deployment (With Slack Alerts)

```hcl
module "reloader" {
  source     = "github.com/aslammmuhammed/run-secret-reloader/terraform/modules/reloader"
  project_id = "your-gcp-project-id"
  
  # Slack alerts with Trace ID for debugging 
  slack_webhook_url_secret_id      = "slack-webhook-url"
  slack_webhook_url_secret_version = "latest"
}
```

### 2. Local Module Usage

```hcl
module "reloader" {
  source                           = "./terraform/modules/reloader"
  project_id                       = "your-project-id"
  cloud_run_image                  = "keyzersoze/run-secret-reloader:v0.1.0"
  
  # Slack alerts with Trace ID for debugging 
  slack_webhook_url_secret_id      = "slack-webhook-url"
  slack_webhook_url_secret_version = "latest"
}
```

### 3. Advanced Configuration

```hcl
module "reloader" {
  source     = "github.com/aslammmuhammed/run-secret-reloader/terraform/modules/reloader"
  project_id = "your-gcp-project-id"
  region     = "us-central1"
  
  # Custom configuration
  cloud_run_service_name = "my-secret-reloader"
  cloud_run_image        = "gcr.io/your-project/run-secret-reloader:latest"
  cloud_run_memory       = "512Mi"
  cloud_run_cpu          = "1"
  
  # Pub/Sub configuration
  topic_name        = "secret-events"
  subscription_name = "secret-events-push"
  
  # Service account
  service_account_name = "reloader-sa"
  
  # Slack configuration
  slack_webhook_url_secret_id      = "slack-webhook-url"
  slack_webhook_url_secret_version = "latest"
}
```

### 4. Deploy

```bash
# Initialize Terraform
terraform init

# Plan the deployment
terraform plan

# Apply the changes
terraform apply
```

## 📋 Module Variables

| Variable | Type | Description | Default | Required |
|----------|------|-------------|---------|----------|
| `project_id` | `string` | GCP project ID | - | ✅ |
| `region` | `string` | GCP region for Cloud Run | `us-central1` | ❌ |
| `topic_name` | `string` | Pub/Sub topic name | `run-secret-reloader` | ❌ |
| `subscription_name` | `string` | Pub/Sub subscription name | `run-secret-reloader-push` | ❌ |
| `cloud_run_service_name` | `string` | Cloud Run service name | `run-secret-reloader` | ❌ |
| `cloud_run_image` | `string` | Container image | `keyzersoze/run-secret-reloader:v0.1.0` | ❌ |
| `cloud_run_memory` | `string` | Memory limit | `256Mi` | ❌ |
| `cloud_run_cpu` | `string` | CPU limit | `1` | ❌ |
| `push_endpoint_path` | `string` | Webhook endpoint path | `/v1/webhook` | ❌ |
| `service_account_name` | `string` | Service account name | `run-secret-reloader-sa` | ❌ |
| `slack_webhook_url_secret_id` | `string` | Secret ID for Slack webhook URL | - | ✅ |
| `slack_webhook_url_secret_version` | `string` | Secret version for Slack webhook URL | - | ✅ |

> **Note:**  Alerts include a **Trace ID** for easy debugging. To disable alerts, leave the slack_webhook_url_secret_id empty.

## 📤 Module Outputs

| Output | Description |
|--------|-------------|
| `cloud_run_url` | URL of the deployed Cloud Run service |
| `pubsub_topic` | Name of the created Pub/Sub topic |
| `pubsub_subscription` | Name of the created Pub/Sub subscription |

## ⚙️ Post-Deployment Configuration

### 1. Configure Secret Manager Events

After deployment, configure your secrets to publish events to the created topic:

```bash
# For each secret that should trigger reloads
gcloud secrets update user-service-env \
  --add-topics=projects/YOUR_PROJECT_ID/topics/${TOPIC_NAME}
```
Make sure secret manager service account has publish to pubsub permissions , Refer
- [Secret Manager Event Notifications](https://cloud.google.com/secret-manager/docs/event-notifications)

### 2. Label Your Cloud Run Services

For any Cloud Run service that should reload when secrets change:

```bash
# Get the secret name hash for your secret name
SECRET_NAME_HASH=$(curl -X POST ${RELOADER_CLOUDRUN_URL}/v1/hash \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $(gcloud auth print-identity-token)" \
  -d '{"secretName": "user-service-env"}' | jq -r '.secretNameHash')

# Add the following label to your service
"run_secret_reloader-${SECRET_NAME_HASH}=true"
```

## 📚 Additional Resources


- [Pub/Sub Push Subscriptions](https://cloud.google.com/pubsub/docs/push)
- [Main Project Documentation](../README.md)

