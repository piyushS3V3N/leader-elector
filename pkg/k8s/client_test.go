package k8s

import (
	"testing"

	"github.com/yourorg/leader-elector/internal/config"
)

func TestNewClient_InvalidKubeconfig(t *testing.T) {
	cfg := &config.Config{
		KubeconfigPath: "/nonexistent/file",
	}

	_, err := NewClient(cfg)
	if err == nil {
		t.Error("Expected error for invalid kubeconfig path, got nil")
	}
}
