output "bucket_name" {
  description = "Name of the S3 bucket created for static website hosting"
  value       = module.s3_frontend.bucket_name
}

output "website_endpoint" {
  description = "The website endpoint URL for the static website"
  value       = module.s3_frontend.website_endpoint
} 