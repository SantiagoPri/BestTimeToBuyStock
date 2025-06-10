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
  description = "Optional override URL for the backend API. If not provided, will be constructed from the EC2 instance's public IP"
  type        = string
  default     = null
}

variable "frontend_public_url" {
  description = "URL of the deployed frontend, used for CORS"
  type        = string
} 