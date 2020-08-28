// Package awsalbingresscontroller provides a Helm chart for installing:
// - https://github.com/kubernetes-sigs/aws-alb-ingress-controller
// - https://github.com/helm/charts/tree/master/incubator/aws-alb-ingress-controller
package awsalbingresscontroller

import (
	"time"

	"github.com/oslokommune/okctl/pkg/helm"
)

// New returns an initialised Helm chart for installing aws-alb-ingress-controller
func New(values *Values) *helm.Chart {
	return &helm.Chart{
		RepositoryName: "incubator",
		RepositoryURL:  "https://kubernetes-charts-incubator.storage.googleapis.com",
		ReleaseName:    "aws-alb-ingress-controller",
		Version:        "1.0.2",
		Chart:          "aws-alb-ingress-controller",
		Namespace:      "kube-system",
		Timeout:        5 * time.Minute, // nolint: gomnd
		Values:         values,
	}
}

// NewDefaultValues returns the mapped values.yml containing the default
// values
func NewDefaultValues(clusterName, vpcID, region string) *Values {
	return &Values{
		ClusterName:           clusterName,
		AwsRegion:             region,
		AutoDiscoverAwsRegion: false,
		AwsVpcID:              vpcID,
		AutoDiscoverAwsVpcID:  false,
		Scope: Scope{
			IngressClass:    "alb",
			SingleNamespace: false,
		},
		EnableReadinessProbe:       true,
		ReadinessProbeInterval:     60, // nolint: gomnd
		ReadinessProbeTimeout:      3,  // nolint: gomnd
		ReadinessProbeInitialDelay: 30, // nolint: gomnd
		EnableLivenessProbe:        true,
		LivenessProbeInitialDelay:  30, // nolint: gomnd
		LivenessProbeTimeout:       1,
		Rbac: Rbac{
			Create: true,
			ServiceAccount: ServiceAccount{
				Create: false,
				Name:   "alb-ingress-controller", // this is too fragile, refactor
			},
		},
		Image: Image{
			Repository: "docker.io/amazon/aws-alb-ingress-controller",
			Tag:        "v1.1.8",
			PullPolicy: "IfNotPresent",
		},
		ReplicaCount: 1,
		Resources: Resources{
			Limits: ResourceEntry{
				CPU:    "100m",
				Memory: "128Mi",
			},
			Requests: ResourceEntry{
				CPU:    "100m",
				Memory: "128Mi",
			},
		},
	}
}

// Values maps up the aws-alb-ingress-controller helm chart values.yml
// - https://github.com/helm/charts/blob/master/incubator/aws-alb-ingress-controller/values.yaml
// nolint: maligned
type Values struct {
	ClusterName                string            `yaml:"clusterName"`
	AwsRegion                  string            `yaml:"awsRegion,omitempty"`
	AutoDiscoverAwsRegion      bool              `yaml:"autoDiscoverAwsRegion"`
	AwsVpcID                   string            `yaml:"awsVpcID,omitempty"`
	AutoDiscoverAwsVpcID       bool              `yaml:"autoDiscoverAwsVpcID"`
	Scope                      Scope             `yaml:"scope"`
	ExtraArgs                  map[string]string `yaml:"extraArgs,omitempty"`
	ExtraEnv                   map[string]string `yaml:"extraEnv,omitempty"`
	PodAnnotations             map[string]string `yaml:"podAnnotations,omitempty"`
	PodLabels                  map[string]string `yaml:"podLabels,omitempty"`
	EnableReadinessProbe       bool              `yaml:"enableReadinessProbe"`
	ReadinessProbeInterval     int               `yaml:"readinessProbeInterval"`
	ReadinessProbeTimeout      int               `yaml:"readinessProbeTimeout"`
	ReadinessProbeInitialDelay int               `yaml:"readinessProbeInitialDelay"`
	EnableLivenessProbe        bool              `yaml:"enableLivenessProbe"`
	LivenessProbeInitialDelay  int               `yaml:"livenessProbeInitialDelay"`
	LivenessProbeTimeout       int               `yaml:"livenessProbeTimeout"`
	Rbac                       Rbac              `yaml:"rbac"`
	Image                      Image             `yaml:"image"`
	ReplicaCount               int               `yaml:"replicaCount"`
	NameOverride               string            `yaml:"nameOverride,omitempty"`
	FullnameOverride           string            `yaml:"fullnameOverride,omitempty"`
	Resources                  Resources         `yaml:"resources"`
	NodeSelector               map[string]string `yaml:"nodeSelector,omitempty"`
	Tolerations                []interface{}     `yaml:"tolerations,omitempty"`
	Affinity                   interface{}       `yaml:"affinity,omitempty"`
	VolumeMounts               []interface{}     `yaml:"volumeMounts,omitempty"`
	Volumes                    []interface{}     `yaml:"volumes,omitempty"`
	PriorityClassName          string            `yaml:"priorityClassName,omitempty"`
	SecurityContext            interface{}       `yaml:"securityContext,omitempty"`
	ContainerSecurityContext   interface{}       `yaml:"containerSecurityContext,omitempty"`
}

// Scope determines what namespaces the controller should watch
// for ingress definitions
type Scope struct {
	IngressClass    string `yaml:"ingressClass"`
	SingleNamespace bool   `yaml:"singleNamespace"`
	WatchNamespace  string `yaml:"watchNamespace,omitempty"`
}

// ServiceAccount determines how the service account is created
type ServiceAccount struct {
	Create      bool        `yaml:"create"`
	Name        string      `yaml:"name"`
	Annotations interface{} `yaml:"annotations,omitempty"`
}

// Rbac determines if rbac should be setup
type Rbac struct {
	Create         bool           `yaml:"create"`
	ServiceAccount ServiceAccount `yaml:"serviceAccount"`
}

// Image defines the container image to use
type Image struct {
	Repository string `yaml:"repository"`
	Tag        string `yaml:"tag"`
	PullPolicy string `yaml:"pullPolicy"`
}

// Resources determines the limits and requests of the pod
type Resources struct {
	Limits   ResourceEntry `yaml:"limits"`
	Requests ResourceEntry `yaml:"requests"`
}

// ResourceEntry determines the cpu and memory usage
type ResourceEntry struct {
	CPU    string `yaml:"cpu"`
	Memory string `yaml:"memory"`
}
