# Contributing to Terraform Provider for Mailtrap

Thank you for your interest in contributing! This document provides guidelines and information for contributors.

## Development Environment Setup

### Prerequisites

- [Go](https://golang.org/doc/install) 1.21 or later
- [Terraform](https://www.terraform.io/downloads.html) 1.0 or later
- Git

### Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/mailtrap-terraform.git
   cd mailtrap-terraform
   ```

3. Add the upstream repository:
   ```bash
   git remote add upstream https://github.com/yanchuk/mailtrap-terraform.git
   ```

4. Install dependencies:
   ```bash
   go mod download
   ```

5. Build the provider:
   ```bash
   go build -o terraform-provider-mailtrap
   ```

## Development Workflow

### Making Changes

1. Create a new branch for your feature or bugfix:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes and ensure they follow the existing code style

3. Write or update tests for your changes

4. Run tests to ensure everything works:
   ```bash
   go test ./...
   ```

5. Build the provider to ensure it compiles:
   ```bash
   go build -o terraform-provider-mailtrap
   ```

### Testing

We have comprehensive unit tests covering:
- HTTP client functionality
- Provider configuration
- Resource operations (CRUD)
- Data source operations
- Model validation

Run tests with:
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -run TestClientNewClient ./internal/client
```

### Code Standards

- Follow standard Go conventions and formatting
- Use `gofmt` to format your code
- Write meaningful commit messages
- Add tests for new functionality
- Update documentation when adding new features

### Commit Messages

Use clear and descriptive commit messages:
- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters or less
- Reference issues and pull requests liberally after the first line

Example:
```
Add support for inbox email forwarding

- Add email_forwarding field to inbox resource
- Update tests for new functionality
- Update documentation

Fixes #123
```

## Testing with Real Mailtrap API

For integration testing, you'll need:
1. A Mailtrap account
2. API token from your Mailtrap account
3. Account ID

Set environment variables:
```bash
export MAILTRAP_API_TOKEN="your-api-token"
export MAILTRAP_ACCOUNT_ID="your-account-id"
```

**Note**: Integration tests will create real resources in your Mailtrap account.

## Submitting Changes

1. Push your changes to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

2. Create a pull request on GitHub
3. Ensure the PR description clearly describes the problem and solution
4. Include any relevant issue numbers

### Pull Request Guidelines

- Keep pull requests focused on a single feature or bugfix
- Include tests for new functionality
- Update documentation as needed
- Ensure all tests pass
- Make sure the code builds successfully

## Code Review Process

1. Pull requests require review from maintainers
2. Address any feedback from reviewers
3. Once approved, a maintainer will merge the PR

## Reporting Issues

When reporting issues, please include:
- Go version
- Terraform version
- Provider version
- Minimal reproduction case
- Expected vs actual behavior
- Relevant logs or error messages

## Architecture Overview

The provider follows standard Terraform provider patterns:

```
├── internal/
│   ├── client/          # HTTP client for Mailtrap API
│   │   ├── client.go    # Main client implementation
│   │   └── models.go    # API models
│   └── provider/        # Terraform provider implementation
│       ├── provider.go  # Main provider
│       ├── resource_*.go      # Resource implementations
│       └── data_source_*.go   # Data source implementations
├── examples/            # Example Terraform configurations
└── docs/               # Documentation
```

## Resources and Data Sources

### Adding New Resources

1. Create `resource_newresource.go` in `internal/provider/`
2. Implement the Resource interface methods
3. Add models to `internal/client/models.go` if needed
4. Add tests in `resource_newresource_test.go`
5. Update documentation

### Adding New Data Sources

1. Create `data_source_newresource.go` in `internal/provider/`
2. Implement the DataSource interface methods
3. Add tests in `data_source_newresource_test.go`
4. Update documentation

## Release Process

Releases are handled by maintainers:
1. Update version in relevant files
2. Update CHANGELOG.md
3. Create and push a new tag
4. GitHub Actions will build and publish the release

## Getting Help

- Check existing issues and discussions
- Join our community discussions
- Reach out to maintainers

Thank you for contributing!