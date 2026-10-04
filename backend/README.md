# Reca Backend (Go + PostgreSQL + Redis)

A high-performance Spaced Repetition engine and API built with Go, PostgreSQL, and Redis.

## Directory Structure & Code Navigation

- **`internal/srs/engine.go`**: Core mathematical implementation of the **SM-2 (SuperMemo-2)** spaced repetition algorithm (interval multipliers, ease factor adjustment, lapse handling).
- **`internal/srs/engine_test.go`**: Unit tests verifying ease factor formulas, repetition progression, and reset on failed recall.
- **`internal/model/models.go`**: Data models for `Topic`, `SRSItem`, `ReviewLog`, `TimelineDay`, and API DTOs.
- **`internal/repository/postgres.go`**: PostgreSQL connection pool (pgx) with automated schema migrations and indexing.
- **`internal/repository/topic_repo.go`**: Database queries for creating topics, fetching due revisions, atomically executing reviews with row locks (`FOR UPDATE`), and generating timeline data.
- **`internal/repository/redis.go`**: Redis caching for timeline day summaries and daily streak bitfields.
- **`internal/handler/router.go`**: Chi router setup with CORS support for the Next.js frontend (`http://localhost:3000`).
- **`internal/handler/topic_handler.go`**: REST HTTP endpoints for topics, reviews, and timeline.
- **`cmd/server/main.go`**: Application entry point with configuration loading and graceful shutdown.
- **`docker-compose.yml`**: Pre-configured PostgreSQL 16 and Redis 7 containers.

---

## Quick Start

### 1. Start PostgreSQL & Redis
```bash
cd backend
docker compose up -d
```

### 2. Run Tests
```bash
go test -v ./...
```

### 3. Run Development Server
```bash
make dev
# or
go run ./cmd/server
```

The API will be live at `http://localhost:8080`.
