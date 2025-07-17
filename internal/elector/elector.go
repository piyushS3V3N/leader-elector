package elector

import (
	"context"

	"github.com/yourorg/leader-elector/internal/config"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/klog/v2"
)

type Elector struct {
	cfg            *config.Config
	client         kubernetes.Interface
	leadershipLost chan struct{}
}

// New creates a new Elector instance with an optional buffered leadershipLost channel.
func New(cfg *config.Config, client kubernetes.Interface) *Elector {
	return &Elector{
		cfg:            cfg,
		client:         client,
		leadershipLost: make(chan struct{}, 1), // buffered to avoid blocking
	}
}

// SetLeadershipLostChannel lets you replace the internal leadershipLost channel with your own.
func (e *Elector) SetLeadershipLostChannel(ch chan struct{}) {
	e.leadershipLost = ch
}

// Run starts the leader election loop.
// The onStartedLeading callback is called when this instance acquires leadership.
func (e *Elector) Run(ctx context.Context, onStartedLeading func(ctx context.Context)) error {
	lock := &resourcelock.LeaseLock{
		LeaseMeta: v1.ObjectMeta{
			Name:      e.cfg.LeaseName,
			Namespace: e.cfg.Namespace,
		},
		Client: e.client.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: e.cfg.Identity,
		},
	}

	klog.Infof("Starting leader election as %s", e.cfg.Identity)

	leaderelectionCfg := leaderelection.LeaderElectionConfig{
		Lock:            lock,
		ReleaseOnCancel: true,
		LeaseDuration:   e.cfg.LeaseDuration,
		RenewDeadline:   e.cfg.RenewDeadline,
		RetryPeriod:     e.cfg.RetryPeriod,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: onStartedLeading,
			OnStoppedLeading: func() {
				klog.Warningf("%s: lost leadership", e.cfg.Identity)
				select {
				case e.leadershipLost <- struct{}{}:
				default:
					// avoid blocking if no listener is ready
				}
				// Optionally exit or handle shutdown here
				//	os.Exit(1)
			},
			OnNewLeader: func(identity string) {
				if identity != e.cfg.Identity {
					klog.Infof("%s: new leader is %s", e.cfg.Identity, identity)
				}
			},
		},
	}

	leaderelection.RunOrDie(ctx, leaderelectionCfg)
	return nil
}

