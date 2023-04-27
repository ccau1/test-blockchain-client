# test-blockchain-client

## Setup

set .env
```
cp .env.example .env.dev
```

## Start

### Docker Compose

```
docker-compose --env-file .env.dev up --build
```

### gomon

```
gomon .
```

## Terraform

1. setup aws account for terraform
2. set IAM user for terraform (will need to define policy based on terraform needs)
3. set IAM credential to `~/.aws/credentials` with profile name `blockchain-client-terraform`
4. run `terraform plan` to prepare changes to AWS
5. run `terraform apply` to apply prepared changes to AWS
6. copy `${env}_deploy_access_key_id` from output of step 5 to store in github action secrets `AWS_ACCESS_KEY_ID`
7. copy secret from the following command to store in github action secrets `AWS_SECRET_ACCESS_KEY`: `terraform state pull | jq '.resources[] | select(.type == "aws_iam_access_key") | .instances[0].attributes'`

## Production Ready Requirement

- add unit tests
- setup github actions to run test, build and deploy to AWS ECR
- terraform:
  - change all prefix "dev-" to "prod-". The purpose for the prefix is to have the option to run env (ie. dev, test, staging) on the same account.
  - use parameter store and secret manager to store env variables that can be populated during deployment phase. This way, variables are stored securely and instances will have the latest variables on service start.
  - move terraform to its own repo so it separates DevOps from developers
- manually set secrets into secrets manager and read secrets into env on deploy
- instead of retrieving IAM user's access key to set in github, can change to using OIDC

## TODO

[x] throw error responses for requests
[x] use validator
[ ] test cases
[ ] set terraform for deployment to ecs fargate, write HCL
[ ] set terraform for different env
[ ] handle rpc endpoints based on endpoints returned errors
[ ] handle rpc endpoints based on speed
[ ] add a proxy (separate service for each provider, single entry point for wallet side to communicate with)