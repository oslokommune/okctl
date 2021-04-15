package v1alpha1

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"regexp"
)

const (
	// ApplicationKind is a string value that represents the resource type
	ApplicationKind = "Application"
	// ApplicationAPIVersion defines the versioned schema of this representation
	ApplicationAPIVersion = "okctl.io/v1alpha1"

	minimumPossiblePort     = 1
	maximumPossiblePort     = 65535
	minimumPossibleReplicas = 0
)

// Application represents an application that can be deployed with okctl
type Application struct {
	metav1.TypeMeta `json:",inline"`

	Metadata ApplicationMeta `json:"metadata"`

	Images          []ApplicationImage `json:"images"`
	ImagePullSecret string             `json:"ImagePullSecret"`

	SubDomain string `json:"subDomain"`
	Port      int32  `json:"port"`

	Replicas int32 `json:"replicas"`

	Environment map[string]string   `json:"environment"`
	Volumes     []map[string]string `json:"volumes"`
}

func (a Application) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Metadata, validation.Required),
		validation.Field(&a.Images),
		validation.Field(&a.ImagePullSecret, is.Subdomain),
		validation.Field(&a.Port, validation.Min(minimumPossiblePort), validation.Max(maximumPossiblePort)),
		validation.Field(&a.Replicas, validation.Min(minimumPossibleReplicas)),
	)
}

// ApplicationMeta describes a unique application
type ApplicationMeta struct {
	// Name is a descriptive value given to the cluster, e.g., the name
	// of the team, product, project, etc.
	Name string `json:"name"`

	// Namespace defines which Kubernetes namespace to place the application in
	Namespace string `json:"namespace"`
}

// Validate ensures ApplicationMeta contains the right information
func (a ApplicationMeta) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Name,
			validation.Required,
			validation.Match(regexp.MustCompile("^[a-zA-Z-]{3,64}$")).Error("must consist of 3-64 characters (a-z, A-Z, -)")),
		validation.Field(&a.Namespace, validation.Required, is.Subdomain),
	)
}

// ApplicationImage describes a single image associated with the application
type ApplicationImage struct {
	// Name is the name of the image
	Name string
	// URI is the location of the image
	URI string
}

// Validate ensures ApplicationImage contains the right information
func (a ApplicationImage) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Name, validation.Required, is.Subdomain),
		validation.Field(&a.URI, is.URL),
	)
}

// ApplicationTypeMeta returns an initialised TypeMeta object
// for an Application
func ApplicationTypeMeta() metav1.TypeMeta {
	return metav1.TypeMeta{
		Kind:       ApplicationKind,
		APIVersion: ApplicationAPIVersion,
	}
}

// NewApplication returns an initialized application definition
func NewApplication() Application {
	return Application{
		TypeMeta: ApplicationTypeMeta(),
	}
}
