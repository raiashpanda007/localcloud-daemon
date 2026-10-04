package daemon

import (
	"context"
	"log/slog"
	"net"
	"raiashpanda007/local-cloud-daemon/internal/utils"
	worker_network "raiashpanda007/local-cloud-daemon/internal/worker/network"

	"github.com/google/uuid"
	"github.com/grandcat/zeroconf"
)

// TODO: 1. Start a server
// TODO: 2. Register over mDNS

const (
	SERVICE_TYPE  = "_lcworker._tcp"
	WORKER_DOMAIN = "local."
)

func WorkerNode(rootCtx context.Context, sysError utils.SystemError) {

	var SERVICE_NAME = "worker"
	var METADATA = []string{"1.0.0", "local.cloud"}

	SERVICE_NAME += "-" + uuid.New().String()[:8]

	listener, err := utils.GetRandomPort()

	if err != nil {
		sysError.NewDaemonError(err, "Unable to get random port to start worker server", 3)
	}

	defer listener.Close()
	port := listener.Addr().(*net.TCPAddr).Port

	go worker_network.StartServer(rootCtx, sysError, listener)

	mDNS, err := zeroconf.Register(SERVICE_NAME, SERVICE_TYPE, WORKER_DOMAIN, port, METADATA, nil)

	slog.Info("mDNS server restarted and registered", SERVICE_NAME, " :: ", SERVICE_TYPE)

	if err != nil {
		sysError.NewDaemonError(err, "Unable to register worker over mDNS", 2)
	}

	defer mDNS.Shutdown()
}
