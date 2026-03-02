# Open-PM

Open-source project management platform. Manage tasks, sprints, epics, docs, and triage — built with Go, Vue 3, and Tailwind CSS.

## Table of Contents

- [Stack](#stack)
- [Local Development](#local-development)
- [Docker Compose](#docker-compose)
- [Production Build](#production-build)
- [Kubernetes](#kubernetes)
- [Environment Variables](#environment-variables)

---

## Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go 1.24, Chi router, pgx/v5 |
| Frontend | Vue 3, TypeScript, Vite, Tailwind CSS |
| Database | PostgreSQL 16 |
| Cache | Redis 7 |
| Object Storage | MinIO (S3-compatible) |
| Email (dev) | MailHog |

---

## Local Development

### Prerequisites

- [Go 1.24+](https://go.dev/dl/)
- [Node.js 22+](https://nodejs.org/)
- [Docker + Docker Compose](https://docs.docker.com/get-docker/)

### First-time setup

```bash
# 1. Clone the repo
git clone https://github.com/open-pm/open-pm.git
cd open-pm

# 2. Start infrastructure (PostgreSQL, Redis, MinIO, MailHog)
make dev

# 3. Run database migrations
make db-migrate

# 4. Install frontend dependencies
cd web && npm install && cd ..
```

### Copy environment file

```bash
cp server/.env.example server/.env   # edit values as needed
```

> If no `.env.example` exists yet, the server falls back to its built-in defaults (see [Environment Variables](#environment-variables)).

### Run

Open **two terminals**:

```bash
# Terminal 1 — backend (http://localhost:8080)
make dev-server

# Terminal 2 — frontend (http://localhost:3000)
make dev-web
```

### Service URLs

| Service | URL |
|---------|-----|
| Frontend | <http://localhost:3000> |
| Backend API | <http://localhost:8080> |
| MailHog UI | <http://localhost:8025> |
| MinIO Console | <http://localhost:9001> |
| PostgreSQL | `localhost:5432` |
| Redis | `localhost:6379` |

### Default seed credentials (development only)

| Field | Value |
|-------|-------|
| Email | `admin@open-pm.dev` |
| Password | `password123` |

### Useful make targets

```bash
make db-migrate          # Apply all pending migrations
make db-migrate-down     # Roll back one migration
make db-migrate-create   # Scaffold a new migration file
make sqlc-generate       # Re-generate sqlc query code
make lint                # Lint Go + TypeScript
make test                # Run all tests
make build               # Build server binary + frontend bundle
```

---

## Docker Compose

The `docker-compose.yml` at the root starts only **infrastructure** services (database, cache, storage, mail). This is intentional for development so you run the server and frontend natively with hot-reload.

### Infrastructure only (dev mode)

```bash
docker compose up -d
```

### Full stack (app + infrastructure)

Create a `docker-compose.prod.yml` (or copy the snippet below) to add the application containers:

```yaml
# docker-compose.prod.yml
version: '3.9'

services:
  server:
    build:
      context: ./server
      dockerfile: Dockerfile
    container_name: openpm-server
    environment:
      PORT: "8080"
      ENVIRONMENT: production
      DATABASE_URL: postgres://openpm:openpm_secret@postgres:5432/openpm?sslmode=disable
      JWT_SECRET: "change-me-in-production"
      REDIS_URL: redis://redis:6379/0
      STORAGE_ENDPOINT: minio:9000
      STORAGE_ACCESS_KEY: minioadmin
      STORAGE_SECRET_KEY: minioadmin
      STORAGE_BUCKET: open-pm-assets
      SMTP_HOST: mailhog
      SMTP_PORT: "1025"
    ports:
      - "8080:8080"
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    restart: unless-stopped

  web:
    build:
      context: ./web
      dockerfile: Dockerfile
    container_name: openpm-web
    ports:
      - "80:80"
    depends_on:
      - server
    restart: unless-stopped
```

Run the full stack:

```bash
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# Apply migrations (first run or after updates)
docker compose exec server ./server migrate up
```

---

## Production Build

### Build Docker images

```bash
# Backend
docker build -t open-pm/server:latest ./server

# Frontend
docker build -t open-pm/web:latest ./web
```

### Build with a specific version tag

```bash
VERSION=$(git rev-parse --short HEAD)

docker build -t open-pm/server:${VERSION} ./server
docker build -t open-pm/web:${VERSION} ./web
```

### Run standalone containers

```bash
# Backend — requires a running PostgreSQL and the env vars below
docker run -d \
  --name openpm-server \
  -p 8080:8080 \
  -e DATABASE_URL="postgres://user:pass@host:5432/openpm?sslmode=require" \
  -e JWT_SECRET="your-secret" \
  -e REDIS_URL="redis://host:6379/0" \
  -e STORAGE_ENDPOINT="s3.amazonaws.com" \
  -e STORAGE_ACCESS_KEY="AKID..." \
  -e STORAGE_SECRET_KEY="secret..." \
  -e STORAGE_BUCKET="open-pm-assets" \
  open-pm/server:latest

# Frontend (proxies /api and /auth to the server container)
docker run -d \
  --name openpm-web \
  -p 80:80 \
  open-pm/web:latest
```

> The frontend nginx config proxies `/api/` and `/auth/` to `http://server:8080`. When running standalone (not on the same Docker network), update `nginx.conf` to point to the correct server address.

---

## Kubernetes

The manifests below provide a minimal production deployment. For a real cluster, replace in-cluster PostgreSQL/Redis/MinIO with managed services (e.g., RDS, ElastiCache, S3).

### 1. Namespace

```yaml
# k8s/namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: open-pm
```

### 2. Secrets

```yaml
# k8s/secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: openpm-secrets
  namespace: open-pm
type: Opaque
stringData:
  DATABASE_URL: "postgres://user:pass@your-db-host:5432/openpm?sslmode=require"
  JWT_SECRET: "change-me-to-a-long-random-string"
  STORAGE_ACCESS_KEY: "your-access-key"
  STORAGE_SECRET_KEY: "your-secret-key"
  SMTP_USER: ""
  SMTP_PASS: ""
```

### 3. ConfigMap

```yaml
# k8s/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: openpm-config
  namespace: open-pm
data:
  PORT: "8080"
  ENVIRONMENT: "production"
  SITE_URL: "https://your-domain.com"
  REDIS_URL: "redis://your-redis-host:6379/0"
  STORAGE_ENDPOINT: "s3.amazonaws.com"
  STORAGE_BUCKET: "open-pm-assets"
  SMTP_HOST: "smtp.your-provider.com"
  SMTP_PORT: "587"
  SMTP_FROM_NAME: "Open-PM"
  SMTP_FROM_ADDR: "noreply@your-domain.com"
```

### 4. Server Deployment

```yaml
# k8s/server-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: openpm-server
  namespace: open-pm
spec:
  replicas: 2
  selector:
    matchLabels:
      app: openpm-server
  template:
    metadata:
      labels:
        app: openpm-server
    spec:
      containers:
        - name: server
          image: open-pm/server:latest
          ports:
            - containerPort: 8080
          envFrom:
            - configMapRef:
                name: openpm-config
            - secretRef:
                name: openpm-secrets
          readinessProbe:
            httpGet:
              path: /api/health
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 10
          livenessProbe:
            httpGet:
              path: /api/health
              port: 8080
            initialDelaySeconds: 15
            periodSeconds: 20
          resources:
            requests:
              cpu: 100m
              memory: 128Mi
            limits:
              cpu: 500m
              memory: 512Mi
---
apiVersion: v1
kind: Service
metadata:
  name: openpm-server
  namespace: open-pm
spec:
  selector:
    app: openpm-server
  ports:
    - port: 8080
      targetPort: 8080
```

### 5. Web Deployment

```yaml
# k8s/web-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: openpm-web
  namespace: open-pm
spec:
  replicas: 2
  selector:
    matchLabels:
      app: openpm-web
  template:
    metadata:
      labels:
        app: openpm-web
    spec:
      containers:
        - name: web
          image: open-pm/web:latest
          ports:
            - containerPort: 80
          resources:
            requests:
              cpu: 50m
              memory: 64Mi
            limits:
              cpu: 200m
              memory: 128Mi
---
apiVersion: v1
kind: Service
metadata:
  name: openpm-web
  namespace: open-pm
spec:
  selector:
    app: openpm-web
  ports:
    - port: 80
      targetPort: 80
```

### 6. Ingress

```yaml
# k8s/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: openpm-ingress
  namespace: open-pm
  annotations:
    nginx.ingress.kubernetes.io/proxy-body-size: "50m"
    cert-manager.io/cluster-issuer: "letsencrypt-prod"  # if using cert-manager
spec:
  ingressClassName: nginx
  tls:
    - hosts:
        - your-domain.com
      secretName: openpm-tls
  rules:
    - host: your-domain.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: openpm-web
                port:
                  number: 80
```

> **Note:** The nginx inside the `web` container proxies `/api/` and `/auth/` to `http://server:8080`. Since `openpm-server` is a Kubernetes Service, this resolves automatically within the same namespace.

### 7. Apply all manifests

```bash
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/secret.yaml
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/server-deployment.yaml
kubectl apply -f k8s/web-deployment.yaml
kubectl apply -f k8s/ingress.yaml
```

### 8. Run migrations in Kubernetes

```bash
kubectl run openpm-migrate \
  --image=open-pm/server:latest \
  --restart=Never \
  --namespace=open-pm \
  --env-from=configmap/openpm-config \
  --env-from=secret/openpm-secrets \
  -- ./server migrate up

# Wait for completion, then clean up
kubectl wait --for=condition=complete pod/openpm-migrate -n open-pm --timeout=120s
kubectl delete pod openpm-migrate -n open-pm
```

---

## Environment Variables

All server configuration is loaded from environment variables (prefix-free, with the `envconfig` library).

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server listen port |
| `ENVIRONMENT` | `development` | `development` or `production` |
| `SITE_URL` | `http://localhost:3000` | Public URL (used in emails) |
| `DATABASE_URL` | *(required)* | PostgreSQL connection string |
| `DATABASE_MAX_OPEN_CONNS` | `25` | Max open DB connections |
| `DATABASE_MAX_IDLE_CONNS` | `5` | Max idle DB connections |
| `JWT_SECRET` | *(required)* | JWT signing secret |
| `JWT_EXP` | `3600` | Access token TTL (seconds) |
| `JWT_REFRESH_EXP` | `604800` | Refresh token TTL (seconds) |
| `REDIS_URL` | `redis://localhost:6379/0` | Redis connection URL |
| `STORAGE_ENDPOINT` | `localhost:9000` | S3/MinIO endpoint (host:port) |
| `STORAGE_ACCESS_KEY` | `minioadmin` | S3 access key |
| `STORAGE_SECRET_KEY` | `minioadmin` | S3 secret key |
| `STORAGE_BUCKET` | `open-pm-assets` | S3 bucket name |
| `SMTP_HOST` | `localhost` | SMTP host |
| `SMTP_PORT` | `1025` | SMTP port |
| `SMTP_USER` | | SMTP username |
| `SMTP_PASS` | | SMTP password |
| `SMTP_FROM_NAME` | `Open-PM` | Sender display name |
| `SMTP_FROM_ADDR` | `noreply@open-pm.dev` | Sender email address |
| `OAUTH_GITHUB_ENABLED` | `false` | Enable GitHub OAuth |
| `OAUTH_GITHUB_CLIENT_ID` | | GitHub OAuth client ID |
| `OAUTH_GITHUB_CLIENT_SECRET` | | GitHub OAuth client secret |
| `OAUTH_GITHUB_REDIRECT_URL` | | GitHub OAuth redirect URL |
| `OAUTH_GOOGLE_ENABLED` | `false` | Enable Google OAuth |
| `OAUTH_GOOGLE_CLIENT_ID` | | Google OAuth client ID |
| `OAUTH_GOOGLE_CLIENT_SECRET` | | Google OAuth client secret |
| `OAUTH_GOOGLE_REDIRECT_URL` | | Google OAuth redirect URL |
| `LLM_ENABLED` | `false` | Enable AI chat assistant |
| `LLM_PROVIDER` | `openai` | LLM provider (`openai`, etc.) |
| `LLM_API_KEY` | | LLM API key |
| `LLM_MODEL` | `gpt-4o-mini` | LLM model name |
| `SEED_ADMIN_EMAIL` | `admin@open-pm.dev` | Seed admin email (dev only) |
| `SEED_ADMIN_PASSWORD` | `password123` | Seed admin password (dev only) |

---

## License

MIT
