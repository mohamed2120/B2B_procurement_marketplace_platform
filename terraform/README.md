# Terraform Infrastructure as Code

This directory contains Terraform configurations for deploying the B2B Procurement Marketplace Platform to AWS staging environment.

## Architecture Overview

- **VPC**: 2 Availability Zones with public and private subnets
- **ECS Fargate**: Container orchestration for microservices
- **ALB**: Application Load Balancer with path-based routing
- **RDS PostgreSQL**: Managed database in private subnets
- **ElastiCache Redis**: Caching and pub/sub
- **OpenSearch**: Search and indexing
- **S3**: Document and media storage
- **Cognito**: User authentication
- **EventBridge + SQS**: Event-driven architecture

## Prerequisites

1. **AWS CLI** configured with appropriate credentials
2. **Terraform** >= 1.5.0 installed
3. **AWS Account** with appropriate permissions
4. **S3 Bucket** for Terraform state (configure in backend)
5. **ACM Certificate** for HTTPS (optional but recommended)
6. **ECR Repositories** created for each service (or update ECS task definitions after)

## Setup

### 1. Configure Backend

Edit `main.tf` to configure the S3 backend for Terraform state:

```hcl
backend "s3" {
  bucket = "your-terraform-state-bucket"
  key    = "b2b-platform/staging/terraform.tfstate"
  region = "us-east-1"
}
```

### 2. Create Terraform Variables File

Copy the example variables file:

```bash
cp terraform.tfvars.example terraform.tfvars
```

Edit `terraform.tfvars` with your values:

```hcl
aws_region  = "us-east-1"
environment = "staging"
vpc_cidr    = "10.0.0.0/16"

# Database credentials (use AWS Secrets Manager in production)
db_password = "your-secure-password"

# Domain and certificate
domain_name     = "staging.yourdomain.com"
certificate_arn = "arn:aws:acm:us-east-1:123456789012:certificate/xxxx"
```

### 3. Initialize Terraform

```bash
cd terraform
terraform init
```

### 4. Review Plan

```bash
terraform plan
```

This will show you what resources will be created. **Review carefully before applying.**

### 5. Apply Infrastructure

**⚠️ WARNING: This will create AWS resources and incur costs.**

```bash
terraform apply
```

Type `yes` when prompted to confirm.

## Module Structure

```
terraform/
├── main.tf                 # Main configuration
├── variables.tf            # Variable definitions
├── outputs.tf              # Output values
├── terraform.tfvars.example # Example variables
└── modules/
    ├── vpc/                # VPC, subnets, NAT gateways
    ├── security-groups/    # Security groups for all services
    ├── alb/                # Application Load Balancer
    ├── ecs/                # ECS cluster, services, task definitions
    ├── rds/                # RDS PostgreSQL
    ├── redis/              # ElastiCache Redis
    ├── opensearch/         # OpenSearch domain
    ├── s3/                 # S3 buckets
    ├── cognito/            # Cognito user pool
    └── eventbridge/        # EventBridge bus and SQS queues
```

## Key Resources Created

### Networking
- VPC with CIDR 10.0.0.0/16
- 2 Public subnets (one per AZ)
- 2 Private subnets (one per AZ)
- Internet Gateway
- 2 NAT Gateways (one per AZ)
- Route tables and associations

### Compute
- ECS Fargate cluster
- ECS services for each microservice
- Task definitions with IAM roles
- CloudWatch log groups

### Load Balancing
- Application Load Balancer
- Target groups (one per service)
- Listener rules for path-based routing
- HTTP to HTTPS redirect

### Database & Cache
- RDS PostgreSQL 15.4 instance
- ElastiCache Redis 7.0
- OpenSearch 2.11 domain

### Storage
- S3 bucket: `{env}-docs-private-{account_id}`
- S3 bucket: `{env}-media-{account_id}`

### Security
- Security groups for ALB, ECS, RDS, Redis, OpenSearch
- IAM roles for ECS tasks
- Cognito user pool and client

### Events
- EventBridge custom bus
- SQS queues for notification and search-indexer services
- Dead letter queues
- Event routing rules

## Outputs

After applying, Terraform will output:

- `alb_dns_name`: ALB DNS name for accessing services
- `rds_endpoint`: RDS connection endpoint
- `redis_endpoint`: Redis connection endpoint
- `opensearch_endpoint`: OpenSearch endpoint
- `s3_buckets`: S3 bucket names
- `cognito_user_pool_id`: Cognito user pool ID
- `eventbridge_bus_name`: EventBridge bus name

View outputs:

```bash
terraform output
```

## Service Configuration

Services are configured in `variables.tf` under the `services` variable. Each service has:

- Port number
- CPU and memory allocation
- Desired task count
- Path pattern for ALB routing
- Health check path

## Path-Based Routing

The ALB routes traffic based on path patterns:

- `/identity/*` → identity-service
- `/company/*` → company-service
- `/catalog/*` → catalog-service
- `/procurement/*` → procurement-service
- `/logistics/*` → logistics-service
- `/collaboration/*` → collaboration-service
- `/notification/*` → notification-service

## Environment Variables for ECS Tasks

ECS tasks receive environment variables for:

- Database connection (host, port, name)
- Redis connection
- OpenSearch URL
- Secrets from AWS Secrets Manager (DB password, JWT secret)

## Secrets Management

**Important**: In production, use AWS Secrets Manager for:

- Database passwords
- JWT secrets
- API keys
- Other sensitive credentials

Update task definitions to reference Secrets Manager ARNs.

## Cost Considerations

This configuration creates resources that incur AWS costs:

- NAT Gateways: ~$32/month each (2 = ~$64/month)
- RDS: Depends on instance class (~$50-200/month)
- ECS Fargate: Pay per vCPU/hour and GB/hour
- OpenSearch: Depends on instance type and count
- ElastiCache: Depends on node type
- ALB: ~$16/month + data transfer
- S3: Pay per GB stored and requests

**Estimated monthly cost for staging: $200-500** (varies by usage)

## Cleanup

To destroy all resources:

```bash
terraform destroy
```

**⚠️ WARNING: This will delete all resources. Ensure you have backups.**

## Next Steps

After infrastructure is created:

1. **Build and push Docker images** to ECR
2. **Update ECS task definitions** with actual image URIs
3. **Configure Secrets Manager** for sensitive values
4. **Set up CI/CD pipeline** for automated deployments
5. **Configure Route53** DNS records pointing to ALB
6. **Set up monitoring and alerts** in CloudWatch
7. **Configure backup policies** for RDS and S3

## Troubleshooting

### Terraform State Lock

If you see a state lock error:

```bash
terraform force-unlock <LOCK_ID>
```

### Module Not Found

Ensure all modules are in the `modules/` directory and run:

```bash
terraform init -upgrade
```

### Resource Already Exists

If a resource already exists, import it:

```bash
terraform import <resource_type>.<name> <resource_id>
```

## Security Best Practices

1. **Never commit** `terraform.tfvars` with real credentials
2. **Use AWS Secrets Manager** for all secrets
3. **Enable encryption** at rest for RDS, S3, ElastiCache
4. **Use private subnets** for all application resources
5. **Restrict security groups** to minimum required access
6. **Enable VPC Flow Logs** for network monitoring
7. **Use IAM roles** instead of access keys
8. **Enable CloudTrail** for audit logging

## Support

For issues or questions:
- Review Terraform documentation: https://www.terraform.io/docs
- Check AWS service documentation
- Review module-specific README files in `modules/` directories
