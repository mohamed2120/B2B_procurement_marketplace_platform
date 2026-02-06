# EventBridge Custom Bus
resource "aws_cloudwatch_event_bus" "main" {
  name = "${var.environment}-b2b-platform"

  tags = {
    Name = "${var.environment}-eventbridge-bus"
  }
}

# SQS Queue for Notification Service
resource "aws_sqs_queue" "notification" {
  name                      = "${var.environment}-notification-queue"
  message_retention_seconds = 1209600 # 14 days
  visibility_timeout_seconds = 300

  tags = {
    Name = "${var.environment}-notification-queue"
  }
}

# SQS Queue for Search Indexer Service
resource "aws_sqs_queue" "search_indexer" {
  name                      = "${var.environment}-search-indexer-queue"
  message_retention_seconds = 1209600 # 14 days
  visibility_timeout_seconds = 300

  tags = {
    Name = "${var.environment}-search-indexer-queue"
  }
}

# Dead Letter Queue for Notification Service
resource "aws_sqs_queue" "notification_dlq" {
  name = "${var.environment}-notification-dlq"

  tags = {
    Name = "${var.environment}-notification-dlq"
  }
}

# Dead Letter Queue for Search Indexer Service
resource "aws_sqs_queue" "search_indexer_dlq" {
  name = "${var.environment}-search-indexer-dlq"

  tags = {
    Name = "${var.environment}-search-indexer-dlq"
  }
}

# Configure DLQ for notification queue
resource "aws_sqs_queue_redrive_policy" "notification" {
  queue_url = aws_sqs_queue.notification.id

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.notification_dlq.arn
    maxReceiveCount     = 3
  })
}

# Configure DLQ for search indexer queue
resource "aws_sqs_queue_redrive_policy" "search_indexer" {
  queue_url = aws_sqs_queue.search_indexer.id

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.search_indexer_dlq.arn
    maxReceiveCount     = 3
  })
}

# EventBridge Rule: Route to Notification Queue
resource "aws_cloudwatch_event_rule" "notification" {
  name           = "${var.environment}-route-to-notification"
  description    = "Route events to notification service queue"
  event_bus_name = aws_cloudwatch_event_bus.main.name

  event_pattern = jsonencode({
    source = [
      "procurement-service",
      "logistics-service",
      "collaboration-service",
      "billing-service",
      "company-service",
      "catalog-service"
    ]
  })
}

resource "aws_cloudwatch_event_target" "notification" {
  rule      = aws_cloudwatch_event_rule.notification.name
  target_id = "NotificationQueue"
  arn       = aws_sqs_queue.notification.arn
  event_bus_name = aws_cloudwatch_event_bus.main.name
}

# EventBridge Rule: Route to Search Indexer Queue
resource "aws_cloudwatch_event_rule" "search_indexer" {
  name           = "${var.environment}-route-to-search-indexer"
  description    = "Route events to search indexer service queue"
  event_bus_name = aws_cloudwatch_event_bus.main.name

  event_pattern = jsonencode({
    source = [
      "catalog-service",
      "company-service",
      "procurement-service",
      "logistics-service"
    ]
  })
}

resource "aws_cloudwatch_event_target" "search_indexer" {
  rule      = aws_cloudwatch_event_rule.search_indexer.name
  target_id = "SearchIndexerQueue"
  arn       = aws_sqs_queue.search_indexer.arn
  event_bus_name = aws_cloudwatch_event_bus.main.name
}

# IAM Policy for EventBridge to send to SQS
resource "aws_iam_role_policy" "eventbridge_sqs" {
  name = "${var.environment}-eventbridge-sqs-policy"
  role = aws_iam_role.eventbridge.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "sqs:SendMessage"
        ]
        Resource = [
          aws_sqs_queue.notification.arn,
          aws_sqs_queue.search_indexer.arn
        ]
      }
    ]
  })
}

resource "aws_iam_role" "eventbridge" {
  name = "${var.environment}-eventbridge-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "events.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })
}
