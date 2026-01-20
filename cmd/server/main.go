package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/singh-anurag-7991/data-guard/internal/api"
	"github.com/singh-anurag-7991/data-guard/internal/engine"
	"github.com/singh-anurag-7991/data-guard/pkg/logger"
)

func main() {
	// Initialize Logger
	logger.Init()
	slog.Info("Starting DataGuard Server...")

	// Initialize Engine
	exec := engine.NewExecutor()

	// Initialize API Handler
	handler := api.NewHandler(exec)

	// Register Routes
	mux := http.NewServeMux()
	mux.HandleFunc("/ingest/api", handler.Ingest)

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
