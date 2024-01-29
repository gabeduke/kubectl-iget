package kube

import (
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// KubeClient wraps Kubernetes client functionality
type KubeClient struct {
	clientset *kubernetes.Clientset
}

// NewKubeClient creates a new KubeClient with appropriate configuration
func NewKubeClient() (*KubeClient, error) {
	var config *rest.Config
	var err error

	if kubeconfig := os.Getenv("KUBECONFIG"); kubeconfig != "" {
		// Use out-of-cluster configuration
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		// Use in-cluster configuration
		config, err = rest.InClusterConfig()
	}

	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &KubeClient{clientset: clientset}, nil
}

// FetchObjectSchema retrieves the schema for a given Kubernetes object
func (k *KubeClient) FetchObjectSchema(objectType string) (metav1.APIResource, error) {
	groupVersion := "v1" // For core objects, or determine dynamically for CRDs

	resList, err := k.clientset.Discovery().ServerResourcesForGroupVersion(groupVersion)
	if err != nil {
		return metav1.APIResource{}, fmt.Errorf("error fetching resources for group/version %s: %v", groupVersion, err)
	}

	for _, res := range resList.APIResources {
		if res.Kind == objectType {
			return res, nil
		}
	}

	return metav1.APIResource{}, fmt.Errorf("object type %s not found in group/version %s", objectType, groupVersion)
}
