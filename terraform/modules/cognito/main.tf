# Cognito User Pool
resource "aws_cognito_user_pool" "main" {
  name = "${var.environment}-b2b-platform-users"

  username_attributes      = ["email"]
  auto_verified_attributes = ["email"]

  password_policy {
    minimum_length    = 8
    require_lowercase = true
    require_numbers   = true
    require_symbols   = true
    require_uppercase = true
  }

  account_recovery_setting {
    recovery_mechanism {
      name     = "verified_email"
      priority = 1
    }
  }

  schema {
    name                = "email"
    attribute_data_type = "String"
    required            = true
    mutable             = true
  }

  schema {
    name                = "tenant_id"
    attribute_data_type = "String"
    required            = false
    mutable             = true
  }

  tags = {
    Name = "${var.environment}-cognito-user-pool"
  }
}

# Cognito User Pool Client
resource "aws_cognito_user_pool_client" "main" {
  name         = "${var.environment}-b2b-platform-client"
  user_pool_id = aws_cognito_user_pool.main.id

  generate_secret                      = true
  allowed_oauth_flows_user_pool_client = true
  allowed_oauth_flows                  = ["code", "implicit"]
  allowed_oauth_scopes                 = ["email", "openid", "profile"]
  supported_identity_providers          = ["COGNITO"]

  callback_urls = [
    "https://${var.domain_name}/callback",
    "http://localhost:3000/callback"
  ]

  logout_urls = [
    "https://${var.domain_name}/logout",
    "http://localhost:3000/logout"
  ]

  explicit_auth_flows = [
    "ALLOW_USER_PASSWORD_AUTH",
    "ALLOW_REFRESH_TOKEN_AUTH",
    "ALLOW_USER_SRP_AUTH"
  ]
}

# Cognito User Pool Domain
resource "aws_cognito_user_pool_domain" "main" {
  domain       = "${var.environment}-b2b-platform"
  user_pool_id = aws_cognito_user_pool.main.id
}
