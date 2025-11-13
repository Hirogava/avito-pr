package shoutdown

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Hirogava/avito-pr/internal/config/logger"
)

func Graceful(server *http.Server, timeout time.Duration) {
    serverErr := make(chan error, 1)
    
    go func() {
        logger.Logger.Info("Starting HTTP server", "port", server.Addr)
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            serverErr <- err
        }
    }()
    
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    
    select {
    case err := <-serverErr:
        logger.Logger.Fatal("Server error", "error", err)
    case sig := <-quit:
        logger.Logger.Info("Received signal, shutting down", "signal", sig)
        
        ctx, cancel := context.WithTimeout(context.Background(), timeout)
        defer cancel()
        
        if err := server.Shutdown(ctx); err != nil {
            logger.Logger.Error("Graceful shutdown failed", "error", err)
            if err := server.Close(); err != nil {
                logger.Logger.Error("Force shutdown failed", "error", err)
            }
        }
        
        logger.Logger.Info("Server stopped gracefully")
    }
}
