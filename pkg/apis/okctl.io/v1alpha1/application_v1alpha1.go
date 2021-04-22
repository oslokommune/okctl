package v1alpha1

import (
	"fmt"
	"net/url"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	Image ApplicationImage `json:"image"`

	ImagePullSecret string `json:"ImagePullSecret"`
	SubDomain       string `json:"subDomain"`

	Port     int32 `json:"port"`
	Replicas int32 `json:"replicas"`

	Environment map[string]string `json:"environment"`

	Volumes []map[string]string `json:"volumes"`

	cluster Cluster
}

// Validate ensures Application contains the right information
func (a Application) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Metadata, validation.Required),
		validation.Field(&a.Image, validation.Required),
		validation.Field(&a.ImagePullSecret, is.Subdomain),
		validation.Field(&a.Port, validation.Min(minimumPossiblePort), validation.Max(maximumPossiblePort)),
		validation.Field(&a.Replicas, validation.Min(minimumPossibleReplicas)),
	)
}

// ApplicationMeta describes a unique application
type ApplicationMeta struct {
	// Name is a descriptive value given to the application
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

// ApplicationTypeMeta returns an initialised TypeMeta object
// for an Application
func ApplicationTypeMeta() metav1.TypeMeta {
	return metav1.TypeMeta{
		Kind:       ApplicationKind,
		APIVersion: ApplicationAPIVersion,
	}
}

// ApplicationImage defines which Docker image the application should use
type ApplicationImage struct {
	Name string `json:"name"`
	URI  string `json:"uri"`
}

// Validate ensures ApplicationImage contains the right information
func (a ApplicationImage) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Name, is.Subdomain),
		validation.Field(&a.URI), // TODO: exclusive OR
	)
}

// HasIngress returns true if the application has an ingress
func (a Application) HasIngress() bool {
	return a.SubDomain != ""
}

// HasService returns true if the application has a service
func (a Application) HasService() bool {
	return a.Port > 0
}

// Url returns the URL where the application is made available
func (a Application) Url() (url.URL, error) {
	tmpUrl, err := url.Parse(fmt.Sprintf("%s.%s", a.SubDomain, a.cluster.ClusterRootDomain))
	if err != nil {
		return url.URL{}, fmt.Errorf("parsing application URL: %w", err)
	}

	return *tmpUrl, nil
}

// NewApplication returns an initialized application definition
func NewApplication(cluster Cluster) Application {
	return Application{
		TypeMeta: ApplicationTypeMeta(),
		Image:    ApplicationImage{},
		cluster:  cluster,
	}
}
