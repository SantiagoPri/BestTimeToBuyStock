variable "bucket_name" {
  description = "Name of the S3 bucket for static website hosting. If not provided, defaults to the hardcoded value in main.tf"
  type        = string
  default     = null
}

variable "database_url" {
  description = "URL for the database connection"
  type        = string
}

variable "redis_password" {
  description = "Password for Redis"
  type        = string
}

variable "openrouter_api_key" {
  description = "API key for OpenRouter"
  type        = string
  sensitive   = true
}

variable "openrouter_model_name" {
  description = "Model name for OpenRouter"
  type        = string
}

variable "backend_api_url" {
  description = "URL of the backend, used to configure the frontend"
  type        = string
  validation {
    condition     = length(var.backend_api_url) > 0
    error_message = "The backend_api_url variable must be set in terraform.tfvars"
  }
}

variable "frontend_public_url" {
  description = "URL of the deployed frontend, used for CORS"
  type        = string
  validation {
    condition     = length(var.frontend_public_url) > 0
    error_message = "The frontend_public_url variable must be set in terraform.tfvars"
  }
}

variable "duck_dns_token" {
  description = "Token for DuckDNS service"
  type        = string
  sensitive   = true
}

variable "domain_name" {
  description = "Domain name for the application"
  type        = string
} 