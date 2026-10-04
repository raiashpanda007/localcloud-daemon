package worker_network

import (
	"context"
	"github.com/google/uuid"
	"github.com/grandcat/zeroconf"
	"log/slog"
	"net"
	"raiashpanda007/local-cloud-daemon/internal/utils"
)

func PublishWorker(rootCtx context.Context, sysError utils.SystemError, SERVICE_TYPE, WORKER_DOMAIN string) {

	var SERVICE_NAME = "worker"
	var METADATA = []string{"1.0.0", "local.cloud"}

	SERVICE_NAME += "-" + uuid.New().String()[:8]

	listener, err := utils.GetRandomPort()

	if err != nil {
		sysError.NewDaemonError(err, "Unable to get random port to start worker server", 3)
		return
	}

	defer listener.Close()
	port := listener.Addr().(*net.TCPAddr).Port

	go StartServer(rootCtx, sysError, listener)

	mDNS, err := zeroconf.Register(SERVICE_NAME, SERVICE_TYPE, WORKER_DOMAIN, port, METADATA, nil)

	slog.Info("mDNS server restarted and registered", SERVICE_NAME, " :: ", SERVICE_TYPE)

	if err != nil {
		sysError.NewDaemonError(err, "Unable to register worker over mDNS", 2)
		return
	}

	<-rootCtx.Done()

	defer mDNS.Shutdown()
}
