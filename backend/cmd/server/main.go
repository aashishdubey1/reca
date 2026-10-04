package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"reca-backend/internal/config"
	"reca-backend/internal/handler"
	"reca-backend/internal/repository"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("[Server] Initializing reca backend...")

	// 1. Connect to PostgreSQL
	pgDB, err := repository.NewPostgresDB(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("[Postgres] Failed to initialize database: %v\nMake sure Postgres is running (e.g., via `docker compose up -d`)", err)
	}
	defer pgDB.Close()

	// 2. Connect to Redis
	redisClient := repository.NewRedisClient(cfg.RedisURL)

	// 3. Initialize Repositories and Handlers
	topicRepo := repository.NewTopicRepository(pgDB)
	topicHandler := handler.NewTopicHandler(topicRepo, redisClient)

	// 4. Initialize Router
	router := handler.NewRouter(topicHandler)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 5. Run Server with Graceful Shutdown
	go func() {
		log.Printf("[Server] Server listening on http://localhost:%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[Server] Listen error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[Server] Shutting down gracefully...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("[Server] Forced shutdown: %v", err)
	}

	log.Println("[Server] Exited cleanly")
}
