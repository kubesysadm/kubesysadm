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

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CmMonitorSpec defines the desired state of CmMonitor
type CmMonitorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of CmMonitor. Edit cmmonitor_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// CmMonitorStatus defines the observed state of CmMonitor
type CmMonitorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// CmMonitor is the Schema for the cmmonitors API
type CmMonitor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CmMonitorSpec   `json:"spec,omitempty"`
	Status CmMonitorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CmMonitorList contains a list of CmMonitor
type CmMonitorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CmMonitor `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CmMonitor{}, &CmMonitorList{})
}
