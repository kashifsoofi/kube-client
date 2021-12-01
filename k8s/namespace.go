package k8s

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) GetNamespaces() ([]string, error) {
	clientset, err := getClientset(c.configPath)
	if err != nil {
		return nil, err
	}

	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error getting namespaces: %w", err)
	}

	nsNames := []string{}
	for _, ns := range namespaces.Items {
		nsNames = append(nsNames, ns.Name)
	}

	return nsNames, nil
}
