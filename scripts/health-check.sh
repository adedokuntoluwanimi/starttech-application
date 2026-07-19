#!/usr/bin/env bash
set -euo pipefail

: "${CLOUDFRONT_DOMAIN:?CLOUDFRONT_DOMAIN is required}"

base_url="https://${CLOUDFRONT_DOMAIN}"
curl --fail --silent --show-error --retry 10 --retry-delay 10 "${base_url}/" >/dev/null
curl --fail --silent --show-error --retry 10 --retry-delay 10 "${base_url}/api/v1/health"
