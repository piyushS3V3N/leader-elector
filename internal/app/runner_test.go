package app

import (
	"context"
	"testing"
	"time"
)

func TestRun_CancelContextStopsExecution(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(1 * time.Second)
		cancel()
	}()

	done := make(chan struct{})
	go func() {
		Run(ctx, "test-leader")
		done <- struct{}{}
	}()

	select {
	case <-done:
		// Test passed
	case <-time.After(3 * time.Second):
		t.Error("Run() did not exit after context cancellation")
	}
}
