# local-cloud-daemon

The daemon writes structured JSON logs to stdout by default. Each entry includes
`time`, `level`, `msg`, and `service`. Docker and systemd can collect these logs
without an application log file. Invalid logging configuration is reported to
stderr and exits with status 1.

| Environment variable | Default | Accepted values |
| --- | --- | --- |
| `LOG_LEVEL` | `info` | `debug`, `info`, `warn`, `error` |
| `LOG_FORMAT` | `json` | `json`, `text` |

Run locally with human-readable logs:

```sh
LOG_LEVEL=debug LOG_FORMAT=text go run . worker
```

Pass `master` to discover workers or `worker` to publish a worker node:

```sh
go run . master
go run . worker
```

Missing or invalid arguments print usage and exit with status 2.

Build the daemon:

```sh
go build -o bin/local-cloud-daemon .
./bin/local-cloud-daemon worker
```

The current daemon scaffold waits for SIGINT or SIGTERM and logs startup and
shutdown. Send SIGINT with Ctrl+C when running locally. In a container, run the
compiled binary as the main process so it receives Docker's shutdown signal.
Read container output with `docker logs --follow <container>`; configure log
retention and rotation in Docker or your system service manager.

Use Go's standard structured logger throughout daemon code:

```go
slog.Info("operation completed", "operation", "sync", "duration", elapsed)
slog.Error("operation failed", "operation", "sync", "error", err)
```

Use debug for diagnostic details, info for lifecycle and normal operations,
warn for recoverable problems, and error for failed operations. Avoid logging
credentials, tokens, or sensitive payloads.
