# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Information

This is an **unofficial** Terraform provider for Mailtrap created by Oleksii Ianchuk. The repository is located at https://github.com/yanchuk/mailtrap-terraform.

## Building and Development

This is a Terraform provider for Mailtrap written in Go. Build the provider using:

```bash
go mod tidy
go build -o terraform-provider-mailtrap
```

The provider requires Go 1.21+ and uses the Terraform Plugin Framework v1.5.0.
Make sure that we strictly follow data types.
Make sure that code is covered by tests.

## Testing

The project has comprehensive unit test coverage. Run tests using:

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...
```

Test coverage includes:
- HTTP client functionality (`internal/client/client_test.go`)
- Provider configuration (`internal/provider/provider_test.go`)
- Resource operations - CRUD (`internal/provider/resource_*_test.go`)
- Data source operations (`internal/provider/data_source_*_test.go`)
- Model validation

## Provider Architecture

The codebase follows standard Terraform provider patterns:

- **main.go**: Entry point that serves the provider using Terraform's plugin framework
- **internal/provider/**: Core provider implementation
  - `provider.go`: Main provider struct with configuration (API token, account ID)
  - `resource_*.go`: Resource implementations (projects, inboxes, sending domains)
  - `data_source_*.go`: Data source implementations for reading existing resources
- **internal/client/**: HTTP client for Mailtrap API
  - `client.go`: HTTP client with authentication and multi-endpoint support
  - `models.go`: API response/request models

The provider supports multiple Mailtrap API endpoints:
- Main API: `https://mailtrap.io` (default)
- Sending API: `https://send.api.mailtrap.io`
- Sandbox API: `https://sandbox.api.mailtrap.io`

## Resources and Data Sources

The provider manages three main resource types:
1. **Projects**: Top-level containers for inboxes
2. **Inboxes**: Email testing environments with SMTP credentials
3. **Sending Domains**: Production email domains with DNS verification

Each resource has a corresponding data source for reading existing resources.

## Authentication

Uses API token authentication via:
- `api_token` provider configuration
- `MAILTRAP_API_TOKEN` environment variable

Account ID can be set via:
- `account_id` provider configuration
- `MAILTRAP_ACCOUNT_ID` environment variable

## Examples

The `examples/` directory contains integration examples:
- `basic/`: Simple project and inbox creation
- `aws-integration/`: Storing SMTP credentials in AWS Parameter Store
- `cloudflare-dns/`: Configuring DNS records for sending domains

## Mailtrap OpenAPI specifications
/docs/mailtrap-open-api directory consists YML files of Mailtrap OpenAPI specifications