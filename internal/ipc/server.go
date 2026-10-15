package ipc

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"os"
	"sync"

	"raiashpanda007/local-cloud-daemon/internal/utils"
)

type Ipc struct {
	SocketServerStart any
}

func (r *Ipc) SocketServer(centralErrCtx context.Context, sysError utils.SystemError, socketpath string) {
	listener, err := net.Listen("unix", socketpath)

	if err != nil {
		sysError.NewDaemonError(err, "Unable to start the unix socket server "+socketpath, 3)
		return
	}

	var (
		wg    sync.WaitGroup
		mu    sync.Mutex
		conns = make(map[net.Conn]struct{})
	)

	closeListener := sync.OnceFunc(func() {
		_ = listener.Close()
	})
	defer func() {
		if err := os.Remove(socketpath); err != nil && !os.IsNotExist(err) {
			sysError.NewDaemonError(err, "Unable to remove unix socket file "+socketpath, 2)
		}
	}()
	defer closeListener()

	go func() {
		<-centralErrCtx.Done()
		closeListener()
	}()

	slog.Info("Started socket server at path :: ", socketpath)

	for {
		conn, err := listener.Accept()
		if err != nil {
			if centralErrCtx.Err() != nil || errors.Is(err, net.ErrClosed) {
				break
			}

			sysError.NewDaemonError(err, "Unable to accept unix socket connections on "+socketpath, 3)
			closeListener()
			break
		}

		mu.Lock()
		conns[conn] = struct{}{}
		mu.Unlock()

		wg.Add(1)
		go func(c net.Conn) {
			defer wg.Done()
			defer func() {
				mu.Lock()
				delete(conns, c)
				mu.Unlock()
				_ = c.Close()
			}()

			r.handleConnection(centralErrCtx, c, sysError)
		}(conn)
	}

	mu.Lock()
	for c := range conns {
		_ = c.Close()
	}
	mu.Unlock()

	wg.Wait()
	slog.Info("Socket server stopped at path :: ", socketpath)
}

func (r *Ipc) handleConnection(ctx context.Context, conn net.Conn, sysError utils.SystemError) {
	slog.Info("Client connected")

	buf := make([]byte, 4*1024)

	for {
		if ctx.Err() != nil {
			return
		}

		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				slog.Info("Client disconnected smoothly")
			} else if ctx.Err() != nil || errors.Is(err, net.ErrClosed) {
				return
			} else {
				sysError.NewDaemonError(err, "Unable to read from unix socket connection", 1)
			}
			return
		}

		receivedData := buf[:n]
		slog.Info("Received :: ", string(receivedData))
	}
}
