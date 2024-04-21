terraform {
  source = "${get_parent_terragrunt_dir()}/terraform/event"
}

locals {
  common_vars = read_terragrunt_config(find_in_parent_folders("common.hcl"))
  module_name = local.common_vars.locals.module_name
  common_tags = local.common_vars.locals.common_tags
}

include {
  path = find_in_parent_folders()
}

inputs = {
  module_name = local.module_name
  topics = ["user_created", "user_password_changed", "user_password_recovered", "user_email_confirmed" ]
  delivery_policy = jsonencode({
    http = {
      defaultHealthyRetryPolicy = {
        minDelayTarget = 20,
        maxDelayTarget = 20,
        numRetries = 3,
        numMaxDelayRetries = 0,
        numNoDelayRetries = 0,
        numMinDelayRetries = 0,
        backoffFunction = "linear"
      },
      disableSubscriptionOverrides = false,
      defaultThrottlePolicy = {
        maxReceivesPerSecond = 1
      }
    }
  })
  common_tags = local.common_tags
}
