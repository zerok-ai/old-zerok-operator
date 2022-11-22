/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	operatorv1alpha1 "github.com/zerokdotai/zerok-operator/api/v1alpha1"
	opclients "github.com/zerokdotai/zerok-operator/opclients"
	appsv1 "k8s.io/api/apps/v1"
)

// ZerokopReconciler reconciles a Zerokop object
type ZerokopReconciler struct {
	Client client.Client
	Scheme *runtime.Scheme
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Zerokop object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile

// +kubebuilder:rbac:groups=operator.zerok.ai,resources=zerokops,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=operator.zerok.ai,resources=zerokops/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=operator.zerok.ai,resources=zerokops/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=events,verbs=create;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch
func (r *ZerokopReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	zerokop := &operatorv1alpha1.Zerokop{}
	err := r.Client.Get(ctx, req.NamespacedName, zerokop)
	if err != nil {
		fmt.Printf("Error in getting the zerokop %v.\n", err)
	} else {
		fmt.Printf("zerokop spec %v.\n", zerokop.Spec)
	}
	opclients.ApplyEnvoyConfig(zerokop.Spec)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ZerokopReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&operatorv1alpha1.Zerokop{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
