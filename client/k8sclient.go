package client

import (
	"context"
	"encoding/json"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func LabelSpillAndSoakPods(podList *v1.PodList) {
	if podList == nil {
		fmt.Printf("Given podList is nil\n")
	} else if len(podList.Items) < 2 {
		fmt.Printf("Not enough pods to apply the configuration.\n")
	} else {
		spillPod := podList.Items[0]
		LabelPod(&spillPod, "/metadata/labels/zk-status", "enabled")
		LabelPod(&spillPod, "/metadata/labels/zk-route-mark", "spill")

		for i := 1; i < len(podList.Items); i++ {
			soakPod := podList.Items[i]
			LabelPod(&soakPod, "/metadata/labels/zk-status", "enabled")
			LabelPod(&soakPod, "/metadata/labels/zk-route-mark", "soak")
		}
	}
}

func LabelPod(pod *v1.Pod, path string, value string) {
	k8sClient := GetK8sClient().CoreV1()
	payload := []patchStringValue{{
		Op:    "replace",
		Path:  path,
		Value: value,
	}}
	payloadBytes, _ := json.Marshal(payload)
	_, updateErr := k8sClient.Pods(pod.GetNamespace()).Patch(context.Background(), pod.GetName(), types.JSONPatchType, payloadBytes, metav1.PatchOptions{})
	if updateErr == nil {
		fmt.Println(fmt.Sprintf("Pod %s labeled successfully for Path %s and Value %s.", pod.GetName(), path, value))
	} else {
		fmt.Println(updateErr)
	}
}

func LabelSpillAndSoakPodsForDeployment(Name string, Namespace string) {
	podList := GetPodsForDeployment(Name, Namespace)
	if podList == nil {
		fmt.Printf("Error while fetching podList for deployment %v.\n", Name)
	} else {
		LabelSpillAndSoakPods(podList)
	}
}

func LabelSpillAndSoakPodsForService(Name string, Namespace string) {
	podList := GetPodsForService(Name, Namespace)
	if podList == nil {
		fmt.Printf("Error while fetching podList for service %v.\n", Name)
	} else {
		LabelSpillAndSoakPods(podList)
	}
}

func GetPodsForDeployment(Name string, Namespace string) *v1.PodList {
	clientSet := GetK8sClient()

	k8sClient := clientSet.AppsV1()

	deployment, _ := k8sClient.Deployments(Namespace).Get(context.Background(), Name, metav1.GetOptions{})

	labelSelector := labels.Set(deployment.Spec.Selector.MatchLabels)

	options := metav1.ListOptions{
		LabelSelector: string(labelSelector.AsSelector().String()),
	}

	podList, err := clientSet.CoreV1().Pods(Namespace).List(context.Background(), options)

	if err != nil {
		fmt.Printf("Get Pods of deployment[%s] error:%v\n", deployment.GetName(), err)
		return nil
	} else {
		for _, pod := range podList.Items {
			fmt.Println(fmt.Sprintf("Pod found for deployment %s with name %s.", deployment.GetName(), pod.GetName()))
		}
	}

	return podList
}

func GetPodsForService(Name string, Namespace string) *v1.PodList {
	k8sClient := GetK8sClient().CoreV1()

	listOptions := metav1.GetOptions{}

	service, _ := k8sClient.Services(Namespace).Get(context.Background(), Name, listOptions)

	set := labels.Set(service.Spec.Selector)

	pods, err := k8sClient.Pods(Namespace).List(context.Background(), metav1.ListOptions{LabelSelector: set.AsSelector().String()})

	if err != nil {
		fmt.Printf("Get Pods of service[%s] error:%v\n", service.GetName(), err)
		return nil
	} else {
		for _, pod := range pods.Items {
			fmt.Println(fmt.Sprintf("Pod found for service %s with name %s.", service.GetName(), pod.GetName()))
		}
	}

	return pods
}

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
					fmt.Println(fmt.Sprintf("Pod %s labelled successfully.", pod.GetName()))
				} else {
					fmt.Println(updateErr)
				}
			}
		}
	}

}
