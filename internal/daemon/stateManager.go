package daemon

import (
	"context"
	"os"
)

type ROLE string

const (
	MASTER ROLE = "master"
	WORKER ROLE = "worker"
)

type StateResponse struct {
	Role   ROLE   `json:"role"`
	Status string `json:"status"`
	PID    int    `json:"pid"`
}

// command is one change or query on the state. apply runs only inside State.Run,
// so the fields are never touched by two goroutines at once.
type command interface {
	apply(s *State)
}

type State struct {
	role     ROLE
	status   string // TODO: Right now keeping status as string later on add enum.
	pid      int
	commands chan command
}

func InitState(role ROLE) *State {
	return &State{
		role:     role,
		status:   "INITIALIZED",
		pid:      os.Getpid(),
		commands: make(chan command, 16),
	}
}

// Run owns the state. Start it once with `go state.Run(rootCtx)`.
func (s *State) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case cmd := <-s.commands:
			cmd.apply(s)
		}
	}
}

// send queues a command for Run. It gives up if ctx is cancelled first.
func (s *State) send(ctx context.Context, cmd command) error {
	select {
	case s.commands <- cmd:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// GetRole is safe without the loop: the role is set once in InitState and never changes.
func (s *State) GetRole() ROLE {
	return s.role
}

// UpdateStatus waits until Run has applied the change, so the caller knows it took effect.
func (s *State) UpdateStatus(ctx context.Context, status string) error {
	return setField(ctx, s, func(st *State) *string { return &st.status }, status)
}

// setField sends a change to one field of the state and waits until Run has applied it.
// field picks the field to write, e.g. func(st *State) *string { return &st.status }.
func setField[T any](ctx context.Context, s *State, field func(*State) *T, value T) error {
	reply := make(chan error, 1)
	if err := s.send(ctx, set[T]{field: field, value: value, reply: reply}); err != nil {
		return err
	}
	select {
	case err := <-reply:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// GetInfo returns a copy of the state, read by the loop.
func (s *State) GetInfo(ctx context.Context) (StateResponse, error) {
	reply := make(chan StateResponse, 1)
	if err := s.send(ctx, getInfo{reply: reply}); err != nil {
		return StateResponse{}, err
	}
	select {
	case info := <-reply:
		return info, nil
	case <-ctx.Done():
		return StateResponse{}, ctx.Err()
	}
}

// --- commands: one type per domain. To add a domain, add a type and its apply method. ---

// set writes value into the field picked by field. T is the field's type, so one
// type covers every setter.
type set[T any] struct {
	field func(*State) *T
	value T
	reply chan error
}

func (c set[T]) apply(s *State) {
	*c.field(s) = c.value
	c.reply <- nil
}

type getInfo struct {
	reply chan StateResponse
}

func (c getInfo) apply(s *State) {
	c.reply <- StateResponse{
		Role:   s.role,
		Status: s.status,
		PID:    s.pid,
	}
}
