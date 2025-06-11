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
  backend_api_resolved = var.backend_api_url
  frontend_url_resolved = var.frontend_public_url
}

module "ec2_backend" {
  source = "./modules/ec2_backend"

  instance_name         = "besttime-backend"
  docker_image         = "elardeluz/besttime-backend:latest"
  database_url         = var.database_url
  redis_password       = var.redis_password
  openrouter_api_key   = var.openrouter_api_key
  openrouter_model_name = var.openrouter_model_name
  frontend_public_url   = local.frontend_url_resolved
}

module "s3_frontend" {
  source = "./modules/s3_frontend"

  bucket_name     = var.bucket_name
  backend_api_url = local.backend_api_resolved
}

module "cloudfront_frontend" {
  source = "./modules/cloudfront_frontend"

  bucket_domain_name = module.s3_frontend.website_endpoint
  region            = "us-east-1"
}

module "cloudfront_backend" {
  source = "./modules/cloudfront_backend"

  backend_domain_name = module.ec2_backend.public_dns
}