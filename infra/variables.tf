variable "bucket_name" {
  description = "Name of the S3 bucket for static website hosting. If not provided, defaults to the hardcoded value in main.tf"
  type        = string
  default     = null
} 