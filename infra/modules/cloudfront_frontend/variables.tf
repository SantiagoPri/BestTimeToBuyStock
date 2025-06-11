variable "bucket_domain_name" {
  description = "The domain name of the S3 website endpoint (e.g. my-bucket.s3-website-us-east-1.amazonaws.com)"
  type        = string
}

variable "region" {
  description = "The AWS region where the S3 bucket is located"
  type        = string
}