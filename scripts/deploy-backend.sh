#!/usr/bin/env bash
set -euo pipefail

: "${AWS_ACCOUNT_ID:?AWS_ACCOUNT_ID is required}"
: "${AWS_REGION:?AWS_REGION is required}"
: "${EKS_CLUSTER_NAME:=starttech-cluster}"
: "${ECR_REPOSITORY:=starttech-backend-api}"
: "${MONGO_URI:?MONGO_URI is required}"
: "${JWT_SECRET_KEY:?JWT_SECRET_KEY is required}"
: "${CLOUDFRONT_DOMAIN:?CLOUDFRONT_DOMAIN is required}"

image_tag="${IMAGE_TAG:-$(git rev-parse HEAD)}"
registry="${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com"
image="${registry}/${ECR_REPOSITORY}:${image_tag}"

aws ecr get-login-password --region "${AWS_REGION}" |
  docker login --username AWS --password-stdin "${registry}"
docker build --tag "${image}" backend
docker push "${image}"

aws eks update-kubeconfig --region "${AWS_REGION}" --name "${EKS_CLUSTER_NAME}"
redis_host="$(aws elasticache describe-cache-clusters \
  --cache-cluster-id starttech-redis \
  --show-cache-node-info \
  --query 'CacheClusters[0].CacheNodes[0].Endpoint.Address' \
  --output text)"

kubectl create configmap backend-config \
  --from-literal=PORT=8080 \
  --from-literal=DB_NAME=much_todo_db \
  --from-literal=ENABLE_CACHE=true \
  --from-literal=REDIS_HOST="${redis_host}" \
  --from-literal=REDIS_PORT=6379 \
  --from-literal=LOG_LEVEL=INFO \
  --from-literal=LOG_FORMAT=json \
  --from-literal=SECURE_COOKIE=true \
  --from-literal=ALLOWED_ORIGINS="https://${CLOUDFRONT_DOMAIN}" \
  --dry-run=client -o yaml | kubectl apply -f -

kubectl create secret generic backend-secrets \
  --from-literal=MONGO_URI="${MONGO_URI}" \
  --from-literal=JWT_SECRET_KEY="${JWT_SECRET_KEY}" \
  --dry-run=client -o yaml | kubectl apply -f -

sed "s|starttech-backend-api:latest|${image}|" k8s/deployment.yaml | kubectl apply -f -
kubectl apply -f k8s/service.yaml
kubectl rollout status deployment/backend-api --timeout=5m
