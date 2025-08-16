# Run Secret Reloader - Terraform Module

This is the core Terraform module for deploying Run Secret Reloader to Google Cloud Platform.

## Usage

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

## Important Setup

After deploying the module, you must configure Secret Manager to publish events to the created Pub/Sub topic:

```bash
# For each secret that should trigger reloads
gcloud secrets update SECRET_NAME \
  --add-topics=projects/YOUR_PROJECT_ID/topics/run-secret-reloader
```

**📖 For complete documentation, see [../README.md](../README.md)**

## Requirements

- Google Cloud APIs enabled:
  - Cloud Run API
  - Pub/Sub API  
  - Secret Manager API
- Service account with proper permissions
- Terraform >= 1.0

## Resources Created

- Cloud Run v2 service
- Service account with Cloud Run admin permissions
- Pub/Sub topic and push subscription
- IAM policy bindings

Refer to [Secret Manager Event Notifications](https://cloud.google.com/secret-manager/docs/event-notifications) for additional setup requirements.