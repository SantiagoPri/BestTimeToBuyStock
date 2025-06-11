terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

# CloudFront distribution for backend
resource "aws_cloudfront_distribution" "backend" {
  enabled          = true
  is_ipv6_enabled  = true
  price_class      = "PriceClass_100"

  origin {
    domain_name = var.backend_domain_name
    origin_id   = "EC2Origin"

    custom_origin_config {
      http_port              = 8080
      https_port             = 443
      origin_protocol_policy = "http-only"
      origin_ssl_protocols   = ["TLSv1.2"]
    }

    custom_header {
      name  = "X-Custom-Header"
      value = "CloudFront"
    }
  }

  default_cache_behavior {
    allowed_methods  = ["DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT"]
    cached_methods   = ["HEAD", "GET"]
    target_origin_id = "EC2Origin"

    forwarded_values {
      query_string = true
      headers      = ["Authorization", "Content-Type", "Host", "Origin"]
      cookies {
        forward = "all"
      }
    }

    viewer_protocol_policy = "redirect-to-https"
    min_ttl               = 0
    default_ttl           = 0
    max_ttl               = 0
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    cloudfront_default_certificate = true
  }

  tags = {
    Name = "backend-distribution"
  }
}

# Create a CORS policy for CloudFront
resource "aws_cloudfront_response_headers_policy" "cors_policy" {
  name    = "cors-policy"
  comment = "Policy for handling CORS headers"

  cors_config {
    access_control_allow_credentials = true
    
    access_control_allow_headers {
      items = ["Authorization", "Content-Type", "Origin"]
    }
    
    access_control_allow_methods {
      items = ["POST", "OPTIONS"]
    }
    
    access_control_allow_origins {
      items = ["*"]
    }
    
    origin_override = true
  }
} 