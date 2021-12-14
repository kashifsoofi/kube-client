package k8s

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) GetDeployments(ns string) ([]string, error) {
	clientset, err := getClientset(c.configPath)
	if err != nil {
		return nil, err
	}

	deployments, err := clientset.AppsV1().Deployments(ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error getting deployments: %w", err)
	}

	deploymentNames := []string{}
	for _, d := range deployments.Items {
		deploymentNames = append(deploymentNames, d.Name)
	}

	return deploymentNames, nil
}
