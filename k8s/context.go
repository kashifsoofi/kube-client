package k8s

import (
	"fmt"

	"k8s.io/client-go/tools/clientcmd"
)

func GetContexts(kubeConfigPath string) ([]string, string, error) {
	loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeConfigPath}
	configOverrides := &clientcmd.ConfigOverrides{}

	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	config, err := kubeConfig.RawConfig()
	if err != nil {
		return nil, "", fmt.Errorf("error getting RawConfig: %w", err)
	}

	contexts := []string{}
	for k := range config.Contexts {
		contexts = append(contexts, k)
	}

	return contexts, config.CurrentContext, nil
}
