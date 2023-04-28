module "blockchain-client-dev" {
  source = "./modules/blockchain-client"

  env    = "dev"
  region = var.region
}

module "blockchain-client-staging" {
  source = "./modules/blockchain-client"

  env    = "staging"
  region = var.region
}