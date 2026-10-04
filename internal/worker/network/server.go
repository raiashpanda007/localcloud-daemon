package worker_network

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"raiashpanda007/local-cloud-daemon/internal/utils"
)

func StartServer(rootCtx context.Context, sysError utils.SystemError, listener net.Listener) {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthy", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := &http.Server{
		Handler: mux,
	}

	go func() {
		err := server.Serve(listener)
		if err != nil {
			sysError.NewDaemonError(err, "Cannnot serve the worker server", 3)
		}
		slog.Info("Started worker server at port :: ", string(listener.Addr().(*net.TCPAddr).Port))

	}()

	serverShutdownCtx, cancel := context.WithCancel(context.Background())

	defer cancel()
	<-rootCtx.Done()

	err := server.Shutdown(serverShutdownCtx)

	if err != nil {
		sysError.NewDaemonError(err, "Unable to shutdown worker server", 3)
	}

}
