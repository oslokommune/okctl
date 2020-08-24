// Package awsalbingresscontroller provides a Helm chart for installing:
// - https://github.com/kubernetes-sigs/aws-alb-ingress-controller
// - https://github.com/helm/charts/tree/master/incubator/aws-alb-ingress-controller
package awsalbingresscontroller

import "github.com/oslokommune/okctl/pkg/helm"

// New returns an initialised Helm chart for installing aws-alb-ingress-controller
func New(values *Values) *helm.Chart {
	return &helm.Chart{
		RepositoryName: "incubator",
		RepositoryURL:  "https://kubernetes-charts-incubator.storage.googleapis.com",
		ReleaseName:    "aws-alb-ingress-controller",
		Version:        "1.0.2",
		Chart:          "aws-alb-ingress-controller",
		Namespace:      "kube-system",
		Values:         values,
	}
}

// NewDefaultValues returns the mapped values.yml containing the default
// values
func NewDefaultValues(clusterName, serviceAccountName string) *Values {
	return &Values{
		ClusterName:           clusterName,
		AutoDiscoverAwsRegion: true,
		AutoDiscoverAwsVpcID:  true,
		Scope: Scope{
			IngressClass:    "alb",
			SingleNamespace: false,
		},
		EnableReadinessProbe:       true,
		ReadinessProbeInterval:     60,
		ReadinessProbeTimeout:      3,
		ReadinessProbeInitialDelay: 30,
		EnableLivenessProbe:        true,
		LivenessProbeInitialDelay:  30,
		LivenessProbeTimeout:       1,
		Rbac: Rbac{
			Create: true,
			ServiceAccount: ServiceAccount{
				Create: false,
				Name:   serviceAccountName,
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

type ServiceAccount struct {
	Create      bool        `yaml:"create"`
	Name        string      `yaml:"name"`
	Annotations interface{} `yaml:"annotations,omitempty"`
}

type Rbac struct {
	Create         bool           `yaml:"create"`
	ServiceAccount ServiceAccount `yaml:"serviceAccount"`
}

type Image struct {
	Repository string `yaml:"repository"`
	Tag        string `yaml:"tag"`
	PullPolicy string `yaml:"pullPolicy"`
}

type Resources struct {
	Limits   ResourceEntry `yaml:"limits"`
	Requests ResourceEntry `yaml:"requests"`
}

type ResourceEntry struct {
	CPU    string `yaml:"cpu"`
	Memory string `yaml:"memory"`
}
