package daemon

import (
	"context"
	"testing"
	"time"
)

func TestStateUpdateAndGetInfo(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := InitState(MASTER)
	go s.Run(ctx)

	info, err := s.GetInfo(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if info.Role != MASTER || info.Status != "INITIALIZED" || info.PID == 0 {
		t.Fatalf("initial info = %+v", info)
	}

	if err := s.UpdateStatus(ctx, "RUNNING"); err != nil {
		t.Fatal(err)
	}

	info, err = s.GetInfo(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if info.Status != "RUNNING" {
		t.Fatalf("status after update = %q, want RUNNING", info.Status)
	}
}

func TestStateCallsReturnAfterShutdown(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	s := InitState(WORKER)
	go s.Run(ctx)
	cancel() // Run has stopped; nothing will read the channel.

	done := make(chan struct{})
	go func() {
		_ = s.UpdateStatus(ctx, "STOPPING")
		_, _ = s.GetInfo(ctx)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("state calls hung after shutdown")
	}
}

func TestGetRoleWithoutRun(t *testing.T) {
	if got := InitState(WORKER).GetRole(); got != WORKER {
		t.Fatalf("GetRole = %q, want %q", got, WORKER)
	}
}
