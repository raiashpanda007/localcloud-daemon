# TODO

Basic things to build, one by one, to complete the project.
(Future features go in [FEATURE_REQUESTS.md](FEATURE_REQUESTS.md), not here.)

## IPC: support `localcloud master start` / `worker start`

The CLI side is done (`localcloud-cli/internals/networks`). Protocol: one JSON
object per line (`json.Marshal` + `\n`), one response per request.

```go
type Request  struct { Type string `json:"type"` }            // "status"
type Response struct {
	OK    bool   `json:"ok"`
	Role  string `json:"role,omitempty"`   // "master" | "worker"
	PID   int    `json:"pid,omitempty"`
	Error string `json:"error,omitempty"`
}
```

- [X] **Answer `status` on the IPC socket.** Pass the role into the socket
      server. Read requests line by line (`bufio`) and reply
      `{"ok":true,"role":"<role>","pid":<os.Getpid()>}`. Reply to an unknown
      type with `{"ok":false,"error":"unknown request type"}`. The CLI gives up
      after 2 s without a reply.
- [ ] **Remove a stale `/tmp/daemon.sock` before `Listen`** (ignore
      not-exist). Code 3 calls `os.Exit` without running defers, so a crash
      leaves the file behind and the next start fails with "address already in
      use".
- [x] **Fix the `internal/daemon/master.go` build break.** The duplicate
      socket server was removed; `main.go` runs the only one.
- [ ] **Keep the master alive until `rootCtx` is done.** `DiscoverWorkers` returns
      after its 5 s browse, so `MasterNode` returns and the daemon shuts down.
