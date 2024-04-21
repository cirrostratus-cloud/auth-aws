variable module_name {
  type        = string
  description = "The name of the module"
}

variable topics {
  type        = list(string)
  description = "The name of the SNS topics"
}

variable delivery_policy {
  type        = string
  description = "The delivery policy for the SNS topic"
}

variable common_tags {
  type        = map(string)
  description = "The common tags for the SNS topic"
}