# 🚀 Terraform Deployment

This directory contains the Terraform module for deploying Run Secret Reloader to Google Cloud Platform.

## 📋 Prerequisites

- Google Cloud Project with billing enabled
- Terraform >= 1.0
- `gcloud` CLI configured and authenticated
- Required APIs enabled:
  ```bash
  gcloud services enable run.googleapis.com
  gcloud services enable pubsub.googleapis.com
  gcloud services enable secretmanager.googleapis.com
  gcloud services enable cloudbuild.googleapis.com
  ```

## 🏗️ What Gets Created

The Terraform module provisions:

- **Cloud Run v2 service** with health checks and auto-scaling (0-5 instances)
- **Service Account** with required permissions:
  - `roles/run.admin` - To update Cloud Run services
  - `roles/iam.serviceAccountUser` - To manage service accounts
- **Pub/Sub topic** for receiving Secret Manager events
- **Push subscription** with OIDC authentication to Cloud Run service
- **IAM bindings** for secure service-to-service communication

## 🚀 Quick Start

### 1. Basic Deployment (With Slack Alerts)

```hcl
module "reloader" {
  source     = "github.com/aslammmuhammed/run-secret-reloader//terraform/modules/reloader"
  project_id = "your-gcp-project-id"
  
  slack_webhook_url_secret_id      = "slack-webhook-url"
  slack_webhook_url_secret_version = "latest"
}
```

### 2. Advanced Configuration

```hcl
module "reloader" {
  source     = "github.com/aslammmuhammed/run-secret-reloader//terraform/modules/reloader"
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
}
```

### 3. Deploy

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
| `secrets` | `list(string)` | List of secret names to monitor | - | ✅ |
| `topic_name` | `string` | Pub/Sub topic name | `run-secret-reloader` | ❌ |
| `subscription_name` | `string` | Pub/Sub subscription name | `run-secret-reloader-push` | ❌ |
| `cloud_run_service_name` | `string` | Cloud Run service name | `run-secret-reloader` | ❌ |
| `cloud_run_image` | `string` | Container image | `keyzersoze/run-secret-reloader:v0.1.0` | ❌ |
| `cloud_run_memory` | `string` | Memory limit | `256Mi` | ❌ |
| `cloud_run_cpu` | `string` | CPU limit | `1` | ❌ |
| `push_endpoint_path` | `string` | Webhook endpoint path | `/v1/webhook` | ❌ |
| `service_account_name` | `string` | Service account name | `run-secret-reloader-sa` | ❌ |

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
# Get the topic name from Terraform output
TOPIC_NAME=$(terraform output -raw pubsub_topic)

# For each secret that should trigger reloads
gcloud secrets update SECRET_NAME \
  --add-topics=projects/YOUR_PROJECT_ID/topics/${TOPIC_NAME}
```

### 2. Label Your Cloud Run Services

For any Cloud Run service that should reload when secrets change:

```bash
# Get the reloader service URL
RELOADER_URL=$(terraform output -raw cloud_run_url)

# Get the hash for your secret name
SECRET_HASH=$(curl -X POST ${RELOADER_URL}/v1/hash \
  -H "Content-Type: application/json" \
  -d '{"secretName": "db-password"}' | jq -r '.secretNameHash')

# Add the label to your service
gcloud run services update YOUR_SERVICE_NAME \
  --region=YOUR_REGION \
  --update-labels="run_secret_reloader-${SECRET_HASH}=true"
```

## 🔍 Verification

### Test the Deployment

```bash
# Check if the service is running
RELOADER_URL=$(terraform output -raw cloud_run_url)
curl ${RELOADER_URL}/health

# Test the hash endpoint
curl -X POST ${RELOADER_URL}/v1/hash \
  -H "Content-Type: application/json" \
  -d '{"secretName": "test-secret"}'
```

### Monitor Logs

```bash
gcloud logs read "resource.type=cloud_run_revision AND resource.labels.service_name=run-secret-reloader" \
  --project=YOUR_PROJECT_ID \
  --limit=20
```

## 🛠️ Troubleshooting

### Common Issues

**1. Permission Denied Errors**
```bash
# Ensure APIs are enabled
gcloud services enable run.googleapis.com pubsub.googleapis.com secretmanager.googleapis.com

# Check your authentication
gcloud auth list
gcloud auth application-default login
```

**2. Service Account Issues**
```bash
# Verify service account permissions
gcloud projects get-iam-policy YOUR_PROJECT_ID \
  --filter="bindings.members:serviceAccount:run-secret-reloader-sa@YOUR_PROJECT_ID.iam.gserviceaccount.com"
```

**3. Pub/Sub Subscription Errors**
```bash
# Check if subscription exists and is configured correctly
gcloud pubsub subscriptions describe run-secret-reloader-push
```

## 📚 Additional Resources

- [Secret Manager Event Notifications](https://cloud.google.com/secret-manager/docs/event-notifications)
- [Cloud Run IAM Roles](https://cloud.google.com/run/docs/reference/iam/roles)
- [Pub/Sub Push Subscriptions](https://cloud.google.com/pubsub/docs/push)

## 🔄 Cleanup

To remove all created resources:

```bash
terraform destroy
```

**Note**: This will delete the Cloud Run service, Pub/Sub resources, and service account. Make sure to remove the topic references from your secrets first:

```bash
gcloud secrets update SECRET_NAME --clear-topics
```