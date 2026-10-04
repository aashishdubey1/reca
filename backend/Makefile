.PHONY: dev test docker-up docker-down build

# Start Postgres and Redis
docker-up:
	docker compose up -d

# Stop Postgres and Redis
docker-down:
	docker compose down

# Run tests
test:
	go test -v ./...

# Build binary
build:
	go build -o bin/server ./cmd/server

# Run backend development server
dev:
	DATABASE_URL="postgres://postgres:postgrespassword@localhost:5432/reca?sslmode=disable" \
	REDIS_URL="localhost:6379" \
	PORT="8080" \
	go run ./cmd/server
