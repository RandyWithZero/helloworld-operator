/*
Copyright 2021.

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

package v1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	Pending = "mypending"
	Running = "myrunning"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HelloworldSpec defines the desired state of Helloworld
type HelloworldSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Helloworld. Edit helloworld_types.go to remove/update
	//Foo string `json:"foo,omitempty"`
	Template v1.PodTemplateSpec `json:"template" protobuf:"bytes,6,opt,name=template"`
}

// HelloworldStatus defines the observed state of Helloworld
type HelloworldStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Phase string `json:"phase,omitempty"`
}

func (status *HelloworldStatus) SetDefault(helloworld *Helloworld) bool {
	changed := false
	if helloworld.Status.Phase == "" {
		helloworld.Status.Phase = Pending
		changed = true
	}
	return changed
}
func (helloworld *Helloworld) StatusSetDefault() bool {
	return helloworld.Status.SetDefault(helloworld)
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Helloworld is the Schema for the helloworlds API
type Helloworld struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HelloworldSpec   `json:"spec,omitempty"`
	Status HelloworldStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// HelloworldList contains a list of Helloworld
type HelloworldList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Helloworld `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Helloworld{}, &HelloworldList{})
}
