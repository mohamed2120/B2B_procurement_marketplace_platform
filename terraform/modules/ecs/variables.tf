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

variable "alb_target_group_arns" {
  description = "Map of target group ARNs by service name"
  type        = map(string)
}

variable "alb_security_group_id" {
  description = "ALB security group ID"
  type        = string
}

variable "ecs_security_group_id" {
  description = "ECS security group ID"
  type        = string
}

variable "log_group_name" {
  description = "CloudWatch log group name"
  type        = string
}

variable "services" {
  description = "Map of services with their configuration"
  type = map(object({
    port              = number
    cpu               = number
    memory            = number
    desired_count     = number
    health_check_path = string
  }))
}

variable "ecr_repository_url" {
  description = "ECR repository URL base"
  type        = string
  default     = ""
}

variable "db_host" {
  description = "RDS database host"
  type        = string
}

variable "db_name" {
  description = "RDS database name"
  type        = string
}

variable "db_password_secret_arn" {
  description = "ARN of Secrets Manager secret for DB password"
  type        = string
}

variable "jwt_secret_arn" {
  description = "ARN of Secrets Manager secret for JWT secret"
  type        = string
}

variable "redis_host" {
  description = "Redis host"
  type        = string
}

variable "opensearch_url" {
  description = "OpenSearch URL"
  type        = string
}

variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}
