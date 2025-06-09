variable "bucket_name" {
  description = "Name of the S3 bucket to be created for static website hosting"
  type        = string
} 

variable "backend_api_url" {
  description = "The base URL of the backend Go service that the frontend will consume"
  type        = string
} 