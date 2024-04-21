terraform {
  backend "s3" {}
}

resource "aws_cloudwatch_log_group" "log" {
  name = "/aws/lambda/${var.module_name}-${var.log_name}"
  retention_in_days = var.retention_in_days
  tags = var.common_tags
}