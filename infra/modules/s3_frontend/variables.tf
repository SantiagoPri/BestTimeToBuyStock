variable "bucket_name" {
  description = "Name of the S3 bucket to be created for static website hosting"
  type        = string
} 

variable "backend_api_url" {
  type        = string
  description = "The URL of the backend API that the frontend will communicate with"
} 