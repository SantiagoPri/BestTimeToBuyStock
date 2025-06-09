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

`