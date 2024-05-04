#!/bin/sh

# Build the Docker image
docker build -t sentiment-analysis -f ./docker/Dockerfile .

# Authenticate with ECR
aws ecr get-login-password --region ap-northeast-2 | docker login --username AWS --password-stdin 867976255728.dkr.ecr.ap-northeast-2.amazonaws.com

# Tag the Docker image with the ECR repository URL
docker tag sentiment-analysis:latest 867976255728.dkr.ecr.ap-northeast-2.amazonaws.com/sentiment-analysis-fastapi:latest

# Push the Docker image to ECR
docker push 867976255728.dkr.ecr.ap-northeast-2.amazonaws.com/sentiment-analysis-fastapi:latest