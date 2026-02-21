.PHONY: dev dev-server dev-web setup db-up db-down db-migrate db-migrate-down sqlc-generate lint test build

# Start infrastructure services
dev:
	docker compose up -d postgres redis minio mailhog minio-init
	@echo "Infrastructure started:"
	@echo "  PostgreSQL: localhost:5432"
	@echo "  Redis:      localhost:6379"
	@echo "  MinIO:      localhost:9000 (console: localhost:9001)"
	@echo "  MailHog:    localhost:8025"
	@echo ""
	@echo "Run 'make dev-server' and 'make dev-web' in separate terminals."

dev-server:
	cd server && go run ./cmd/server/

dev-web:
	cd web && npm run dev

# Setup from scratch
setup: dev db-migrate sqlc-generate
	cd web && npm install
	@echo "Setup complete!"

# Database
db-up:
	docker compose up -d postgres

db-down:
	docker compose down

db-migrate:
	cd server && go run ./cmd/migrate/ up

db-migrate-down:
	cd server && go run ./cmd/migrate/ down 1

db-migrate-create:
	@read -p "Migration name: " name; \
	cd server && go run github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
		create -ext sql -dir ./migrations -seq $$name

# Code generation
sqlc-generate:
	cd server && go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate

# Quality
lint:
	cd server && golangci-lint run ./...
	cd web && npm run lint

test:
	cd server && go test ./... -v -race
	cd web && npm run test

# Build
build-server:
	cd server && go build -o ./tmp/server ./cmd/server/

build-web:
	cd web && npm run build

build: build-server build-web
