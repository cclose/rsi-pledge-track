# Future deployment (IaC + CI/CD)

This repo currently deploys to AWS (S3/Lambda/API Gateway) but does not include infrastructure-as-code or an automated deployment pipeline.

This doc outlines a recommended approach to add auto-deploys from `main`.

## Goals

- Version-controlled infrastructure
- Repeatable deployments
- Minimal manual console work
- Safe secret handling

## Recommended direction

## Infrastructure as Code (pick one)

- AWS CDK (Go/TypeScript)
- Terraform
- AWS SAM / CloudFormation

Whichever you choose should define:

- API Gateway
- Lambda for `rsiAPI`
- EventBridge schedule + Lambda for `rsiPullFunding`
- RDS connectivity wiring (VPC subnets, security groups)
- IAM roles/policies

## Packaging model

- Build Go binaries for Linux (Lambda runtime).
- Package as:
  - zip artifacts for Lambda, or
  - container images (ECR) if preferred.

## CI/CD (GitHub Actions)

Typical pipeline steps:

- Run tests (`go test ./...`)
- Build deployment artifacts
- Deploy (assume-role to AWS via OIDC)

## Secrets

Current state uses plain Lambda env vars for DB config.

Recommended follow-up:

- Move DB credentials to **AWS Secrets Manager**.
- Optionally use RDS-managed secrets and rotation.

## S3 usage

If you host any static assets (or eventually a small status page), define an S3 bucket plus CloudFront in IaC.
