# StartTech Application

React/Vite frontend and Go/Gin backend adapted from the supplied MuchToDo `feature/full-stack` branch for delivery through Amazon S3, CloudFront, ECR, and EKS.

## Runtime architecture

- CloudFront serves the static SPA from private S3 over HTTPS.
- Browser requests use the relative `/api/v1` base path, avoiding mixed content and cross-domain session issues.
- CloudFront forwards `/api/*` to the backend ALB with headers, cookies, and query strings intact and caching disabled.
- The ALB targets EKS worker nodes on fixed NodePort `30080`; Kubernetes forwards traffic to container port `8080`.
- The Go service reads `MONGO_URI`, `REDIS_HOST`, and other runtime values from Kubernetes Secret and ConfigMap resources.
- Application logs use structured JSON on stdout in production.

## Layout

```text
frontend/                 React 19 and Vite SPA
backend/                  Go 1.25 Gin API and Dockerfile
k8s/                      Deployment, Service, and ALB Ingress manifests
scripts/                  Deploy, health-check, and rollback helpers
.github/workflows/        Frontend and backend CI/CD
```

## Local development

Frontend:

```bash
cd frontend
cp .env.example .env
npm ci
npm run dev
```

Backend:

```bash
cd backend
cp .env.example .env
go test ./...
go run ./cmd/api
```

The local frontend environment points to `http://localhost:8080/api/v1`; production leaves `VITE_API_BASE_URL` unset and uses `/api/v1` automatically.

## GitHub Actions configuration

Repository variables:

- `AWS_REGION`
- `AWS_APPLICATION_ROLE_ARN`
- `FRONTEND_BUCKET`
- `CLOUDFRONT_DISTRIBUTION_ID`
- `CLOUDFRONT_DOMAIN`

Repository secrets:

- `MONGO_URI`: MongoDB Atlas connection string
- `JWT_SECRET_KEY`: long, randomly generated signing secret

Never commit either secret. The backend workflow creates or updates Kubernetes runtime secrets without printing their values.

## Delivery behavior

- Frontend changes run `npm ci`, `npm audit`, and `npm run build`, synchronize `dist/` to S3, and invalidate CloudFront.
- Backend or manifest changes run `go test ./...`, build a SHA-tagged image, scan it with Trivy, push it to `starttech-backend-api`, apply Kubernetes manifests, and wait for a successful rollout.
- `scripts/health-check.sh` verifies both the SPA and `/api/v1/health` through the CloudFront domain.
- `scripts/rollback.sh` rolls the backend deployment back and waits for the rollback to finish.

The `ingress.yaml` manifest documents the AWS Load Balancer Controller path required by the brief. The deployed infrastructure uses the explicit NodePort target path supported by the same brief, which avoids creating a second ALB and allows CloudFront to receive the ALB DNS name during the first Terraform apply.
