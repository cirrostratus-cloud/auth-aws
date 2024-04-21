variable log_name {
  type        = string
  description = "Log group name"
}

variable module_name {
  type        = string
  description = "Module name"
}


variable retention_in_days {
  type        = number
  default     = 30
  description = "Log retention in days"
}


variable common_tags {
  type        = map(string)
  description = "Common tags for all resources"
}
