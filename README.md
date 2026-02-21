# Open-PM

Open-source project management platform. Manage tasks, sprints, docs, and triage.

Built with Go, Vue.js 3, and Tailwind CSS.

## Quick Start

```bash
# Start infrastructure (PostgreSQL, Redis, MinIO, MailHog)
make dev

# Run database migrations
make db-migrate

# Generate sqlc code
make sqlc-generate

# Start backend (in one terminal)
make dev-server

# Start frontend (in another terminal)
make dev-web
```

## Stack

- **Backend**: Go, Chi, PostgreSQL, Redis, sqlc
- **Frontend**: Vue 3, TypeScript, Tailwind CSS, Pinia
- **Infrastructure**: Docker Compose, MinIO (S3), MailHog

## Services (Development)

| Service | URL |
|---------|-----|
| Frontend | http://localhost:3000 |
| Backend API | http://localhost:8080 |
| MailHog UI | http://localhost:8025 |
| MinIO Console | http://localhost:9001 |

## License

MIT
