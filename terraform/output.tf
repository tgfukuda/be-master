output "ecr_repository_id" {
  description = "ECR Repository ID"
  value       = aws_ecr_repository.ecr_repo.id
}

output "rds_instance_id" {
  description = "RDS Instance ID"
  value       = aws_db_instance.rds_instance.id
}

output "secret_manager_id" {
  description = "Secret Manager ID"
  value       = aws_secretsmanager_secret.secret_manager.id
}

output "ecr_publish_role_arn" {
  description = "ECR Publish Role ARN"
  value       = aws_iam_role.ecr_publish_role.arn
}
