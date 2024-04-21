variable function_name {
  type = string
  description = "Name of lambda function"
}

variable module_name {
    type = string
    description = "Project name"
}

variable iam_policy {
    type = string
    description = "IAM policy"
}

variable memory_size {
    type = number
    default = 128
    description = "Memory size"
}

variable timeout {
    type = number
    default = 30
    description = "Timeout"
}

variable environment_variables {
    type = map(string)
    description = "Environment variables"
}

variable module_bucket {
    type = string
    description = "Module bucket"
}

variable file_location {
    type = string
    description = "File location"
}

variable zip_location {
    type = string
    description = "Zip location"
}

variable zip_name {
    type = string
    description = "Zip name"
}

variable common_tags {
    type = map(string)
    description = "Common tags"
}

variable event_sources_arn {
    type = list(string)
    default = []
    description = "Event sources arn"
}

