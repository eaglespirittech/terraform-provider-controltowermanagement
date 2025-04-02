# AWS Account Data Source Example

This example demonstrates how to use the `controltowermanagement_aws_account` data source to retrieve information about AWS accounts in your organization.

## Usage

To run this example:

1. Make sure you have the provider installed:
```bash
terraform init
```

2. Configure your AWS credentials either through environment variables or in the provider block:
```bash
export AWS_ACCESS_KEY="your-access-key"
export AWS_SECRET_ACCESS_KEY="your-secret-key"
export AWS_REGION="us-west-2"
```

3. Run the example:
```bash
terraform plan
```

## Features Demonstrated

- Basic data source usage
- Provider configuration with assume role support
- Environment variable support for credentials
- Output examples showing different ways to use the account data:
  - List all accounts
  - Filter active accounts
  - Find specific account by email

## Outputs

The example will output:
- `accounts`: List of all AWS accounts in the organization
- `active_accounts`: List of only the active AWS accounts
- `specific_account`: Details of a specific account (if found)

## Requirements

- Terraform >= 1.0.0
- AWS credentials with appropriate permissions
- AWS Organizations access 