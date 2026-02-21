package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/open-pm/open-pm/server/internal/api"
	"github.com/open-pm/open-pm/server/internal/config"
	"github.com/open-pm/open-pm/server/internal/db"
	"github.com/open-pm/open-pm/server/internal/repository"
	"github.com/open-pm/open-pm/server/internal/seed"
)

func main() {
	// Pretty logging for development
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	// Load .env file (optional, ignores if not found)
	_ = godotenv.Load()
	_ = godotenv.Load("../.env")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load configuration")
	}

	// Connect to database
	database, err := db.New(cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer database.Close()

	// Create repository (implements api.Queries)
	repo := repository.New(database.Pool)

	// Seed default admin account and workspace
	if cfg.Environment == "development" {
		if err := seed.Run(context.Background(), cfg, repo); err != nil {
			log.Error().Err(err).Msg("failed to seed default data")
		}
	}

	// Create API
	a := api.NewAPI(cfg, repo)

	// Start server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      a.Handler(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info().Msgf("starting Open-PM server on :%d", cfg.Port)
		log.Info().Msgf("environment: %s", cfg.Environment)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server error")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("server forced to shutdown")
	}

	log.Info().Msg("server stopped")
}
