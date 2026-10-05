package master_network

import (
	"context"
	"log/slog"
	"raiashpanda007/local-cloud-daemon/internal/utils"
	"time"

	"github.com/grandcat/zeroconf"
)

func DiscoverWorkers(serviceProtocol, domain string, rootCtx context.Context, sysError utils.SystemError) {
	resolver, err := zeroconf.NewResolver(nil)

	if err != nil {
		sysError.NewDaemonError(err, "Can't create mDNS resolver to discover workers", 3)
		return
	}

	entriesRecieved := make(chan *zeroconf.ServiceEntry)

	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			slog.Info("Discovered worker",
				"instance", entry.Instance,
				"hostname", entry.HostName,
				"ipv4", entry.AddrIPv4,
				"ipv6", entry.AddrIPv6,
				"port", entry.Port,
				"metadata", entry.Text,
			)
		}
	}(entriesRecieved)

	ctx, cancel := context.WithTimeout(rootCtx, time.Second*5)
	defer cancel()

	slog.Info("Starting worker discovery", "service", serviceProtocol, "domain", domain)
	err = resolver.Browse(ctx, serviceProtocol, domain, entriesRecieved)
	if err != nil {
		sysError.NewDaemonError(err, "Unable to browse for workers over mDNS", 2)
		return
	}
	select {
	case <-ctx.Done():
		slog.Info("Worker discovery stopped", "service", serviceProtocol, "domain", domain, "reason", ctx.Err())
		return
	case <-rootCtx.Done():
		return
	}

}
