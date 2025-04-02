terraform {
  required_providers {
    controltowermanagement = {
      source = "registry.terraform.io/eaglespirittech/controltowermanagement"
    }
  }
}

provider "controltowermanagement" {
  # AWS credentials can be provided via environment variables:
  # AWS_ACCESS_KEY
  # AWS_SECRET_ACCESS_KEY
  # AWS_REGION
  # Or directly in the configuration:
  # access_key = "your-access-key"
  # secret_key = "your-secret-key"
  # region     = "us-west-2"

  # Optional: Assume role configuration
  # assume_role {
  #   role_arn = "arn:aws:iam::1111111111:role/AWSControlTowerExecution"
  # }
}

# List all AWS accounts in the organization
data "controltowermanagement_aws_account" "example" {}

# Output the account information
output "accounts" {
  description = "List of AWS accounts in the organization"
  value = data.controltowermanagement_aws_account.example.accounts
}

# Example of filtering accounts by status
output "active_accounts" {
  description = "List of active AWS accounts"
  value = [
    for account in data.controltowermanagement_aws_account.example.accounts :
    account if account.status == "ACTIVE"
  ]
}

# Example of finding a specific account by email
output "specific_account" {
  description = "Details of a specific account"
  value = [
    for account in data.controltowermanagement_aws_account.example.accounts :
    account if account.email == "admin@example.com"
  ]
}