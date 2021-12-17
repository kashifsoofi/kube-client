package k8s

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

func (c *Client) GetPods(ns string) ([]string, error) {
	clientset, err := getClientset(c.configPath)
	if err != nil {
		return nil, err
	}

	pods, err := clientset.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error getting pods: %w", err)
	}

	podNames := []string{}
	for _, pod := range pods.Items {
		podNames = append(podNames, pod.Name)
	}

	return podNames, nil
}

func (c *Client) GetPodLogStream(ns, podName string) (io.ReadCloser, error) {
	clientset, err := getClientset(c.configPath)
	if err != nil {
		return nil, err
	}

	count := int64(100)
	podLogOptions := v1.PodLogOptions{
		Follow:    true,
		TailLines: &count,
	}

	podLogRequest := clientset.CoreV1().Pods(ns).GetLogs(podName, &podLogOptions)
	stream, err := podLogRequest.Stream(context.TODO())
	if err != nil {
		return nil, err
	}

	return stream, nil
}

type PortForwardAPodRequest struct {
	Namespace string
	PodName   string
	// LocalPort is the local port that will be selected to expose the PodPort
	LocalPort string
	// PodPort is the target port for the pod
	PodPort string
	// Steams configures where to write or read input from
	Out    io.Writer
	ErrOut io.Writer
	// StopCh is the channel used to manage the port forward lifecycle
	StopCh <-chan struct{}
	// ReadyCh communicates when the tunnel is ready to receive traffic
	ReadyCh chan struct{}
}

func (c *Client) PortForwardAPod(req PortForwardAPodRequest) error {
	config, err := clientcmd.BuildConfigFromFlags("", c.configPath)
	if err != nil {
		return fmt.Errorf("error building config from flags: %w", err)
	}

	path := fmt.Sprintf(
		"/api/v1/namespaces/%s/pods/%s/portforward",
		req.Namespace,
		req.PodName)
	hostIP := strings.TrimLeft(config.Host, "htps:/")

	transport, upgrader, err := spdy.RoundTripperFor(config)
	if err != nil {
		return err
	}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, &url.URL{Scheme: "https", Path: path, Host: hostIP})
	fw, err := portforward.New(dialer, []string{fmt.Sprintf("%s:%s", req.LocalPort, req.PodPort)}, req.StopCh, req.ReadyCh, req.Out, req.ErrOut)
	if err != nil {
		return err
	}
	return fw.ForwardPorts()
}
