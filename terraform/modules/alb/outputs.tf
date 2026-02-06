output "dns_name" {
  description = "ALB DNS name"
  value       = aws_lb.main.dns_name
}

output "zone_id" {
  description = "ALB zone ID"
  value       = aws_lb.main.zone_id
}

output "arn" {
  description = "ALB ARN"
  value       = aws_lb.main.arn
}

output "security_group_id" {
  description = "ALB security group ID"
  value       = aws_security_group.alb.id
}

output "target_group_arns" {
  description = "Map of target group ARNs by service name"
  value       = { for k, v in aws_lb_target_group.services : k => v.arn }
}
