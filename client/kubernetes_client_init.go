package client

import (
	"fmt"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// InitK8SClient creates a Kubernetes clientset. If kubeconfigPath is provided it
// will be used to build a config, otherwise in-cluster config is attempted.
func InitK8SClient(kubeconfigPath string) (*kubernetes.Clientset, error) {
	var (
		config *rest.Config
		err    error
	)

	if kubeconfigPath != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load kubeconfig from %s: %v", kubeconfigPath, err)
		}
	} else if env := os.Getenv("KUBECONFIG"); env != "" {
		config, err = clientcmd.BuildConfigFromFlags("", env)
		if err != nil {
			return nil, fmt.Errorf("failed to load kubeconfig from %s: %v", env, err)
		}
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to load in-cluster config: %v", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %v", err)
	}
	return clientset, nil
}
