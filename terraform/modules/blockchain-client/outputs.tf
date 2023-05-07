output "load_balancer_ip" {
  description = "external IP"
  value       = aws_lb.default.dns_name
}

output "deploy_user_id" {
  description = "The deploy user's id"
  value       = aws_iam_user.deploy.unique_id
}

output "iam_user_name" {
  description = "The deploy user's name"
  value       = aws_iam_user.deploy.name
}

output "deploy_user_arn" {
  description = "The deploy user's ARN"
  value       = aws_iam_user.deploy.arn
}

output "deploy_access_key_id" {
  description = "The deploy user's access key id"
  value       = aws_iam_access_key.deploy.id
}

output "deploy_access_key_secret" {
  description = "The deploy user's access key secret"
  value       = aws_iam_access_key.deploy.encrypted_secret
}

output "deploy_access_key_key_fingerprint" {
  description = "The fingerprint of the PGP key used to encrypt the secret"
  value       = aws_iam_access_key.deploy.key_fingerprint
}

output "secrets_manager_name" {
  description = "The secrets manager name"
  value       = aws_secretsmanager_secret.deploy_access_key.name
}