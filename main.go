package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"raiashpanda007/local-cloud-daemon/internal/logging"
	"raiashpanda007/local-cloud-daemon/internal/utils"
)

func main() {
	// system wide definations

	logger, err := logging.New(os.Stdout, os.Getenv("LOG_LEVEL"), os.Getenv("LOG_FORMAT"))
	_, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err != nil {
		slog.New(slog.NewJSONHandler(os.Stderr, nil)).Error("invalid logging configuration", "error", err)
		os.Exit(1)
	}
	slog.SetDefault(logger)

	sysError := utils.Init()
	go sysError.Handler()

	started := time.Now()
	logger.Info("daemon started", "pid", os.Getpid())

	logger.Info("shutdown requested", "signal")
	logger.Info("daemon stopped", "uptime", time.Since(started).String())
}
