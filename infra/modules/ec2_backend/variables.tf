variable "instance_name" {
  description = "Name of the EC2 instance"
  type        = string
}

variable "docker_image" {
  description = "Docker image for the backend application"
  type        = string
}

variable "database_url" {
  description = "URL for the database connection"
  type        = string
}

variable "redis_password" {
  description = "Password for Redis"
  type        = string
  sensitive   = true
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

variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "frontend_public_url" {
  description = "URL of the deployed frontend, used for CORS"
  type        = string
} 