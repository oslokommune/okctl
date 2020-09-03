//go:generate controller-gen object paths=$GOFILE
package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

//+kubebuilder:object:generate:=true
type ExternalSecretData struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

//+kubebuilder:object:generate:=true
type ExternalSecretSpec struct {
	BackendType string               `json:"backendType"`
	Data        []ExternalSecretData `json:"data"`
}

//+kubebuilder:object:root:=true
//+kubebuilder:object:generate:=true
type ExternalSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ExternalSecretSpec `json:"spec"`
}

//+kubebuilder:object:root:=true
//+kubebuilder:object:generate:=true
type ExternalSecretList struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	
	Items []ExternalSecret `json:"items"`
}