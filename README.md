# 🔄 Run Secret Reloader

[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go](https://img.shields.io/badge/go-%23007BCD.svg?style=flat&logo=go&logoColor=white)](https://golang.org/)
[![Google Cloud](https://img.shields.io/badge/GoogleCloud-%234285F4.svg?style=flat&logo=google-cloud&logoColor=white)](https://cloud.google.com/)

**Automatically redeploy Cloud Run services when secrets change in Google Secret Manager.**

Stop manual redeployments and let your services automatically pick up fresh secrets! This service listens for Secret Manager events and intelligently reloads only the Cloud Run services that depend on the updated secret.

## ✨ Features

- 🔄 **Automatic Redeployment** - No manual intervention needed
- 🏷️ **Label-based Targeting** - Only reloads services that depend on the changed secret
- 📦 **Zero Code Changes** - Works with existing applications using `os.Getenv()`
- 🔒 **Secure** - Uses OIDC authentication and proper IAM roles
- 📊 **Smart Versioning** - Skips redeployment if service already has newer secret version
- 🚨 **Alerts (Expandable)** - Receive notifications about successful or failed reloads, with Trace ID for debugging. Currently, only Slack is supported.
- 🌐 **Multi-service Support** - Handle multiple Cloud Run services per secret
- ☁️ **Serverless** - Runs as a Cloud Run service, scales to zero when idle

## 🚀 Quick Start

### 1. Deploy with Terraform

→ **[See complete deployment guide](./terraform/README.md)**

### 2. Update Your Secret

```bash
echo -n "new-password-value" | gcloud secrets versions add user-service-env --data-file=-
```

### 3. Label Your Services

Add the label to services that should reload when secrets change:

```bash
# For a service using "user-service-env" secret
SECRET_NAME_HASH=$(curl -X POST https://YOUR_RELOADER_URL/v1/hash \
  -H "Content-Type: application/json" \
  -d '{"secretName": "user-service-env"}' | jq -r '.secretNameHash')

"run_secret_reloader-${SECRET_NAME_HASH}=true"
```

🎉 **That's it!** Your Cloud Run service will automatically redeploy with the new secret.

Cloud Run has built-in [secret hot-reload capabilities](https://medium.com/google-cloud/cloud-run-hot-reload-your-secret-manager-secrets-ff2c502df666), but most applications read secrets only once at startup. This project bridges that gap by automatically triggering redeployments when secrets change.

**The Problem:**
- Environment variables are read once at startup
- File mounts fetch fresh values, but apps cache secrets in memory
- Manual redeployments are error-prone and slow

**The Solution:**
- Automatic redeployments triggered by Secret Manager events
- Works with existing applications using `os.Getenv()`
- Label-based targeting ensures only dependent services reload
- Mount each secret as a volume to make it available to the container as files. Reading a volume always fetches the secret value from Secret Manager, so it can be used with the latest version.

## 🛠️ API Reference

### Hash Endpoint

```bash
POST /v1/hash
{
  "secretName": "my-secret-name"
}
```

**Response:**
```json
{
  "secretNameHash": "2dbfe2930811fa77b4fc5105492cc6e6"
}
```

### Webhook Endpoint

```bash
POST /v1/webhook
```

## 🏗️ Architecture & Technical Details

### Error Handling & Retries

- **Retryable errors**: `ABORTED`, `OUT_OF_RANGE`, `FAILED_PRECONDITION`
- **Max retries**: 3 attempts with exponential backoff
- **Base delay**: 10 seconds between retries
- **Optimistic locking**: Uses etags to prevent conflicts in Cloud Run update

### How It Works Internally

1. **Event Processing** - Receives Pub/Sub push notifications from Secret Manager
2. **Service Discovery** - Uses Cloud Run v1 API to list services with matching labels
3. **Template Updates** - Uses Cloud Run v2 API to update service annotations/labels
4. **Revision Creation** - Cloud Run automatically creates new revision when template changes

### Label and Annotation Schema

**Service Labels** (for targeting):
- `run_secret_reloader-<hash>=true` (where `<hash>` is MD5 of secret name)

**Template Annotations** (for tracking):
- `run_secret_reloader-<hash>/name` → secret name
- `run_secret_reloader-<hash>/version` → secret version

**Template Labels** (for tracking):
- `run_secret_reloader-<hash>_version` → secret version

### Project Structure

```
cmd/app/                – program entrypoint
internal/               – application code
  app/                  – application initialization
  controller/           – HTTP handlers (webhook, health, hash)
  usecase/              – business logic (webhook processing)
  repo/                 – Cloud Run API client implementations
  dto/                  – data transfer objects
  utils/                – helpers (events, retry, secrets)
  constants/            – constants (retries, prefixes, errors)
  dependencies/         – dependency injection
  errors/               – custom error types
  middlewares/          – HTTP middlewares
pkg/                    – reusable packages
  alert/                – alerting providers (Slack, extensible)
  cloudrun/             – Cloud Run API client wrapper
  server/               – HTTP server implementation  
  logger/               – structured logging
config/                 – runtime configuration
terraform/              – infrastructure as code
```

## 📜 License

Apache License 2.0 © 2025 Aslam Muhammed

---

Built with ❤️ for the Google Cloud community
