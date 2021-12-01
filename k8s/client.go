package k8s

import (
	"fmt"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/util/homedir"
)

type Client struct {
	configPath string
	RawConfig  api.Config
}

func NewClient() (*Client, error) {
	home := homedir.HomeDir()
	kubeConfigPath := filepath.Join(home, ".kube", "config")

	loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeConfigPath}
	configOverrides := &clientcmd.ConfigOverrides{}

	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	config, err := kubeConfig.RawConfig()
	if err != nil {
		return nil, fmt.Errorf("error getting RawConfig: %w", err)
	}

	return &Client{
		configPath: kubeConfigPath,
		RawConfig:  config,
	}, nil
}

func getClientset(kubeConfigPath string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, fmt.Errorf("error building config from flags: %w", err)
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error creating clientset: %w", err)
	}

	return clientset, nil
}
