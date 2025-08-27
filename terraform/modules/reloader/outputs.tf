output "cloud_run_url" {
  value = google_cloud_run_v2_service.reloader.uri
}

output "pubsub_topic" {
  value = google_pubsub_topic.topic.name
}

output "pubsub_subscription" {
  value = google_pubsub_subscription.push.name
}