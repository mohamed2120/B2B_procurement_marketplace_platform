output "bus_name" {
  description = "EventBridge bus name"
  value       = aws_cloudwatch_event_bus.main.name
}

output "bus_arn" {
  description = "EventBridge bus ARN"
  value       = aws_cloudwatch_event_bus.main.arn
}

output "notification_queue_url" {
  description = "Notification service SQS queue URL"
  value       = aws_sqs_queue.notification.url
}

output "notification_queue_arn" {
  description = "Notification service SQS queue ARN"
  value       = aws_sqs_queue.notification.arn
}

output "search_indexer_queue_url" {
  description = "Search indexer service SQS queue URL"
  value       = aws_sqs_queue.search_indexer.url
}

output "search_indexer_queue_arn" {
  description = "Search indexer service SQS queue ARN"
  value       = aws_sqs_queue.search_indexer.arn
}
