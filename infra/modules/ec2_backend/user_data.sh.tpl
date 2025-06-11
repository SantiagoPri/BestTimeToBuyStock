#!/bin/bash
set -e

# Redirect all output to log file
exec > >(tee /var/log/user_data.log) 2>&1

echo "[$(date)] Starting user data script execution"

# Update system packages
yum update -y

# Install Docker
yum install -y docker
systemctl start docker
systemctl enable docker

# Install Docker Compose
curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

# Start Redis container
docker run -d \
  --name redis \
  -p 6379:6379 \
  redis:latest \
  redis-server --requirepass "${redis_password}"

# Start backend container
docker run -d \
  --name backend \
  -p 8080:8080 \
  -e DATABASE_URL="${database_url}" \
  -e REDIS_HOST="localhost" \
  -e REDIS_PORT="6379" \
  -e REDIS_PASSWORD="${redis_password}" \
  -e OPENROUTER_API_KEY="${openrouter_api_key}" \
  -e OPENROUTER_MODEL_NAME="${openrouter_model_name}" \
  -e FRONTEND_PUBLIC_URL="${frontend_public_url}" \
  --network="host" \
  ${docker_image}

// echo "[$(date)] User data script execution completed"
