variable "aws_region" {
  description = "AWS region for resources"
  type        = string
  default     = "us-east-1"
}

variable "environment" {
  description = "Environment name (staging, production)"
  type        = string
  default     = "staging"
}

variable "vpc_cidr" {
  description = "CIDR block for VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "certificate_arn" {
  description = "ARN of ACM certificate for ALB"
  type        = string
  default     = ""
}

variable "domain_name" {
  description = "Domain name for the application"
  type        = string
  default     = "b2b-platform.example.com"
}

# RDS Variables
variable "db_instance_class" {
  description = "RDS instance class"
  type        = string
  default     = "db.t3.medium"
}

variable "db_allocated_storage" {
  description = "RDS allocated storage in GB"
  type        = number
  default     = 100
}

variable "db_name" {
  description = "RDS database name"
  type        = string
  default     = "b2b_platform"
}

variable "db_username" {
  description = "RDS master username"
  type        = string
  default     = "b2b_admin"
  sensitive   = true
}

variable "db_password" {
  description = "RDS master password"
  type        = string
  sensitive   = true
}

variable "db_backup_retention" {
  description = "RDS backup retention period in days"
  type        = number
  default     = 7
}

# Redis Variables
variable "redis_node_type" {
  description = "ElastiCache Redis node type"
  type        = string
  default     = "cache.t3.micro"
}

variable "redis_num_cache_nodes" {
  description = "Number of cache nodes"
  type        = number
  default     = 1
}

# OpenSearch Variables
variable "opensearch_instance_type" {
  description = "OpenSearch instance type"
  type        = string
  default     = "t3.small.search"
}

variable "opensearch_instance_count" {
  description = "Number of OpenSearch instances"
  type        = number
  default     = 2
}

variable "opensearch_volume_size" {
  description = "OpenSearch EBS volume size in GB"
  type        = number
  default     = 20
}

# Logging
variable "log_retention_days" {
  description = "CloudWatch log retention in days"
  type        = number
  default     = 7
}

# Service Configuration
variable "services" {
  description = "List of microservices to deploy"
  type = map(object({
    port            = number
    cpu             = number
    memory          = number
    desired_count   = number
    path_pattern    = string
    health_check_path = string
  }))
  default = {
    identity-service = {
      port              = 8001
      cpu               = 256
      memory            = 512
      desired_count     = 2
      path_pattern      = "/identity/*"
      health_check_path = "/health"
    }
    company-service = {
      port              = 8002
      cpu               = 256
      memory            = 512
      desired_count     = 2
      path_pattern      = "/company/*"
      health_check_path = "/health"
    }
    catalog-service = {
      port              = 8003
      cpu               = 256
      memory            = 512
      desired_count     = 2
      path_pattern      = "/catalog/*"
      health_check_path = "/health"
    }
    procurement-service = {
      port              = 8006
      cpu               = 512
      memory            = 1024
      desired_count     = 2
      path_pattern      = "/procurement/*"
      health_check_path = "/health"
    }
    logistics-service = {
      port              = 8007
      cpu               = 256
      memory            = 512
      desired_count     = 2
      path_pattern      = "/logistics/*"
      health_check_path = "/health"
    }
    collaboration-service = {
      port              = 8008
      cpu               = 256
      memory            = 512
      desired_count     = 2
      path_pattern      = "/collaboration/*"
      health_check_path = "/health"
    }
    notification-service = {
      port              = 8009
      cpu               = 256
      memory            = 512
      desired_count     = 2
      path_pattern      = "/notification/*"
      health_check_path = "/health"
    }
  }
}
