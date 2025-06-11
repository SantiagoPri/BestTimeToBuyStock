output "bucket_name" {
  description = "Name of the S3 bucket created for static website hosting"
  value       = module.s3_frontend.bucket_name
}

output "website_endpoint" {
  description = "The website endpoint URL for the static website"
  value       = module.s3_frontend.website_endpoint
}

output "ec2_backend_public_ip" {
  description = "Public IP address of the backend EC2 instance"
  value       = module.ec2_backend.public_ip
}

output "cloudfront_frontend_domain" {
  description = "The domain name of the frontend CloudFront distribution"
  value       = module.cloudfront_frontend.cloudfront_domain_name
}

output "cloudfront_backend_domain" {
  description = "The domain name of the backend CloudFront distribution"
  value       = module.cloudfront_backend.cloudfront_domain_name
} 