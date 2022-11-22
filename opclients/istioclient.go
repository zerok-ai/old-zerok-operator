package opclients

import (
	"context"
	"fmt"
	"log"

	"k8s.io/client-go/rest"

	operatorv1alpha1 "github.com/zerokdotai/zerok-operator/api/v1alpha1"
	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type patchStringValue struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"`
}

func ApplyEnvoyConfig(zerokopSpec operatorv1alpha1.ZerokopSpec) {
	ic := GetIstioClient()
	envoyFilterCrd := GetErrorRuleCrd(zerokopSpec)
	_, err := ic.NetworkingV1alpha3().EnvoyFilters(envoyFilterCrd.Namespace).Create(context.Background(), envoyFilterCrd, metav1.CreateOptions{})
	if err == nil {
		fmt.Println("Envoy Filter applied successfully.")
	} else {
		fmt.Println(err)
	}
	fmt.Println("Applied envoy filter crd.")
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
