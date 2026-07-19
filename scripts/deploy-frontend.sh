#!/usr/bin/env bash
set -euo pipefail

: "${FRONTEND_BUCKET:?FRONTEND_BUCKET is required}"
: "${CLOUDFRONT_DISTRIBUTION_ID:?CLOUDFRONT_DISTRIBUTION_ID is required}"

pushd frontend >/dev/null
npm ci
npm audit --audit-level=high
npm run build
aws s3 sync dist/ "s3://${FRONTEND_BUCKET}" --delete
popd >/dev/null

aws cloudfront create-invalidation \
  --distribution-id "${CLOUDFRONT_DISTRIBUTION_ID}" \
  --paths "/*"
