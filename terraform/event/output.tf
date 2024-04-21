output topic_arn {
  value       = aws_sns_topic.event[*].arn
  description = "The ARN of the SNS topic"
}

output queue_arn {
  value       = aws_sqs_queue.event[*].arn
  description = "The ARN of the SQS queue"
}
