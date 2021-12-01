package k8s

import (
	"fmt"

	"k8s.io/client-go/tools/clientcmd"
)

func (c *Client) GetContexts() []string {
	contexts := []string{}
	for k := range c.RawConfig.Contexts {
		contexts = append(contexts, k)
	}
	return contexts
}

func (c *Client) GetCurrentContext() string {
	return c.RawConfig.CurrentContext
}

func (c *Client) SwitchContext(ctx string) error {
	c.RawConfig.CurrentContext = ctx
	err := clientcmd.ModifyConfig(clientcmd.NewDefaultPathOptions(), c.RawConfig, true)
	if err != nil {
		return fmt.Errorf("error modify config: %w", err)
	}

	return nil
}
