package daemon

import (
	"context"
	master_network "raiashpanda007/local-cloud-daemon/internal/master/network"
	"raiashpanda007/local-cloud-daemon/internal/utils"
)

func MasterNode(rootCtx context.Context, sysError utils.SystemError) {

	master_network.DiscoverWorkers(SERVICE_TYPE, WORKER_DOMAIN, rootCtx, sysError)

}
