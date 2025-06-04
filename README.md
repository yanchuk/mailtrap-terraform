# Terraform Provider for Mailtrap

The Terraform Mailtrap provider allows you to manage Mailtrap resources using Terraform.

## Features

- Manage Projects
- Manage Inboxes (with SMTP credentials)
- Manage Sending Domains (with DNS records for configuration)
- Data sources for reading existing resources

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21 (to build the provider plugin)

## Building The Provider

1. Clone the repository
```bash
git clone https://github.com/mailtrap/terraform-provider-mailtrap.git
cd terraform-provider-mailtrap
```

2. Build the provider
```bash
go build -o terraform-provider-mailtrap
```

3. Install the provider locally
```bash
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/mailtrap/mailtrap/1.0.0/darwin_arm64
mv terraform-provider-mailtrap ~/.terraform.d/plugins/registry.terraform.io/mailtrap/mailtrap/1.0.0/darwin_arm64/
```

## Using the Provider

### Provider Configuration

```hcl
terraform {
  required_providers {
    mailtrap = {
      source  = "mailtrap/mailtrap"
      version = "~> 1.0"
    }
  }
}

# Configure the Mailtrap Provider
provider "mailtrap" {
  api_token  = var.mailtrap_api_token
  account_id = var.mailtrap_account_id
}
```

You can also use environment variables:
- `MAILTRAP_API_TOKEN`
- `MAILTRAP_ACCOUNT_ID`

### Example Usage

#### Create a Project with Inbox

```hcl
# Create a project
resource "mailtrap_project" "example" {
  name = "My Test Project"
}

# Create an inbox in the project
resource "mailtrap_inbox" "example" {
  project_id     = mailtrap_project.example.id
  name           = "My Test Inbox"
  email_username = "test-inbox"
}

# Output SMTP credentials
output "smtp_credentials" {
  value = {
    host     = mailtrap_inbox.example.domain
    port     = mailtrap_inbox.example.smtp_ports[0]
    username = mailtrap_inbox.example.username
    password = mailtrap_inbox.example.password
  }
  sensitive = true
}
```

#### Create a Sending Domain

```hcl
# Create a sending domain
resource "mailtrap_sending_domain" "example" {
  name = "example.com"
}

# Output DNS records for configuration
output "dns_records" {
  value = {
    cname_records = [
      for record in mailtrap_sending_domain.example.dns_records.cname : {
        type     = record.record_type
        hostname = record.hostname
        value    = record.value
      }
    ]
    mx_records = [
      for record in mailtrap_sending_domain.example.dns_records.mx : {
        type     = record.record_type
        hostname = record.hostname
        value    = record.value
        priority = record.priority
      }
    ]
    txt_records = [
      for record in mailtrap_sending_domain.example.dns_records.txt : {
        type     = record.record_type
        hostname = record.hostname
        value    = record.value
      }
    ]
  }
}
```

#### Store Credentials in AWS Parameter Store

```hcl
# Store SMTP credentials in AWS Parameter Store
resource "aws_ssm_parameter" "smtp_host" {
  name  = "/mailtrap/smtp/host"
  type  = "String"
  value = mailtrap_inbox.example.domain
}

resource "aws_ssm_parameter" "smtp_username" {
  name  = "/mailtrap/smtp/username"
  type  = "String"
  value = mailtrap_inbox.example.username
}

resource "aws_ssm_parameter" "smtp_password" {
  name  = "/mailtrap/smtp/password"
  type  = "SecureString"
  value = mailtrap_inbox.example.password
}
```

#### Configure Cloudflare DNS with Mailtrap

```hcl
# Configure Cloudflare DNS records for Mailtrap sending domain
resource "cloudflare_record" "mailtrap_cname" {
  for_each = {
    for idx, record in mailtrap_sending_domain.example.dns_records.cname : 
    idx => record
  }

  zone_id = var.cloudflare_zone_id
  name    = each.value.hostname
  value   = each.value.value
  type    = "CNAME"
  ttl     = 3600
}

resource "cloudflare_record" "mailtrap_mx" {
  for_each = {
    for idx, record in mailtrap_sending_domain.example.dns_records.mx : 
    idx => record
  }

  zone_id  = var.cloudflare_zone_id
  name     = each.value.hostname
  value    = each.value.value
  type     = "MX"
  priority = each.value.priority
  ttl      = 3600
}

resource "cloudflare_record" "mailtrap_txt" {
  for_each = {
    for idx, record in mailtrap_sending_domain.example.dns_records.txt : 
    idx => record
  }

  zone_id = var.cloudflare_zone_id
  name    = each.value.hostname
  value   = each.value.value
  type    = "TXT"
  ttl     = 3600
}
```

