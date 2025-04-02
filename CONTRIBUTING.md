# Contributing to Terraform Provider for Control Tower Management

Thank you for your interest in contributing to this project! This document provides guidelines and steps for contributing.

## Development Setup

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/yourusername/terraform-provider-controltowermanagement.git
   cd terraform-provider-controltowermanagement
   ```
3. Create a new branch for your feature:
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Development Guidelines

### Code Style

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` to format your code
- Run `golangci-lint` to check for linting issues

### Testing

1. Run unit tests:
   ```bash
   go test ./internal/client -v
   ```

2. Run acceptance tests:
   ```bash
   export AWS_ACCESS_KEY="your-access-key"
   export AWS_SECRET_ACCESS_KEY="your-secret-key"
   export AWS_REGION="your-region"
   go test ./internal/provider -v
   ```

### Documentation

- Update the README.md if you add new features
- Add examples in the examples directory
- Update the CHANGELOG.md with your changes
- Add inline documentation for complex logic

## Pull Request Process

1. Update the CHANGELOG.md with your changes
2. Ensure all tests pass
3. Update documentation as needed
4. Submit a pull request

## Release Process

1. Update the version in the CHANGELOG.md
2. Create a new tag:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```
3. The GitHub Action will automatically create a release

## Questions?

If you have any questions, please open an issue in the repository. 