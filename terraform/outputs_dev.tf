output "dev_load_balancer_ip" {
  description = "external IP"
  value       = module.blockchain-client-dev.load_balancer_ip
}

output "dev_deploy_user_id" {
  description = "The deploy user's id"
  value       = module.blockchain-client-dev.deploy_user_id
}

output "dev_iam_user_name" {
  description = "The deploy user's name"
  value       = module.blockchain-client-dev.iam_user_name
}

output "dev_deploy_user_arn" {
  description = "The deploy user's ARN"
  value       = module.blockchain-client-dev.deploy_user_arn
}

output "dev_deploy_access_key_id" {
  description = "The deploy user's access key id"
  value       = module.blockchain-client-dev.deploy_access_key_id
}

output "dev_deploy_access_key_secret" {
  description = "The deploy user's access key secret"
  value       = module.blockchain-client-dev.deploy_access_key_secret
}

output "dev_deploy_access_key_key_fingerprint" {
  description = "The fingerprint of the PGP key used to encrypt the secret"
  value       = module.blockchain-client-dev.deploy_access_key_key_fingerprint
}

output "dev_secrets_manager_name" {
  description = "The secrets manager name"
  value       = module.blockchain-client-dev.secrets_manager_name
}