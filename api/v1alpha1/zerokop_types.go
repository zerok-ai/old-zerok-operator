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

// ZerokopSpec defines the desired state of Zerokop
type ZerokopSpec struct {
	Http_response_headers_match HttpResponseHeadersMatch `json:"http_response_headers_match"`
}

// ZerokopStatus defines the observed state of Zerokop
type ZerokopStatus struct {
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
