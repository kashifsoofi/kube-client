package k8s

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) GetServices(ns string) ([]string, error) {
	clientset, err := getClientset(c.configPath)
	if err != nil {
		return nil, err
	}

	services, err := clientset.CoreV1().Services(ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error getting services: %w", err)
	}

	serviceNames := []string{}
	for _, svc := range services.Items {
		serviceNames = append(serviceNames, svc.Name)
	}

	return serviceNames, nil
}
