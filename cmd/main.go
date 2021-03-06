package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// switch context
	// if err := switchContext("my-context-to-switch", *kubeconfig); err != nil {
	// 	fmt.Printf("Error switching context: %s\n", err)
	// }

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	listNamespaces(clientset)

	listAllPods(clientset)
}

func listNamespaces(clientset *kubernetes.Clientset) {
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error getting namespaces: %w\n", err)
		panic(err.Error())
	}

	fmt.Printf("There are %d namespaces in the cluster\n", len(namespaces.Items))

	for _, ns := range namespaces.Items {
		fmt.Printf("%s - %s\n", ns.Name, ns.Kind)
	}
}

func listAllPods(clientset *kubernetes.Clientset) {
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	// Examples for error handling:
	// - Use helper functions like e.g. errors.IsNotFound()
	// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
	namespace := "default"
	pod := "example-xxxxx"
	_, err = clientset.CoreV1().Pods(namespace).Get(context.TODO(), pod, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("Pod %s in namespace %s not found\n", pod, namespace)
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting pod %s in namespace %s: %v\n",
			pod, namespace, statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found pod %s in namespace %s\n", pod, namespace)
	}

	time.Sleep(10 * time.Second)
}

func switchContext(ctx, kubeconfigPath string) (err error) {
	// loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath}
	configOverrides := &clientcmd.ConfigOverrides{}

	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	config, err := kubeConfig.RawConfig()
	if err != nil {
		return fmt.Errorf("error getting RawConfig: %w", err)
	}

	for k := range config.Contexts {
		fmt.Printf("%s", k)
		if k == config.CurrentContext {
			fmt.Printf(" *")
		}
		fmt.Println()
	}

	if config.Contexts[ctx] == nil {
		return fmt.Errorf("context %s doesn't exists", ctx)
	}

	config.CurrentContext = ctx
	err = clientcmd.ModifyConfig(clientcmd.NewDefaultPathOptions(), config, true)
	if err != nil {
		return fmt.Errorf("error ModifyConfig: %w", err)
	}

	fmt.Printf("Switched to context \"%s\"", ctx)
	return nil
}
