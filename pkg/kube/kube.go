package kube

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

// KubeClient wraps Kubernetes client functionality
type KubeClient struct {
	clientset *kubernetes.Clientset
	config    *rest.Config
}

func NewKubeClient() (*KubeClient, error) {
	var config *rest.Config
	var err error

	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		// Use default kubeconfig path
		home := homedir.HomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user home directory: %w", err)
		}
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	// Use out-of-cluster configuration
	config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	return &KubeClient{clientset: clientset, config: config}, nil
}

// ListObjects lists all objects of a given type
func (k *KubeClient) ListObjects(gvk metav1.APIResource, namespace string) (*unstructured.UnstructuredList, error) {
	dynamicClient, err := dynamic.NewForConfig(k.config)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %w", err)
	}

	resourceClient := dynamicClient.Resource(schema.GroupVersionResource{
		Group:    gvk.Group,
		Version:  gvk.Version,
		Resource: gvk.Name,
	}).Namespace(namespace)

	return resourceClient.List(context.Background(), metav1.ListOptions{})
}

// FetchObjectSchema retrieves the schema for a given Kubernetes object
func (k *KubeClient) FetchObjectSchema(objectType string) (metav1.APIResource, error) {
	// Fetch all API resources from the server
	serverResources, err := k.clientset.Discovery().ServerPreferredResources()
	if err != nil {
		return metav1.APIResource{}, fmt.Errorf("error fetching server resources: %v", err)
	}

	// Iterate over aqll group versions
	for _, groupVersion := range serverResources {
		// Iterate over all API resources in this group version
		for _, res := range groupVersion.APIResources {
			// If the kind of this resource matches the objectType, return it
			if MatchesObjectType(res, objectType) {
				// parse group and version from the groupVersion string
				gv, err := schema.ParseGroupVersion(groupVersion.GroupVersion)
				if err != nil {
					return metav1.APIResource{}, fmt.Errorf("error parsing group version: %v", err)
				}

				// return the APIResource with the group and version
				res.Group = gv.Group
				res.Version = gv.Version
				return res, nil
			}
		}
	}

	return metav1.APIResource{}, fmt.Errorf("object type %s not found in group %s", objectType)
}

// MatchesObjectType checks if the given APIResource matches the objectType
func MatchesObjectType(res metav1.APIResource, objectType string) bool {
	if res.Kind == objectType || res.Name == objectType || res.SingularName == objectType {
		return true
	}

	for _, shortName := range res.ShortNames {
		if shortName == objectType {
			return true
		}
	}

	return false
}
