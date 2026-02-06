variable "vpc_id" {
  description = "VPC ID"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
}

variable "alb_security_group_id" {
  description = "ALB security group ID"
  type        = string
}
