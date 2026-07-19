#!/usr/bin/env bash
set -euo pipefail

: "${AWS_REGION:?AWS_REGION is required}"
: "${EKS_CLUSTER_NAME:=starttech-cluster}"

aws eks update-kubeconfig --region "${AWS_REGION}" --name "${EKS_CLUSTER_NAME}"
kubectl rollout undo deployment/backend-api
kubectl rollout status deployment/backend-api --timeout=5m
