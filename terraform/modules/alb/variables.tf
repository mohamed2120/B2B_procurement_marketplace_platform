variable "environment" {
  description = "Environment name"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID"
  type        = string
}

variable "public_subnet_ids" {
  description = "List of public subnet IDs"
  type        = list(string)
}

variable "certificate_arn" {
  description = "ARN of ACM certificate"
  type        = string
  default     = ""
}

variable "domain_name" {
  description = "Domain name"
  type        = string
  default     = ""
}

variable "services" {
  description = "Map of services with their configuration"
  type = map(object({
    port              = number
    path_pattern      = string
    health_check_path = string
    priority          = number
  }))
}
