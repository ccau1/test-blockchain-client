module "blockchain-client-dev" {
  source = "./modules/blockchain-client"

  env = "dev"
}

module "blockchain-client-staging" {
  source = "./modules/blockchain-client"

  env = "staging"
}