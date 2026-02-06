# DB Subnet Group
resource "aws_db_subnet_group" "main" {
  name       = "${var.environment}-db-subnet-group"
  subnet_ids = var.private_subnet_ids

  tags = {
    Name = "${var.environment}-db-subnet-group"
  }
}

# RDS Parameter Group
resource "aws_db_parameter_group" "main" {
  name   = "${var.environment}-postgres15"
  family = "postgres15"

  parameter {
    name  = "log_connections"
    value = "1"
  }

  parameter {
    name  = "log_statement"
    value = "all"
  }

  tags = {
    Name = "${var.environment}-postgres15-params"
  }
}

# RDS Instance
resource "aws_db_instance" "main" {
  identifier             = "${var.environment}-postgres"
  engine                 = "postgres"
  engine_version         = "15.4"
  instance_class         = var.db_instance_class
  allocated_storage      = var.db_allocated_storage
  max_allocated_storage  = var.db_allocated_storage * 2
  storage_type           = "gp3"
  storage_encrypted       = true

  db_name  = var.db_name
  username = var.db_username
  password = var.db_password

  db_subnet_group_name   = aws_db_subnet_group.main.name
  vpc_security_group_ids = [var.security_group_id]
  parameter_group_name   = aws_db_parameter_group.main.name

  backup_retention_period = var.backup_retention
  backup_window          = "03:00-04:00"
  maintenance_window     = "mon:04:00-mon:05:00"

  skip_final_snapshot       = false
  final_snapshot_identifier = "${var.environment}-postgres-final-snapshot-${formatdate("YYYY-MM-DD-hhmm", timestamp())}"
  deletion_protection       = false

  enabled_cloudwatch_logs_exports = ["postgresql", "upgrade"]
  performance_insights_enabled   = true

  tags = {
    Name = "${var.environment}-postgres"
  }
}
