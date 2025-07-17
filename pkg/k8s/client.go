package k8s

import (
	"github.com/yourorg/leader-elector/internal/config"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func NewClient(cfg *config.Config) (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error

	if cfg.KubeconfigPath != "" {
		config, err = clientcmd.BuildConfigFromFlags("", cfg.KubeconfigPath)
	} else {
		config, err = rest.InClusterConfig()
	}

	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}
