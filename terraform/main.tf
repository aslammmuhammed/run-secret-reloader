module "reloader" {
  source                           = "./modules/reloader"
  project_id                       = "your-project-id"
  cloud_run_image                  = "keyzersoze/run-secret-reloader:v0.1.0"
  slack_webhook_url_secret_id      = "slack-webhook-url"
  slack_webhook_url_secret_version = "latest"
}
