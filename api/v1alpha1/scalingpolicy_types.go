/*
Copyright 2025.

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

// ScalingPolicySpec defines the desired state of ScalingPolicy
type ScalingPolicySpec struct {
	MinReplicas     int    `json:"minReplicas"`
	MaxReplicas     int    `json:"maxReplicas"`
	CPUThreshold    int    `json:"cpuThreshold"`
	MemoryThreshold int    `json:"memoryThreshold"`
	MetricSource    string `json:"metricSource"`
	Foo             string `json:"foo,omitempty"`
}

// ScalingPolicyStatus defines the observed state of ScalingPolicy
type ScalingPolicyStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ScalingPolicy is the Schema for the scalingpolicies API
type ScalingPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ScalingPolicySpec   `json:"spec,omitempty"`
	Status ScalingPolicyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ScalingPolicyList contains a list of ScalingPolicy
type ScalingPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ScalingPolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ScalingPolicy{}, &ScalingPolicyList{})
}
