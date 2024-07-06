/*
Copyright 2024 Wayne Wang<net_use@bzhy.com>.

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

package v1beta1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// PodCleanRuleSpec defines the spec for cleaning pod rules
type PodCleanRuleSpec struct {
	// NameSpace the namespace where the pod in is to be monitoring
	NameSpace string `json:"namespace,omitempty" protobuf:"bytes,3,opt,name=namespace"`

	// PrefixName the prefix of the pods' name which will be deleted.
	// normally, prefixname is the workload name. All pods in the NameSpace field specified
	// will be deleted if this field is empty.
	PrefixName string `json:"prefixname,omitempty" protobuf:"bytes,opt,name=prefixname"`

	// Age how long the pod is no-running which will to be deleted
	Age int32 `json:"age,omitempty" protobuf:"opt,name=age"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// PodCleanRule is the Schema for the PodCleanRule API
type PodCleanRule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec PodCleanRuleSpec `json:"spec,omitempty"`
}

//+kubebuilder:object:root=true

// PodCleanRuleList contains a list of PodCleanRule
type PodCleanRuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PodCleanRule `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PodCleanRule{}, &PodCleanRuleList{})
}
