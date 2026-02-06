output "endpoint" {
  description = "OpenSearch endpoint"
  value       = aws_opensearch_domain.main.endpoint
}

output "domain_id" {
  description = "OpenSearch domain ID"
  value       = aws_opensearch_domain.main.domain_id
}

output "arn" {
  description = "OpenSearch domain ARN"
  value       = aws_opensearch_domain.main.arn
}
