package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartServer(ctx context.Context, serverAddress string, readTimeout, writeTimeout, idleTimeout time.Duration) {
	server := &http.Server{
		Addr:         serverAddress,
		ReadTimeout:  readTimeout * time.Second,
		WriteTimeout: writeTimeout * time.Second,
		IdleTimeout:  idleTimeout * time.Second,
	}

	go handleShutdown(ctx, server)

	slog.Info("Server started successfully", "address", serverAddress, "readTimeout", readTimeout, "writeTimeout", writeTimeout, "idleTimeout", idleTimeout)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
	}

}

func handleShutdown(ctx context.Context, server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	slog.Info("Received signal. Initiating graceful shutdown", "signal", sig)

	cctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := server.Shutdown(cctx); err != nil {
		slog.Error("Error during server shutdown", "error", err)
	} else {
		slog.Info("Graceful shutdown completed", "signal", sig)
	}
}
