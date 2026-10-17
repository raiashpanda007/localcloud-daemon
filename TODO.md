# TODO

Basic things to build, one by one, to complete the project.
(Future features go in [FEATURE_REQUESTS.md](FEATURE_REQUESTS.md), not here.)

## Daemon state

- [x] **Role constants.** `MASTER` / `WORKER` in `internal/daemon/stateManager.go`.
- [x] **State holds role, status and PID.** `InitState` stores `os.Getpid()`.
- [x] **Create the state once in `main.go`**, before the socket server starts, and pass the same pointer to the socket server and the role function. Today it's created inside the role `switch`, after the socket server has started.
- [ ] **Only the main routine updates the state.** Other goroutines read it by sending a request over a channel that main owns (request + reply channel). No shared mutable fields.
- [ ] **Add a central message bus owned by main.** Socket requests are published onto it, and main handles them against the state and replies on each request's reply channel. This is the single path into the state, so the socket server never touches the state directly.

## IPC protocol and socket server

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

- [ ] **Build the protocol in `internal/ipc/protocol.go`.** Request/response types, a size limit on request lines, and a single-line encoder. Keep the structs in sync with the CLI copy.
- [ ] **Build the request handling in `internal/ipc/handler.go`.** Read line by line, decode the request, and answer `status` with `{"ok":true,"role":"<role>","pid":<pid>}`. Reply to an unknown type with `{"ok":false,"error":"unknown request type"}`. The CLI gives up after 2 s without a reply.
- [ ] **Connect the handler to `server.go`.** `handleConnection` currently only logs raw bytes. It needs to call the handler and reply on the same connection. Status should go through the message bus above, not read the state directly.
- [ ] **Remove a stale socket before `Listen`** (ignore not-exist). Code 3 calls `os.Exit` without running defers, so a crash leaves the file behind and the next start fails with "address already in use".
- [ ] **Make the socket path configurable.** Needed so master and worker can run on the same machine in dev mode. The CLI has to know the path too, so it's a shared decision (see the CLI notes in `localcloud-cli/README.md`).

## Daemon lifecycle

- [x] **Fix the `internal/daemon/master.go` build break.** The duplicate socket server was removed; `main.go` runs the only one.
- [ ] **Keep the master alive until `rootCtx` is done.** `DiscoverWorkers` returns after its 5 s browse, so `MasterNode` returns and the daemon shuts down.
