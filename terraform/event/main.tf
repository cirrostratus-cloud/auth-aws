terraform {
  backend "s3" {}
}

resource "aws_sns_topic" "event" {
  count = length(var.topics)
  name = "${var.module_name}_${var.topics[count.index]}"
  delivery_policy = var.delivery_policy
  tags = var.common_tags
}

resource "aws_sqs_queue" "event_dlq" {
  count = length(var.topics)
  name = "${var.module_name}_${var.topics[count.index]}-dlq"
  tags = var.common_tags
}

resource "aws_sqs_queue" "event" {
  count = length(var.topics)
  name = "${var.module_name}_${var.topics[count.index]}"
  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.event_dlq[count.index].arn
    maxReceiveCount = 3
  })
  tags = var.common_tags
}

resource "aws_sns_topic_subscription" "event" {
  count = length(var.topics)
  topic_arn = aws_sns_topic.event[count.index].arn
  protocol  = "sqs"
  endpoint  = aws_sqs_queue.event[count.index].arn
}

resource "aws_sqs_queue_policy" "event" {
  count = length(var.topics)
  queue_url = aws_sqs_queue.event[count.index].id
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Sid = "Allow-SNS-Event",
        Effect = "Allow",
        Principal = "*",
        Action = "sqs:SendMessage",
        Resource = aws_sqs_queue.event[count.index].arn,
        Condition = {
          ArnEquals = {
            "aws:SourceArn" = aws_sns_topic.event[count.index].arn
          }
        }
      }
    ]
  })
}