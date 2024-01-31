package ui

import (
	"fmt"
	"github.com/gabeduke/kubectl-iget/pkg/kube"
)

type UIManager struct {
	kubeClient *kube.KubeClient
	config     UIConfig
}

type UIConfig struct {
	Object    string
	Namespace string
}

func NewUIManager(kubeClient *kube.KubeClient, config UIConfig) *UIManager {
	return &UIManager{
		kubeClient: kubeClient,
		config:     config,
	}
}

// Start initiates the UI workflow
func (u *UIManager) Start() {
	fmt.Println("Welcome to the Kubernetes Interactive Browser")

	// Fetch the object schema
	api, err := u.kubeClient.FetchObjectSchema(u.config.Object)
	if err != nil {
		fmt.Printf("Error fetching schema: %v\n", err)
		return
	}

	list, err := u.kubeClient.ListObjects(api, u.config.Namespace)
	for _, item := range list.Items {
		fmt.Println(item.GetName())
	}
}

// RenderMainMenu displays the main menu options to the user
func (u *UIManager) RenderMainMenu() error {
	// Implementation for rendering the main menu
	fmt.Println("1. Browse Kubernetes Objects")
	fmt.Println("2. Exit")
	// Add more options as needed
	return nil
}

// CaptureFieldSelections captures user selections for Kubernetes object fields
func (u *UIManager) CaptureFieldSelections() ([]string, error) {
	// Implementation for capturing field selections
	// This could be a series of prompts or a single input
	return []string{"field1", "field2"}, nil
}

// CaptureFilters captures user input for filters
func (u *UIManager) CaptureFilters() ([]string, error) {
	// Implementation for capturing filters
	return []string{"filter1=value1", "filter2=value2"}, nil
}
