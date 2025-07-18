package client

import (
        "fmt"

        "k8s.io/client-go/kubernetes"
        "k8s.io/client-go/rest"
        "k8s.io/client-go/tools/clientcmd"
)

// InitK8SClient initializes a Kubernetes client. When kubeconfigPath is not
// empty it will be used to create an out-of-cluster client. Otherwise it falls
// back to the in-cluster configuration.
func InitK8SClient(kubeconfigPath string) (*kubernetes.Clientset, error) {
        var (
                config *rest.Config
                err    error
        )

        if kubeconfigPath != "" {
                config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
                if err != nil {
                        return nil, fmt.Errorf("failed to load kubeconfig: %v", err)
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
