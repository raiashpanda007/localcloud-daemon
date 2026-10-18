package daemon

import (
	"context"
	"sync"

	"raiashpanda007/local-cloud-daemon/internal/ipc"
	"raiashpanda007/local-cloud-daemon/internal/utils"
	worker_network "raiashpanda007/local-cloud-daemon/internal/worker/network"
)

const (
	SERVICE_TYPE  = "_lcworker._tcp"
	WORKER_DOMAIN = "local."

	// Shared by master and worker: the IPC socket the CLI talks to.
	SOCKET_SERVER_PATH = "/tmp/daemon.sock"
)

func WorkerNode(rootCtx context.Context, sysError utils.SystemError, state *State) {

	var socketDone sync.WaitGroup
	socketDone.Add(1)
	go func() {
		defer socketDone.Done()
		_ipc := ipc.Ipc{}
		_ipc.SocketServer(rootCtx, sysError, SOCKET_SERVER_PATH)
	}()

	worker_network.PublishWorker(rootCtx, sysError, SERVICE_TYPE, WORKER_DOMAIN)

	<-rootCtx.Done()
	socketDone.Wait()
}
