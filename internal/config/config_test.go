package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad_Defaults(t *testing.T) {
	os.Clearenv()

	cfg := LoadFromEnv()

	if cfg.Namespace != "default" {
		t.Errorf("Expected default namespace, got %s", cfg.Namespace)
	}

	if cfg.LeaseName != "spring-app-leader-election" {
		t.Errorf("Expected default lease name, got %s", cfg.LeaseName)
	}

	if cfg.LeaseDuration != 15*time.Second {
		t.Errorf("Expected default lease duration of 15s, got %v", cfg.LeaseDuration)
	}
}

func TestLoad_FromEnv(t *testing.T) {

	if err := os.Setenv("NAMESPACE", "test-ns"); err != nil {
		t.Fatalf("failed to set env: %v", err)
	}
	if err := os.Setenv("LEASE_NAME", "custom-lease"); err != nil {
		t.Fatalf("failed to set env: %v", err)
	}
	if err := os.Setenv("POD_NAME", "my-pod"); err != nil {
		t.Fatalf("failed to set env: %v", err)
	}

}
