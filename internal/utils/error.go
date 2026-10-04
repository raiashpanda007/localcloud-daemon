package utils

import (
	"log/slog"
	"os"
	"runtime/debug"
)

/*
1. Invalid data sort off.
2. Service not responded (Not a big problem to kil the daemon)
3. Not solvable have to kill the daemon
*/
type DaemonError struct {
	Error   error
	Trace   string
	Message string
	Code    int
}

type SystemError struct {
	rootErrorChan chan DaemonError
}

func Init() SystemError {
	errChan := make(chan DaemonError)

	return SystemError{
		rootErrorChan: errChan,
	}
}

func (r *SystemError) NewDaemonError(err error, message string, code int) {
	var errRecieved = DaemonError{
		Error:   err,
		Message: message,
		Trace:   string(debug.Stack()),
		Code:    code,
	}

	slog.Error("ERROR RECIEVED ::", "message", message, "error", err, "trace", errRecieved.Trace)

	r.rootErrorChan <- errRecieved
}

func (r *SystemError) Handler() {
	for {
		daemonError := <-r.rootErrorChan

		if daemonError.Code == 1 {
			slog.Warn("Invalid data received; continuing daemon", "message", daemonError.Message, "error", daemonError.Error)
			continue
		}

		if daemonError.Code == 2 {
			slog.Warn("Service did not respond; continuing daemon", "message", daemonError.Message, "error", daemonError.Error)
			continue
		}

		if daemonError.Code == 3 {
			slog.Error("Can't resolve it killing the daemon. Please restart the daemon", "message", daemonError.Message, "error", daemonError.Error, "trace", daemonError.Trace)
			os.Exit(1)
		}
	}
}
