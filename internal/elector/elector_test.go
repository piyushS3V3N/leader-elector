package elector_test

import (
	"context"
	"testing"
	"time"

	"github.com/yourorg/leader-elector/internal/config"
	"github.com/yourorg/leader-elector/internal/elector"
	"k8s.io/client-go/kubernetes/fake"
)

func TestElector_Run_Simple(t *testing.T) {
	cfg := &config.Config{
		Namespace:     "default",
		LeaseName:     "test-lease",
		Identity:      "test-leader",
		LeaseDuration: 5 * time.Second,
		RenewDeadline: 3 * time.Second,
		RetryPeriod:   1 * time.Second,
	}

	client := fake.NewSimpleClientset()
	e := elector.New(cfg, client)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 🛠 Pass the required callback as second arg to e.Run
	err := e.Run(ctx, func(ctx context.Context) {
		t.Logf("%s: Acting as leader", cfg.Identity)
		<-ctx.Done() // simulate leader work
	})

	if err != nil {
		t.Fatalf("leader election run failed: %v", err)
	}
}
