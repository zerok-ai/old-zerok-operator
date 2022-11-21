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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type Zerokop_ApplyTo int32

const (
	//Applies the error rule to response code
	RESPONSE_CODE Zerokop_ApplyTo = 0
	//Applies the error rule to response body
	RESPONSE_BODY Zerokop_ApplyTo = 1
)

var (
	Zerokop_ApplyTo_name = map[int32]string{
		0: "RESPONSE_CODE",
		1: "RESPONSE_BODY",
	}
	Zerokop_ApplyTo_value = map[string]int32{
		"RESPONSE_CODE": 0,
		"RESPONSE_BODY": 1,
	}
)

type ZerokopErrorRule struct {
	ApplyTo Zerokop_ApplyTo `json:"ApplyTo"`
	Code    int32           `json:"code"`
}

// ZerokopSpec defines the desired state of Zerokop
type ZerokopSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Size int32 `json:"size"`

	ErrorRules []ZerokopErrorRule `json:"ErrorRules"`
}

// ZerokopStatus defines the observed state of Zerokop
type ZerokopStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// +operator-sdk:csv:customresourcedefinitions:type=status
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

//+kubebuilder:object:root=true

// Zerokop is the Schema for the zerokops API
// +kubebuilder:subresource:status
// +genclient:nonNamespaced
type Zerokop struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ZerokopSpec   `json:"spec,omitempty"`
	Status ZerokopStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ZerokopList contains a list of Zerokop
type ZerokopList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Zerokop `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Zerokop{}, &ZerokopList{})
}
