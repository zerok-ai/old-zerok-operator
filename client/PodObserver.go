package client

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

type Target int

const (
	Deployment Target = 0
	Service           = 1
)

type PodObserver struct {
	informers informers.SharedInformerFactory
	target    Target
	Name      string
	Namespace string
	client    *k8sClient
	ch        chan struct{}
}

func (po *PodObserver) StartObservingPods() {

	podInformer := po.informers.Core().V1().Pods()

	podInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    po.handleAdd,
			DeleteFunc: po.handleDel,
		},
	)

	po.informers.Start(po.ch)
}

func (po *PodObserver) StopObservingPods() {
	close(po.ch)
}

func (po *PodObserver) handleAdd(obj interface{}) {
	pod := obj.(*v1.Pod)
	fmt.Printf("Pod added %v\n", pod.Name)
	po.handleClusterChange()
}

func (po *PodObserver) handleDel(obj interface{}) {
	pod := obj.(*v1.Pod)
	fmt.Printf("Pod deleted %v\n", pod.Name)
	po.handleClusterChange()
}

func (po *PodObserver) handleClusterChange() {
	switch po.target {
	case Deployment:
		po.client.LabelSpillAndSoakPodsForDeployment(po.Name, po.Namespace)
	case Service:
		po.client.LabelSpillAndSoakPodsForService(po.Name, po.Namespace)
	}
}
