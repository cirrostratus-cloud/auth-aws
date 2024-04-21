terraform {
  source = "${get_parent_terragrunt_dir()}/terraform/function"
}

locals {
  common_vars  = read_terragrunt_config(find_in_parent_folders("common.hcl"))
  module_name = local.common_vars.locals.module_name
  function_name = "http-user"
  common_tags = local.common_vars.locals.common_tags
}

include {
  path = find_in_parent_folders()
}

dependency log {
    config_path = "${get_parent_terragrunt_dir()}/terragrunt/user/http/log"
    mock_outputs = {
        log_arn = "log_arn"
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
          "sns:Publish"
        ],
        "Resource": "${get_env("TOPIC_ARN_PREFIX")}*"
      }
    ]
  })
  environment_variables = {
    LOG_LEVEL = "INFO"
    AWS_STAGE = local.common_vars.locals.aws_stage
    CIRROSTRATUS_AUTH_MODULE_NAME = local.common_vars.locals.module_name
    CIRROSTRATUS_AUTH_USER_TABLE = get_env("CIRROSTRATUS_AUTH_USER_TABLE")
    USER_MIN_PASSWORD_LENGTH = get_env("USER_MIN_PASSWORD_LENGTH")
    USER_UPPER_CASE_REQUIRED = get_env("USER_UPPER_CASE_REQUIRED")
    USER_LOWER_CASE_REQUIRED = get_env("USER_LOWER_CASE_REQUIRED")
    USER_NUMBER_REQUIRED = get_env("USER_NUMBER_REQUIRED")
    USER_SPECIAL_CHARACTER_REQUIRED = get_env("USER_SPECIAL_CHARACTER_REQUIRED")
    TOPIC_ARN_PREFIX = get_env("TOPIC_ARN_PREFIX")
    PRIVATE_KEY = get_env("PRIVATE_KEY")
  }
  module_bucket = local.common_vars.locals.module_bucket
  file_location = "${get_parent_terragrunt_dir()}/bin/user/http"
  zip_location = "${get_parent_terragrunt_dir()}/dist/user/http"
  zip_name = "user.zip"
  common_tags = local.common_tags
}
