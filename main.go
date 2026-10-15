package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"raiashpanda007/local-cloud-daemon/internal/daemon"
	"raiashpanda007/local-cloud-daemon/internal/ipc"
	"raiashpanda007/local-cloud-daemon/internal/logging"
	"raiashpanda007/local-cloud-daemon/internal/utils"
)

const (
	SOCKET_SERVER_PATH string = "/tmp/daemon.sock"
)

func main() {
	if len(os.Args) != 2 || (os.Args[1] != "master" && os.Args[1] != "worker") {
		fmt.Fprintf(os.Stderr, "Usage: %s <master|worker>\n", os.Args[0])
		os.Exit(2)
	}
	role := os.Args[1]

	// system wide definations

	logger, err := logging.New(os.Stdout, os.Getenv("LOG_LEVEL"), os.Getenv("LOG_FORMAT"))
	rootCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err != nil {
		slog.New(slog.NewJSONHandler(os.Stderr, nil)).Error("invalid logging configuration", "error", err)
		os.Exit(1)
	}
	slog.SetDefault(logger)

	sysError := utils.Init()
	go sysError.Handler()

	started := time.Now()
	logger.Info("daemon started", "pid", os.Getpid(), "role", role)

	_ipc := ipc.Ipc{}
	var socketDone sync.WaitGroup
	socketDone.Add(1)
	go func() {
		defer socketDone.Done()
		_ipc.SocketServer(rootCtx, sysError, SOCKET_SERVER_PATH)
	}()

	switch role {
	case "master":
		daemon.MasterNode(rootCtx, sysError)
	case "worker":
		daemon.WorkerNode(rootCtx, sysError)
	}

	stop()
	socketDone.Wait()

	logger.Info("shutdown requested")
	logger.Info("daemon stopped", "uptime", time.Since(started).String())
}
