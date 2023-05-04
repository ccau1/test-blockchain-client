# test-blockchain-client

## Table of Content

- [Ideology](##ideology)
- [Setup](##setup)
- [Start Application](##start-application)
- [API Calls](##api-calls)
- [Terraform](##terraform)
- [Production Ready Requirement](#production-ready-requirement)
- [TODO](##todo)

## Ideology

### The Problem
We want a proxy where different wallet chains can go through this service to trigger actions or fetch information. Items to concern includes:

- manage multiple providers per chain type
- manage navigating through providers that cannot connect
- manage selecting provider based on response time
- reduce expenditure on provider usage
- scalable services for each chain type

### Solution

![Proxy Diagram](/proxy-diagram.png)

There are two main layers: Providers and ProviderAccounts

Providers holds a list of different Providers, all of which defines the same methods (ie. GetBlockNumber, GetBlockByNumber). Each provider define how to communicate with its corresponding provider endpoint.

ProviderAccounts holds a list of provider's accounts. This way, we can use multiple accounts per provider, based on our needs (ie. rotate between them on reaching quota limit).

Both Providers and ProviderAccounts can have a layer of strategy that defines how to select which Provider/ProviderAccount to use.

A proxy helps centralize all chain route methods

## Setup

copy .env.example to your needed env `.env.{environment}`
```
cp .env.example .env.dev
```

## Start Application

### Docker Compose

```
docker-compose --env-file .env.dev up --build
```

### gomon

```
gomon .
```

## API Calls

### Block Number
fetch latest block number

#### Queries
`batchId` (optional) - a unique number in case of batching

`jsonrpc` (optional) - json rpc version
```
[GET] localhost:3000/eth/{blockNumber}?batchId=3&jsonrpc=2.0
```

### Block by Number
fetch block number by number

#### Params
`blockNumber` - block number to fetch

#### Queries
`batchId` (optional) - a unique number in case of batching

`jsonrpc` (optional) - json rpc version
```
[GET] localhost:3000/eth/block-number/blockNumber?batchId=3&jsonrpc=2.0
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
- github actions:
  - run test (`go test ./...`) and promote to staging (in this scenario there's dev and staging)
  - cache installed steps
  - add sonarcloud for extra layer to code checking
- terraform:
  - change all prefix "dev-" to "prod-". The purpose for the prefix is to have the option to run env (ie. dev, test, staging) on the same account.
  - use parameter store and secret manager to store env variables that can be populated during deployment phase. This way, variables are stored securely and instances will have the latest variables on service start.
  - move terraform to its own repo so it separates DevOps from developers
- manually set secrets into secrets manager and read secrets into env on deploy
- instead of retrieving IAM user's access key to set in github, can change to using OIDC
- move some of the accounts strategy data to using MemCache so multiple instances can decide on which request to call based on single source of truth

## TODO

[x] throw error responses for requests

[x] use validator

[ ] test cases

[x] set terraform for deployment to ecs fargate, write HCL

[x] set terraform for different env

[x] use multiple accounts under a rpc endpoint

[ ] handle rotating providers based on endpoints returned errors

[ ] handle rotating providers based on speed

[x] add a proxy (separate service for each chain type, single entry point for wallet side to communicate with all chain types)