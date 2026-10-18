package daemon

import (
	"context"
	"sync"

	"raiashpanda007/local-cloud-daemon/internal/ipc"
	master_network "raiashpanda007/local-cloud-daemon/internal/master/network"
	"raiashpanda007/local-cloud-daemon/internal/utils"
)

func MasterNode(rootCtx context.Context, sysError utils.SystemError, state *State) {

	var socketDone sync.WaitGroup
	socketDone.Add(1)
	go func() {
		defer socketDone.Done()
		_ipc := ipc.Ipc{}
		_ipc.SocketServer(rootCtx, sysError, SOCKET_SERVER_PATH)
	}()

	master_network.DiscoverWorkers(SERVICE_TYPE, WORKER_DOMAIN, rootCtx, sysError)

	<-rootCtx.Done()
	socketDone.Wait()
}
