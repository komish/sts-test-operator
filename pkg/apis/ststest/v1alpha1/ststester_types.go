package v1alpha1

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// StsTesterSpec defines the desired state of StsTester
// +k8s:openapi-gen=true
type StsTesterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// Reference: https://github.com/kubernetes/api/blob/master/apps/v1/types.go#L61
	PodManagementPolicy appsv1.PodManagementPolicyType `json:"podManagementPolicy,omitempty"`
}

// StsTesterStatus defines the observed state of StsTester
// +k8s:openapi-gen=true
type StsTesterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// Reference: https://github.com/kubernetes/api/blob/master/apps/v1/types.go#L61
	// We remove the omitempty here because we expect the operator to add a default PodManagementPolicy
	// if none was provided by the user.
	PodManagementPolicy appsv1.PodManagementPolicyType `json:"podManagementPolicy"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// StsTester is the Schema for the ststesters API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=ststesters,scope=Namespaced
type StsTester struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StsTesterSpec   `json:"spec,omitempty"`
	Status StsTesterStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// StsTesterList contains a list of StsTester
type StsTesterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StsTester `json:"items"`
}

func init() {
	SchemeBuilder.Register(&StsTester{}, &StsTesterList{})
}
