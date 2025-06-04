# Cloudflare DNS Integration Example - Configure DNS for Mailtrap Sending Domain

terraform {
  required_providers {
    mailtrap = {
      source  = "mailtrap/mailtrap"
      version = "~> 1.0"
    }
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 4.0"
    }
  }
}

# Configure providers
provider "mailtrap" {
  # Configured via environment variables
}

provider "cloudflare" {
  # Configured via environment variables
}

# Variables
variable "domain_name" {
  description = "Domain name to configure for Mailtrap sending"
  type        = string
  example     = "example.com"
}

variable "cloudflare_zone_id" {
  description = "Cloudflare zone ID for the domain"
  type        = string
}

# Create a Mailtrap sending domain
resource "mailtrap_sending_domain" "main" {
  name = var.domain_name
}

# Configure Cloudflare DNS records for Mailtrap
# CNAME Records
resource "cloudflare_record" "mailtrap_cname" {
  for_each = {
    for idx, record in mailtrap_sending_domain.main.dns_records.cname : 
    "${record.hostname}-${idx}" => record
  }

  zone_id = var.cloudflare_zone_id
  name    = each.value.hostname
  value   = each.value.value
  type    = "CNAME"
  ttl     = 3600
  proxied = false

  lifecycle {
    create_before_destroy = true
  }
}

# MX Records
resource "cloudflare_record" "mailtrap_mx" {
  for_each = {
    for idx, record in mailtrap_sending_domain.main.dns_records.mx : 
    "${record.hostname}-${idx}" => record
  }

  zone_id  = var.cloudflare_zone_id
  name     = each.value.hostname
  value    = each.value.value
  type     = "MX"
  priority = each.value.priority
  ttl      = 3600

  lifecycle {
    create_before_destroy = true
  }
}

# TXT Records
resource "cloudflare_record" "mailtrap_txt" {
  for_each = {
    for idx, record in mailtrap_sending_domain.main.dns_records.txt : 
    "${record.hostname}-${idx}" => record
  }

  zone_id = var.cloudflare_zone_id
  name    = each.value.hostname
  value   = each.value.value
  type    = "TXT"
  ttl     = 3600

  lifecycle {
    create_before_destroy = true
  }
}

# Outputs
output "sending_domain_status" {
  value = {
    domain            = mailtrap_sending_domain.main.name
    status           = mailtrap_sending_domain.main.status
    compliance_status = mailtrap_sending_domain.main.compliance_status
    dns_verified = {
      cname = mailtrap_sending_domain.main.dns_status.cname
      mx    = mailtrap_sending_domain.main.dns_status.mx
      txt   = mailtrap_sending_domain.main.dns_status.txt
    }
  }
  description = "Mailtrap sending domain status"
}

output "dns_records_created" {
  value = {
    cname_count = length(cloudflare_record.mailtrap_cname)
    mx_count    = length(cloudflare_record.mailtrap_mx)
    txt_count   = length(cloudflare_record.mailtrap_txt)
  }
  description = "Number of DNS records created in Cloudflare"
}

output "verification_instructions" {
  value = <<-EOT
    DNS records have been created in Cloudflare for ${var.domain_name}.
    
    The domain verification status is:
    - CNAME: ${mailtrap_sending_domain.main.dns_status.cname ? "Verified ✓" : "Pending ⏳"}
    - MX: ${mailtrap_sending_domain.main.dns_status.mx ? "Verified ✓" : "Pending ⏳"}
    - TXT: ${mailtrap_sending_domain.main.dns_status.txt ? "Verified ✓" : "Pending ⏳"}
    
    DNS propagation may take up to 48 hours. Run 'terraform refresh' to check the latest status.
  EOT
  description = "Instructions for domain verification"
}
