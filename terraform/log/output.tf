output log_arn {
  value       = aws_cloudwatch_log_group.log.arn
  description = "The ARN of the CloudWatch Log Group"
}
