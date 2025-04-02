# Terraform Provider for Control Tower Management

This provider allows you to manage AWS Control Tower resources using Terraform. It provides data sources and resources to interact with AWS Control Tower and AWS Organizations.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0.0
- [Go](https://golang.org/doc/install) >= 1.21
- AWS credentials with appropriate permissions
- AWS Organizations access

## Building The Provider

1. Clone the repository
```bash
git clone https://github.com/yourusername/terraform-provider-controltowermanagement.git
cd terraform-provider-controltowermanagement
```

2. Build the provider
```bash
go build -o terraform-provider-controltowermanagement
```

3. Install the provider
```bash
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/yourusername/controltowermanagement/1.0.0/darwin_amd64
cp terraform-provider-controltowermanagement ~/.terraform.d/plugins/registry.terraform.io/yourusername/controltowermanagement/1.0.0/darwin_amd64/
```

## Using the provider

### Provider Configuration

```hcl
provider "controltowermanagement" {
  # AWS credentials can be provided via environment variables:
  # AWS_ACCESS_KEY
  # AWS_SECRET_ACCESS_KEY
  # AWS_REGION
  # Or directly in the configuration:
  access_key = "your-access-key"
  secret_key = "your-secret-key"
  region     = "us-west-2"

  # Optional: Assume role configuration
  assume_role {
    role_arn = "arn:aws:iam::1111111111:role/AWSControlTowerExecution"
  }
}
```

### Data Sources

#### AWS Account Data Source

Use this data source to get information about AWS accounts in your organization.

```hcl
data "controltowermanagement_aws_account" "example" {}

output "accounts" {
  value = data.controltowermanagement_aws_account.example.accounts
}
```

##### Attributes

| Name | Description | Type |
|------|-------------|------|
| accounts | List of AWS accounts | List of Object |
| accounts.account_id | The ID of the AWS account | String |
| accounts.account_name | The name of the AWS account | String |
| accounts.email | The email address associated with the account | String |
| accounts.status | The status of the account (ACTIVE, SUSPENDED, etc.) | String |

### Examples

See the [examples](examples) directory for more detailed examples of using the provider.

## Development

### Running Tests

1. Unit Tests
```bash
go test ./internal/client -v
```

2. Acceptance Tests
```bash
export AWS_ACCESS_KEY="your-access-key"
export AWS_SECRET_ACCESS_KEY="your-secret-key"
export AWS_REGION="your-region"
go test ./internal/provider -v
```

### Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 