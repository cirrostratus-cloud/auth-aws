locals {
  common_vars    = read_terragrunt_config(find_in_parent_folders("common.hcl"))
  aws_region     = local.common_vars.locals.aws_region
  module_name   = local.common_vars.locals.module_name
  module_bucket = local.common_vars.locals.module_bucket
}

generate "provider" {
  path      = "provider.tf"
  if_exists = "overwrite_terragrunt"
  contents  = <<EOF
  provider "aws" {
    region = "${local.aws_region}"
  }
EOF
}

remote_state {
  backend = "s3"
  config = {
    bucket = local.module_bucket
    key = "${path_relative_to_include()}/terraform.tfstate"
    region = "${local.aws_region}"
    encrypt  = true
    dynamodb_table = "${local.module_name}-${local.aws_region}-tfstate-lock"
    skip_region_validation = true
    skip_credentials_validation = true
    skip_metadata_api_check = true
    skip_requesting_account_id = true
    skip_s3_checksum = true
  }
}

inputs = merge(local.common_vars.locals)