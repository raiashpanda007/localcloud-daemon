package daemon

import (
	"context"

	"raiashpanda007/local-cloud-daemon/internal/utils"
	worker_network "raiashpanda007/local-cloud-daemon/internal/worker/network"
)

const (
	SERVICE_TYPE  = "_lcworker._tcp"
	WORKER_DOMAIN = "local."
)

func WorkerNode(rootCtx context.Context, sysError utils.SystemError) {

	worker_network.PublishWorker(rootCtx, sysError, SERVICE_TYPE, WORKER_DOMAIN)
}
