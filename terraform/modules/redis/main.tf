# ElastiCache Subnet Group
resource "aws_elasticache_subnet_group" "main" {
  name       = "${var.environment}-redis-subnet-group"
  subnet_ids = var.private_subnet_ids

  tags = {
    Name = "${var.environment}-redis-subnet-group"
  }
}

# ElastiCache Parameter Group
resource "aws_elasticache_parameter_group" "main" {
  name   = "${var.environment}-redis7"
  family = "redis7"

  parameter {
    name  = "maxmemory-policy"
    value = "allkeys-lru"
  }

  tags = {
    Name = "${var.environment}-redis7-params"
  }
}

# ElastiCache Replication Group (Redis)
resource "aws_elasticache_replication_group" "main" {
  replication_group_id       = "${var.environment}-redis"
  description                = "Redis cluster for ${var.environment}"

  engine               = "redis"
  engine_version       = "7.0"
  node_type            = var.node_type
  num_cache_nodes      = var.num_cache_nodes
  port                 = 6379
  parameter_group_name = aws_elasticache_parameter_group.main.name
  subnet_group_name    = aws_elasticache_subnet_group.main.name
  security_group_ids   = [var.security_group_id]

  at_rest_encryption_enabled = true
  transit_encryption_enabled = true
  auth_token_enabled         = false # Set to true and provide auth_token for production

  automatic_failover_enabled = var.num_cache_nodes > 1
  multi_az_enabled          = var.num_cache_nodes > 1

  snapshot_retention_limit = 3
  snapshot_window          = "03:00-05:00"

  tags = {
    Name = "${var.environment}-redis"
  }
}