## Resources

### mailtrap_project

Creates and manages a Mailtrap project.

#### Arguments

- `name` - (Required) The name of the project (min 2 characters, max 100 characters).
- `account_id` - (Optional) The account ID. If not specified, uses the provider's account_id.

#### Attributes

- `id` - The project ID.
- `share_links` - Share links for the project.
  - `admin` - Admin share link.
  - `viewer` - Viewer share link.

### mailtrap_inbox

Creates and manages a Mailtrap inbox.

#### Arguments

- `project_id` - (Required) The project ID where the inbox will be created.
- `name` - (Required) The name of the inbox.
- `email_username` - (Optional) The email username part (before @) for the inbox email address.
- `account_id` - (Optional) The account ID. If not specified, uses the provider's account_id.

#### Attributes

- `id` - The inbox ID.
- `username` - SMTP username.
- `password` - SMTP password (sensitive).
- `email_username_enabled` - Whether email username is enabled.
- `domain` - Domain for SMTP.
- `email_domain` - Email domain.
- `pop3_domain` - POP3 domain.
- `smtp_ports` - List of available SMTP ports.
- `pop3_ports` - List of available POP3 ports.
- `status` - Inbox status.
- `max_size` - Maximum inbox size.
- `sent_messages_count` - Number of sent messages.
- `forwarded_messages_count` - Number of forwarded messages.
- `forward_from_email_address` - Email address used for forwarding.

### mailtrap_sending_domain

Creates and manages a Mailtrap sending domain.

#### Arguments

- `name` - (Required) The domain name.
- `account_id` - (Optional) The account ID. If not specified, uses the provider's account_id.

#### Attributes

- `id` - The sending domain ID.
- `cname` - CNAME value for domain verification.
- `status` - Domain status.
- `compliance_status` - Compliance status.
- `dns_records` - DNS records for domain verification.
  - `cname` - List of CNAME records.
  - `mx` - List of MX records.
  - `txt` - List of TXT records.
- `dns_status` - DNS verification status.
  - `cname` - CNAME verification status.
  - `mx` - MX verification status.
  - `txt` - TXT verification status.

## Data Sources

### mailtrap_account

Reads account information.

#### Arguments

- `id` - (Required) The account ID.

#### Attributes

- `name` - The account name.

### mailtrap_project

Reads project information.

#### Arguments

- `id` - (Required) The project ID.
- `account_id` - (Optional) The account ID. If not specified, uses the provider's account_id.

#### Attributes

All attributes from the `mailtrap_project` resource.

### mailtrap_inbox

Reads inbox information.

#### Arguments

- `id` - (Required) The inbox ID.
- `account_id` - (Optional) The account ID. If not specified, uses the provider's account_id.

#### Attributes

All attributes from the `mailtrap_inbox` resource.

### mailtrap_sending_domain

Reads sending domain information.

#### Arguments

- `id` - (Required) The sending domain ID.
- `account_id` - (Optional) The account ID. If not specified, uses the provider's account_id.

#### Attributes

All attributes from the `mailtrap_sending_domain` resource.

## Importing Resources

Resources can be imported using the format `account_id/resource_id`.

```bash
# Import a project
terraform import mailtrap_project.example 12345/67890

# Import an inbox
terraform import mailtrap_inbox.example 12345/67890

# Import a sending domain
terraform import mailtrap_sending_domain.example 12345/67890
```

## Future Enhancements

The following features are planned for future releases:

- API Token resource (when API support is available)
- Webhook resource (when API support is available)
- User and permission management
- Contact list management

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This provider is distributed under the [Mozilla Public License 2.0](LICENSE).
