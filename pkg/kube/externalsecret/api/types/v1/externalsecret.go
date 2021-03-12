// nolint
//go:generate controller-gen object paths=$GOFILE
package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

//+kubebuilder:object:generate:=true
type ExternalSecretData struct {
	Key      string `json:"key"`
	Name     string `json:"name"`
	Property string `json:"property,omitempty"`
}

//+kubebuilder:object:generate:=true
type ExternalSecretSpec struct {
	BackendType string `json:"backendType"`
	// +optional
	Template *ExternalSecretTemplate `json:"template,omitempty"`
	Data     []ExternalSecretData    `json:"data"`
}

//+kubebuilder:object:generate:=true
type ExternalSecretTemplate struct {
	Metadata ExternalSecretTemplateMetadata `json:"metadata"`
}

//+kubebuilder:object:generate:=true
type ExternalSecretTemplateMetadata struct {
	Annotations map[string]string `json:"annotations"`
	Labels      map[string]string `json:"labels"`
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
