# AWS Integration Example - Store Mailtrap credentials in AWS Parameter Store

terraform {
  required_providers {
    mailtrap = {
      source  = "mailtrap/mailtrap"
      version = "~> 1.0"
    }
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

# Configure providers
provider "mailtrap" {
  # Configured via environment variables
}

provider "aws" {
  region = var.aws_region
}

# Variables
variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "development"
}

# Create a Mailtrap project for the environment
resource "mailtrap_project" "env_project" {
  name = "${var.environment} Email Testing"
}

# Create separate inboxes for different services
resource "mailtrap_inbox" "app_inbox" {
  project_id     = mailtrap_project.env_project.id
  name           = "Application Emails"
  email_username = "${var.environment}-app"
}

resource "mailtrap_inbox" "notification_inbox" {
  project_id     = mailtrap_project.env_project.id
  name           = "Notification Emails"
  email_username = "${var.environment}-notifications"
}

# Store SMTP credentials in AWS Parameter Store
resource "aws_ssm_parameter" "app_smtp_host" {
  name  = "/${var.environment}/mailtrap/app/smtp_host"
  type  = "String"
  value = mailtrap_inbox.app_inbox.domain
  
  tags = {
    Environment = var.environment
    Service     = "mailtrap"
    Type        = "smtp-config"
  }
}

resource "aws_ssm_parameter" "app_smtp_port" {
  name  = "/${var.environment}/mailtrap/app/smtp_port"
  type  = "String"
  value = tostring(mailtrap_inbox.app_inbox.smtp_ports[0])
  
  tags = {
    Environment = var.environment
    Service     = "mailtrap"
    Type        = "smtp-config"
  }
}

resource "aws_ssm_parameter" "app_smtp_username" {
  name  = "/${var.environment}/mailtrap/app/smtp_username"
  type  = "String"
  value = mailtrap_inbox.app_inbox.username
  
  tags = {
    Environment = var.environment
    Service     = "mailtrap"
    Type        = "smtp-config"
  }
}

resource "aws_ssm_parameter" "app_smtp_password" {
  name  = "/${var.environment}/mailtrap/app/smtp_password"
  type  = "SecureString"
  value = mailtrap_inbox.app_inbox.password
  
  tags = {
    Environment = var.environment
    Service     = "mailtrap"
    Type        = "smtp-config"
  }
}

# Create an IAM policy for accessing the parameters
resource "aws_iam_policy" "mailtrap_params_read" {
  name        = "${var.environment}-mailtrap-params-read"
  description = "Allow reading Mailtrap SMTP parameters"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "ssm:GetParameter",
          "ssm:GetParameters",
          "ssm:GetParametersByPath"
        ]
        Resource = [
          "arn:aws:ssm:${var.aws_region}:*:parameter/${var.environment}/mailtrap/*"
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "kms:Decrypt"
        ]
        Resource = ["*"]
        Condition = {
          StringEquals = {
            "kms:ViaService" = "ssm.${var.aws_region}.amazonaws.com"
          }
        }
      }
    ]
  })
}

# Outputs
output "parameter_paths" {
  value = {
    app_inbox = {
      host     = aws_ssm_parameter.app_smtp_host.name
      port     = aws_ssm_parameter.app_smtp_port.name
      username = aws_ssm_parameter.app_smtp_username.name
      password = aws_ssm_parameter.app_smtp_password.name
    }
  }
  description = "SSM Parameter paths for SMTP configuration"
}

output "iam_policy_arn" {
  value       = aws_iam_policy.mailtrap_params_read.arn
  description = "IAM policy ARN for reading Mailtrap parameters"
}

output "inbox_emails" {
  value = {
    app           = "${mailtrap_inbox.app_inbox.email_username}@${mailtrap_inbox.app_inbox.email_domain}"
    notifications = "${mailtrap_inbox.notification_inbox.email_username}@${mailtrap_inbox.notification_inbox.email_domain}"
  }
  description = "Email addresses for the inboxes"
}
