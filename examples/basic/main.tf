# Basic Example - Create a Project and Inbox

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
  # API token can be set via MAILTRAP_API_TOKEN environment variable
  # api_token = "your-api-token"
  
  # Account ID can be set via MAILTRAP_ACCOUNT_ID environment variable
  # account_id = 12345
}

# Create a project
resource "mailtrap_project" "example" {
  name = "My Development Project"
}

# Create an inbox in the project
resource "mailtrap_inbox" "dev" {
  project_id     = mailtrap_project.example.id
  name           = "Development Inbox"
  email_username = "dev-emails"
}

# Output the project share links
output "project_share_links" {
  value = {
    admin  = mailtrap_project.example.share_links.admin
    viewer = mailtrap_project.example.share_links.viewer
  }
}

# Output SMTP configuration
output "smtp_config" {
  value = {
    host     = mailtrap_inbox.dev.domain
    ports    = mailtrap_inbox.dev.smtp_ports
    username = mailtrap_inbox.dev.username
    email    = "${mailtrap_inbox.dev.email_username}@${mailtrap_inbox.dev.email_domain}"
  }
}

# Output sensitive SMTP password separately
output "smtp_password" {
  value     = mailtrap_inbox.dev.password
  sensitive = true
}
