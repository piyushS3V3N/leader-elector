package main

import (
	"context"
	"flag"
	"github.com/yourorg/leader-elector/internal/config"
	"github.com/yourorg/leader-elector/internal/elector"
	"github.com/yourorg/leader-elector/pkg/k8s"
	"k8s.io/klog/v2"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	klog.InitFlags(nil)
	cfg := config.LoadFromEnv()

	// Flags override env (but only in main)
	flag.StringVar(&cfg.Namespace, "namespace", cfg.Namespace, "Kubernetes namespace")
	flag.StringVar(&cfg.LeaseName, "lease-name", cfg.LeaseName, "Lease name")
	flag.StringVar(&cfg.KubeconfigPath, "kubeconfig", cfg.KubeconfigPath, "Path to kubeconfig file")
	flag.StringVar(&cfg.Identity, "identity", cfg.Identity, "Identity of this instance")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle SIGTERM / SIGINT for graceful shutdown
	go handleSignals(cancel)

	clientset, err := k8s.NewClient(cfg)
	if err != nil {
		klog.Fatalf("Failed to create Kubernetes client: %v", err)
	}

	leaderElector := elector.New(cfg, clientset)
	if err := leaderElector.Run(ctx, func(ctx context.Context) {
		klog.Infof("%s: Acting as leader", cfg.Identity)
		<-ctx.Done()
	}); err != nil {
		klog.Fatalf("Leader election failed: %v", err)
	}
}

func handleSignals(cancel context.CancelFunc) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	<-sigCh
	klog.Info("Received termination signal. Exiting gracefully...")
	cancel()
}
