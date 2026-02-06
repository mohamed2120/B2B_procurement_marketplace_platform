variable "environment" {
  description = "Environment name"
  type        = string
}

variable "domain_name" {
  description = "Domain name for callbacks"
  type        = string
  default     = ""
}
