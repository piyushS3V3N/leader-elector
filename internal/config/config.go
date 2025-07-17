package config

import (
	"os"
	"time"
)

type Config struct {
	Namespace      string
	LeaseName      string
	Identity       string
	KubeconfigPath string

	LeaseDuration time.Duration
	RenewDeadline time.Duration
	RetryPeriod   time.Duration
}

func LoadFromEnv() *Config {
	hostname, _ := os.Hostname()

	return &Config{
		Namespace:      getEnv("NAMESPACE", "default"),
		LeaseName:      getEnv("LEASE_NAME", "spring-app-leader-election"),
		Identity:       getEnv("POD_NAME", hostname),
		KubeconfigPath: getEnv("KUBECONFIG", ""),

		LeaseDuration: 15 * time.Second,
		RenewDeadline: 10 * time.Second,
		RetryPeriod:   2 * time.Second,
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
