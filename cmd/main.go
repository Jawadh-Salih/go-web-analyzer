package main

import (
	"github.com/Jawadh-Salih/go-web-analyzer/internal/logger"
	"github.com/Jawadh-Salih/go-web-analyzer/internal/server"
)

func main() {
	// enable logging and inject here.
	logger := logger.New()

	svr := server.New(":8080", logger)
	logger.Info("Starting server...", "addr", ":8080")

	if err := svr.Start(); err != nil {
		logger.Error("Server failed", "error", err)
	}
}
