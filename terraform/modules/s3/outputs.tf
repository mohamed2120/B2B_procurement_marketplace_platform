output "docs_private_bucket" {
  description = "S3 bucket name for private documents"
  value       = aws_s3_bucket.docs_private.id
}

output "media_bucket" {
  description = "S3 bucket name for media"
  value       = aws_s3_bucket.media.id
}

output "docs_private_bucket_arn" {
  description = "S3 bucket ARN for private documents"
  value       = aws_s3_bucket.docs_private.arn
}

output "media_bucket_arn" {
  description = "S3 bucket ARN for media"
  value       = aws_s3_bucket.media.arn
}
