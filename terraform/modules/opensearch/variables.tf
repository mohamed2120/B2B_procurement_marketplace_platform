variable "environment" {
  description = "Environment name"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID"
  type        = string
}

variable "private_subnet_ids" {
  description = "List of private subnet IDs"
  type        = list(string)
}

variable "security_group_id" {
  description = "Security group ID for OpenSearch"
  type        = string
}

variable "instance_type" {
  description = "OpenSearch instance type"
  type        = string
  default     = "t3.small.search"
}

variable "instance_count" {
  description = "Number of OpenSearch instances"
  type        = number
  default     = 2
}

variable "volume_size" {
  description = "EBS volume size in GB"
  type        = number
  default     = 20
}

variable "master_user_arn" {
  description = "ARN of master user for OpenSearch"
  type        = string
  default     = ""
}
