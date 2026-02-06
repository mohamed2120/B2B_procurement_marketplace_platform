terraform {
  required_version = ">= 1.5.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.0"
    }
  }

  backend "s3" {
    # Configure backend in terraform.tfvars or via CLI
    # bucket = "your-terraform-state-bucket"
    # key    = "b2b-platform/staging/terraform.tfstate"
    # region = "us-east-1"
  }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Environment = var.environment
      Project     = "b2b-procurement-platform"
      ManagedBy   = "terraform"
    }
  }
}

provider "random" {}

# Data sources
data "aws_availability_zones" "available" {
  state = "available"
}

data "aws_caller_identity" "current" {}

# VPC Module
module "vpc" {
  source = "./modules/vpc"

  environment             = var.environment
  vpc_cidr                = var.vpc_cidr
  availability_zones      = slice(data.aws_availability_zones.available.names, 0, 2)
  enable_nat_gateway      = true
  enable_vpn_gateway      = false
  enable_dns_hostnames    = true
  enable_dns_support      = true
}

# Security Groups
module "security_groups" {
  source = "./modules/security-groups"

  vpc_id                  = module.vpc.vpc_id
  environment             = var.environment
  alb_security_group_id   = module.alb.security_group_id
}

# Application Load Balancer
module "alb" {
  source = "./modules/alb"

  environment        = var.environment
  vpc_id             = module.vpc.vpc_id
  public_subnet_ids  = module.vpc.public_subnet_ids
  certificate_arn    = var.certificate_arn
  domain_name        = var.domain_name
  services           = {
    for k, v in var.services : k => {
      port              = v.port
      path_pattern      = v.path_pattern
      health_check_path = v.health_check_path
      priority          = index(keys(var.services), k) + 1
    }
  }
}

# ECS Cluster
module "ecs" {
  source = "./modules/ecs"

  environment           = var.environment
  vpc_id                = module.vpc.vpc_id
  private_subnet_ids    = module.vpc.private_subnet_ids
  alb_target_group_arns = module.alb.target_group_arns
  alb_security_group_id = module.alb.security_group_id
  ecs_security_group_id = module.security_groups.ecs_security_group_id
  log_group_name        = aws_cloudwatch_log_group.ecs.name
  services              = var.services
  db_host               = module.rds.address
  db_name               = var.db_name
  db_password_secret_arn = aws_secretsmanager_secret.db_password.arn
  jwt_secret_arn        = aws_secretsmanager_secret.jwt_secret.arn
  redis_host            = module.redis.endpoint
  opensearch_url        = "https://${module.opensearch.endpoint}"
  aws_region            = var.aws_region

  depends_on = [
    aws_cloudwatch_log_group.ecs,
    module.alb,
    module.rds,
    module.redis,
    module.opensearch,
    aws_secretsmanager_secret.db_password,
    aws_secretsmanager_secret.jwt_secret
  ]
}

# RDS PostgreSQL
module "rds" {
  source = "./modules/rds"

  environment         = var.environment
  vpc_id              = module.vpc.vpc_id
  private_subnet_ids  = module.vpc.private_subnet_ids
  security_group_id   = module.security_groups.rds_security_group_id
  db_instance_class   = var.db_instance_class
  db_allocated_storage = var.db_allocated_storage
  db_name             = var.db_name
  db_username         = var.db_username
  db_password         = var.db_password
  backup_retention    = var.db_backup_retention
}

# ElastiCache Redis
module "redis" {
  source = "./modules/redis"

  environment         = var.environment
  vpc_id              = module.vpc.vpc_id
  private_subnet_ids  = module.vpc.private_subnet_ids
  security_group_id   = module.security_groups.redis_security_group_id
  node_type           = var.redis_node_type
  num_cache_nodes     = var.redis_num_cache_nodes
}

# OpenSearch
module "opensearch" {
  source = "./modules/opensearch"

  environment         = var.environment
  vpc_id              = module.vpc.vpc_id
  private_subnet_ids  = module.vpc.private_subnet_ids
  security_group_id   = module.security_groups.opensearch_security_group_id
  instance_type       = var.opensearch_instance_type
  instance_count      = var.opensearch_instance_count
  volume_size        = var.opensearch_volume_size
}

# S3 Buckets
module "s3" {
  source = "./modules/s3"

  environment = var.environment
  account_id  = data.aws_caller_identity.current.account_id
}

# Cognito
module "cognito" {
  source = "./modules/cognito"

  environment = var.environment
  domain_name = var.domain_name
}

# EventBridge
module "eventbridge" {
  source = "./modules/eventbridge"

  environment = var.environment
}

# CloudWatch Log Group for ECS
resource "aws_cloudwatch_log_group" "ecs" {
  name              = "/ecs/${var.environment}/b2b-platform"
  retention_in_days = var.log_retention_days

  tags = {
    Name = "ecs-logs-${var.environment}"
  }
}

# Secrets Manager for Database Password
resource "aws_secretsmanager_secret" "db_password" {
  name = "${var.environment}/rds/db-password"

  tags = {
    Name = "${var.environment}-db-password-secret"
  }
}

resource "aws_secretsmanager_secret_version" "db_password" {
  secret_id     = aws_secretsmanager_secret.db_password.id
  secret_string = var.db_password
}

# Secrets Manager for JWT Secret
resource "aws_secretsmanager_secret" "jwt_secret" {
  name = "${var.environment}/app/jwt-secret"

  tags = {
    Name = "${var.environment}-jwt-secret"
  }
}

resource "aws_secretsmanager_secret_version" "jwt_secret" {
  secret_id     = aws_secretsmanager_secret.jwt_secret.id
  secret_string = "CHANGE_ME_IN_PRODUCTION_${random_password.jwt_secret.result}"
}

resource "random_password" "jwt_secret" {
  length  = 32
  special = true
}
