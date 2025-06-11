data "aws_ami" "amazon_linux_2" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn2-ami-hvm-*-x86_64-gp2"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

resource "aws_security_group" "backend" {
  name_prefix = "backend-"
  description = "Security group for backend EC2 instance"

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "SSH access"
  }

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "HTTP access"
  }

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "HTTPS access"
  }

  ingress {
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "Backend application access"
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
    description = "Allow all outbound traffic"
  }

  tags = {
    Name = "${var.instance_name}-sg"
  }
}

resource "aws_key_pair" "this" {
  key_name   = "besttime-key"
  public_key = file("~/.ssh/besttime_key.pub")
}

resource "aws_instance" "backend" {
  ami           = data.aws_ami.amazon_linux_2.id
  instance_type = "t2.micro"
  key_name      = aws_key_pair.this.key_name

  vpc_security_group_ids = [aws_security_group.backend.id]

  user_data = templatefile("${path.module}/user_data.sh.tpl", {
    docker_image          = var.docker_image
    database_url          = var.database_url
    redis_password        = var.redis_password
    openrouter_api_key    = var.openrouter_api_key
    openrouter_model_name = var.openrouter_model_name
    frontend_public_url   = var.frontend_public_url
  })

  tags = {
    Name = var.instance_name
  }

  root_block_device {
    volume_size = 8
    volume_type = "gp2"
    encrypted   = true
  }
} 