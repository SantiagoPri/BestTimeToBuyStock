resource "aws_s3_bucket" "website" {
  bucket = var.bucket_name
}

resource "aws_s3_bucket_website_configuration" "website" {
  bucket = aws_s3_bucket.website.id

  index_document {
    suffix = "index.html"
  }

  error_document {
    key = "404.html"
  }
}

resource "aws_s3_bucket_public_access_block" "website" {
  bucket = aws_s3_bucket.website.id

  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}

resource "aws_s3_bucket_policy" "website" {
  bucket = aws_s3_bucket.website.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid       = "PublicReadGetObject"
        Effect    = "Allow"
        Principal = "*"
        Action    = "s3:GetObject"
        Resource  = "${aws_s3_bucket.website.arn}/*"
      }
    ]
  })

  depends_on = [aws_s3_bucket_public_access_block.website]
}

resource "aws_s3_bucket_cors_configuration" "website" {
  bucket = aws_s3_bucket.website.id

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["GET", "HEAD"]
    allowed_origins = ["*"]
    expose_headers  = ["ETag"]
    max_age_seconds = 3000
  }
}

resource "local_file" "env_file" {
  filename = "${path.module}/../frontend/.env"
  content  = "API_URL=${var.backend_api_url}"
  file_permission = "0644"
}

resource "null_resource" "build_and_deploy_frontend" {
  triggers = {
    env_file_hash = filesha256("../frontend/.env")
  }

  provisioner "local-exec" {
    command = <<EOT
      if ! which npm > /dev/null; then
        echo "npm not found"
        exit 1
      fi

      cd ../frontend
      npm install
      npm run build
      aws s3 sync ./dist s3://${aws_s3_bucket.website.id} --delete
EOT
  }

  depends_on = [local_file.env_file]
} 