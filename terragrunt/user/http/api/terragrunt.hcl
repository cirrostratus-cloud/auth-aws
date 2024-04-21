terraform {
  source = "${get_parent_terragrunt_dir()}/terraform/api"
}

locals {
  common_vars = read_terragrunt_config(find_in_parent_folders("common.hcl"))
  aws_stage = local.common_vars.locals.aws_stage
  common_tags = local.common_vars.locals.common_tags
}

include {
  path = find_in_parent_folders()
}

dependency function {
    config_path = "${get_parent_terragrunt_dir()}/terragrunt/user/http/function"
    mock_outputs = {
        invoke_arn = "invoke_arn"
        function_name = "function_name"
        lambda_arn = "lambda_arn"
    }
}

inputs = {
  invoke_arn = dependency.function.outputs.invoke_arn
  function_name = dependency.function.outputs.function_name
  stage_name = local.aws_stage
  routes = [
    {
        path = "/{proxy+}"
        method = "ANY"
    }
  ]
  common_tags = local.common_tags
}
