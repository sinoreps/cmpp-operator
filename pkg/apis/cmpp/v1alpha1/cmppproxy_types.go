package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CMPPProxySpec defines the desired state of CMPPProxy
// +k8s:openapi-gen=true
type CMPPProxySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html

	// +kubebuilder:validation:MinLength=1
	Image string `json:"image"`

	// +kubebuilder:validation:MinLength=1
	ServerAddr string `json:"serverAddr"`

	Account   string `json:"account"`
	Password  string `json:"password"`
	Version   string `json:"version"`
	SourceID  string `json:"sourceID"`
	ServiceID string `json:"serviceID"`

	NumConnections    int32  `json:"numConnections"`
	ReportCallbackURL string `json:"reportCallbackURL"`
	ReplyCallbackURL  string `json:"replyCallbackURL"`
}

// CMPPProxyStatus defines the observed state of CMPPProxy
// +k8s:openapi-gen=true
type CMPPProxyStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	Pods []string `json:"pods"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CMPPProxy is the Schema for the httpproxies API
// +k8s:openapi-gen=true
type CMPPProxy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CMPPProxySpec   `json:"spec,omitempty"`
	Status CMPPProxyStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CMPPProxyList contains a list of CMPPProxy
type CMPPProxyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CMPPProxy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CMPPProxy{}, &CMPPProxyList{})
}
