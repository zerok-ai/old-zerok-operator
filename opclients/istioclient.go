package opclients

import (
	"fmt"
	"log"

	"k8s.io/client-go/rest"

	versionedclient "istio.io/client-go/pkg/clientset/versioned"
)

type patchStringValue struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"`
}

func ApplyEnvoyConfig() {
	GetIstioClient()
	fmt.Println("Create Istio client.")
}

func GetIstioClient() *versionedclient.Clientset {

	restConfig, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Failed to create k8s rest client: %s", err)
	}

	ic, err := versionedclient.NewForConfig(restConfig)
	if err != nil {
		log.Fatalf("Failed to create istio client: %s", err)
	}
	return ic
}
