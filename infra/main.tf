terraform {
  required_version = ">= 1.3.0"
  
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  # Local backend configuration
  backend "local" {}
}

provider "aws" {
  region = "us-east-1"
}

locals {
  backend_api_resolved = var.backend_api_url != null ? var.backend_api_url : "http://${module.ec2_backend.public_ip}:8080"
}

module "ec2_backend" {
  source = "./modules/ec2_backend"

  instance_name         = "besttime-backend"
  docker_image         = "your-dockerhub-user/besttime-backend:latest"
  database_url         = var.database_url
  redis_password       = var.redis_password
  openrouter_api_key   = var.openrouter_api_key
  openrouter_model_name = var.openrouter_model_name
}

module "s3_frontend" {
  source = "./modules/s3_frontend"

  bucket_name     = var.bucket_name
  backend_api_url = local.backend_api_resolved
}

`