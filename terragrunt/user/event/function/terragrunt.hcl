terraform {
  source = "${get_parent_terragrunt_dir()}/terraform/function"
}

locals {
  common_vars  = read_terragrunt_config(find_in_parent_folders("common.hcl"))
  module_name = local.common_vars.locals.module_name
  function_name = "event-user"
  common_tags = local.common_vars.locals.common_tags
}

include {
  path = find_in_parent_folders()
}

dependency log {
    config_path = "${get_parent_terragrunt_dir()}/terragrunt/user/event/log"
    mock_outputs = {
        log_arn = "log_arn"
    }
}

dependency subscription {
    config_path = "${get_parent_terragrunt_dir()}/terragrunt/user/event/subscription"
    mock_outputs = {
        topic_arn = ["topic_arn"]
        queue_arn = ["queue_arn"]
    }
}

inputs = {
  function_name = local.function_name
  module_name = local.module_name
  iam_policy = jsonencode({
    "Version": "2012-10-17",
    "Statement": [
      {
        "Effect": "Allow",
        "Action": [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ],
        "Resource": "${dependency.log.outputs.log_arn}:*"
      },
      {
        "Effect": "Allow",
        "Action": [
          "sqs:ChangeMessageVisibility",
          "sqs:DeleteMessage",
          "sqs:GetQueueAttributes",
          "sqs:ReceiveMessage",
        ],
        "Resource": dependency.subscription.outputs.queue_arn
      },
      {
        "Effect": "Allow",
        "Action": [
          "dynamodb:PutItem",
          "dynamodb:GetItem",
          "dynamodb:UpdateItem",
          "dynamodb:DeleteItem",
          "dynamodb:Query",
          "dynamodb:Scan"
        ],
        "Resource": "arn:aws:dynamodb:*:*:table/${local.common_vars.locals.module_name}-${get_env("CIRROSTRATUS_AUTH_USER_TABLE")}"
      },
      {
        "Effect": "Allow",
        "Action": [
          "ses:SendEmail",
        ],
        "Resource": [
          get_env("SES_EMAIl_ARN"),
          get_env("SES_CONFIGURATION_SET"),
        ]
      }
    ]
  })
  environment_variables = {
    LOG_LEVEL = "INFO"
    AWS_STAGE = local.common_vars.locals.aws_stage
    CIRROSTRATUS_AUTH_MODULE_NAME = local.common_vars.locals.module_name
    CIRROSTRATUS_AUTH_USER_TABLE = get_env("CIRROSTRATUS_AUTH_USER_TABLE")
    SES_EMAIL_FROM = get_env("SES_EMAIL_FROM")
    TOPIC_ARN_PREFIX = get_env("TOPIC_ARN_PREFIX")
    EMAIL_CONFIRMATION_URL = get_env("EMAIL_CONFIRMATION_URL")
    PRIVATE_KEY = get_env("PRIVATE_KEY")
    MAX_AGE_IN_SECONDS = get_env("MAX_AGE_IN_SECONDS")
  }
  module_bucket = local.common_vars.locals.module_bucket
  file_location = "${get_parent_terragrunt_dir()}/bin/user/subscriber"
  zip_location = "${get_parent_terragrunt_dir()}/dist/user/subscriber"
  zip_name = "user.zip"
  common_tags = local.common_tags
  event_sources_arn = dependency.subscription.outputs.queue_arn
}
