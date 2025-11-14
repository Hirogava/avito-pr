// Package main is the entry point of the Avito-PR backend server.
package main

import (
	"net/http"
	"os"
	"time"

	"github.com/Hirogava/avito-pr/internal/config/environment"
	"github.com/Hirogava/avito-pr/internal/config/logger"
	postgres "github.com/Hirogava/avito-pr/internal/repository/postgres"
	"github.com/Hirogava/avito-pr/internal/service/shoutdown"
	router "github.com/Hirogava/avito-pr/internal/transport/http"
)

func main() {
	if err := environment.LoadEnvFile(".env"); err != nil {
		logger.Logger.Fatalf("failed to load .env: %v", err)
	}

	logger.LogInit()
	logger.Logger.Info("Starting Avito-PR backend server")

	dbConnStr := os.Getenv("DB_CONNECT_STRING")
	if dbConnStr == "" {
		logger.Logger.Fatal("DB_CONNECT_STRING environment variable is required")
	}
	logger.Logger.Info("Connecting to database", "connection_string", dbConnStr)

	manager := postgres.NewManager("postgres", dbConnStr)
	logger.Logger.Info("Database connection established successfully")
	defer manager.Close()

	logger.Logger.Info("Running database migrations")
	manager.Migrate()
	logger.Logger.Info("Database migrations completed successfully")

	logger.Logger.Info("Initializing HTTP router")
	r := router.CreateRouter(manager)

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = ":8080"
		logger.Logger.Warn("SERVER_PORT not set, using default port 8080")
	}

	server := &http.Server{
		Addr:         serverPort,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	logger.Logger.Info("Starting HTTP server", "port", serverPort)
	shoutdown.Graceful(server, 30*time.Second)
}
