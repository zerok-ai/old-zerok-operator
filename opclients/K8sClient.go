package opclients

import (
	"context"
	"encoding/json"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func GetK8sClient() *kubernetes.Clientset {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return clientset
}

func PrintPodsInCluster() {

	k8sClient := GetK8sClient().CoreV1()

	listOptions := metav1.ListOptions{}

	name := "default"

	services, _ := k8sClient.Services(name).List(context.Background(), listOptions)

	for _, service := range services.Items {
		if name == "default" && service.GetName() == "kubernetes" {
			continue
		}
		fmt.Println("namespace", name, "serviceName:", service.GetName(), "serviceKind:", service.Kind, "serviceLabels:", service.GetLabels(), service.Spec.Ports, "serviceSelector:", service.Spec.Selector)

		// labels.Parser
		set := labels.Set(service.Spec.Selector)

		if pods, err := k8sClient.Pods(name).List(context.Background(), metav1.ListOptions{LabelSelector: set.AsSelector().String()}); err != nil {
			fmt.Printf("List Pods of service[%s] error:%v\n", service.GetName(), err)
		} else {
			for _, pod := range pods.Items {
				fmt.Println("Pod", pod.GetName(), pod.Spec.NodeName, pod.Spec.Containers)
				payload := []patchStringValue{{
					Op:    "replace",
					Path:  "/metadata/labels/testLabel",
					Value: "897889",
				}}
				payloadBytes, _ := json.Marshal(payload)

				_, updateErr := k8sClient.Pods(pod.GetNamespace()).Patch(context.Background(), pod.GetName(), types.JSONPatchType, payloadBytes, metav1.PatchOptions{})
				if updateErr == nil {
					fmt.Printf("Pod %s labelled successfully.\n", pod.GetName())
				} else {
					fmt.Println(updateErr)
				}
			}
		}
	}

}
