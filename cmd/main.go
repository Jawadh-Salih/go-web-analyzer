package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/Jawadh-Salih/go-web-analyzer/internal/logger"
	"github.com/Jawadh-Salih/go-web-analyzer/internal/server"
)

func main() {
	// enable logging and inject here.
	logger := logger.New()

	svr := server.New(":8080", logger, true)
	logger.Info("Starting server...", "addr", ":8080")

	go func() {
		if err := svr.Start(); err != nil {
			logger.Error("Server failed", "error", err)
			os.Exit(1)
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)
	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := svr.Stop(ctx); err != nil {
		logger.Error("Server shutdown failed", "error", err)
		os.Exit(1)
	}
	logger.Info("Server stopped gracefully")

}
