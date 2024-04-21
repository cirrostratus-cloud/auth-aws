terraform {
  source = "${get_parent_terragrunt_dir()}/terraform/log"
}

locals {
  common_vars  = read_terragrunt_config(find_in_parent_folders("common.hcl"))
  module_name = local.common_vars.locals.module_name
  log_name = "event-user"
  common_tags = local.common_vars.locals.common_tags
}

include {
  path = find_in_parent_folders()
}

inputs = {
  log_name = local.log_name
  module_name = local.module_name
  common_tags = local.common_tags
}
