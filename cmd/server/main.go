package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/singh-anurag-7991/data-guard/internal/api"
	"github.com/singh-anurag-7991/data-guard/internal/engine"
	"github.com/singh-anurag-7991/data-guard/internal/ingest/postgres"
	"github.com/singh-anurag-7991/data-guard/internal/storage"
	"github.com/singh-anurag-7991/data-guard/pkg/logger"
)

func main() {
	// Initialize Logger
	logger.Init()
	slog.Info("Starting DataGuard Server...")

	// Initialize DB (Postgres)
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		slog.Warn("DATABASE_URL not set, storage disabled")
	}

	// Create context for DB connection
	ctx := context.Background() // basic root context
	var repo *storage.Repository

	if dbURL != "" {
		pgClient, err := postgres.NewClient(ctx, dbURL)
		if err != nil {
			slog.Error("Failed to connect to DB", "error", err)
			os.Exit(1)
		}
		defer pgClient.Close()
		repo = storage.NewRepository(pgClient)
	}

	// Initialize Engine
	exec := engine.NewExecutor()

	// Initialize API Handlers
	ingestHandler := api.NewHandler(exec)
	dashboardHandler := api.NewDashboardHandler(repo)

	// Register Routes
	mux := http.NewServeMux()
	mux.HandleFunc("/ingest/api", ingestHandler.Ingest)

	if repo != nil {
		mux.HandleFunc("/api/runs", dashboardHandler.ListRuns)
	}

	// Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	slog.Info("Server listening", "port", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	}
}
