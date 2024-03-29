// Package awslbc provides a Helm chart for installing:
// - https://github.com/kubernetes-sigs/aws-load-balancer-controller
// - https://github.com/aws/eks-charts/tree/master/stable/aws-load-balancer-controller
package awslbc

import (
	"fmt"
	"time"

	"github.com/oslokommune/okctl/pkg/helm"
)

const (
	// ReleaseName is the name of the release
	ReleaseName = "aws-load-balancer-controller"
	// Namespace is the default namespace
	Namespace = "kube-system"
)

// New returns an initialised Helm chart for installing aws-alb-ingress-controller
func New(values *Values, timeout time.Duration) *helm.Chart {
	return &helm.Chart{
		RepositoryName: "eks",
		RepositoryURL:  "https://aws.github.io/eks-charts",
		ReleaseName:    ReleaseName,
		Version:        "1.1.3",
		Chart:          "aws-load-balancer-controller",
		Namespace:      Namespace,
		Timeout:        timeout,
		Values:         values,
	}
}

// NewDefaultValues returns the mapped values.yml containing the default
// values
func NewDefaultValues(clusterName, vpcID, region string) *Values {
	return &Values{
		ReplicaCount:                  1, // nolint: gomnd
		ClusterName:                   clusterName,
		Region:                        region,
		VpcID:                         vpcID,
		TerminationGracePeriodSeconds: 10, // nolint: gomnd
		EnableCertManager:             false,
		ServiceAccount: ServiceAccount{
			Create: false,
			Name:   "aws-load-balancer-controller", // this is too fragile, refactor
		},
		Rbac: Rbac{
			Create: true,
		},
		Image: Image{
			Repository: fmt.Sprintf("602401143452.dkr.ecr.%s.amazonaws.com/amazon/aws-load-balancer-controller", region),
			Tag:        "v2.1.1",
			PullPolicy: "IfNotPresent",
		},
		Resources: resources{
			Limits: resourceEntry{
				CPU:    "200m",
				Memory: "256Mi",
			},
			Requests: resourceEntry{
				CPU:    "100m",
				Memory: "128Mi",
			},
		},
		PodSecurityContext: PodSecurityContext{
			FsGroup: 65534, // nolint: gomnd
		},
		SecurityContext: SecurityContext{
			ReadOnlyRootFilesystem:   true,
			RunAsNonRoot:             true,
			AllowPrivilegeEscalation: false,
		},
		// The ingress class this controller will satisfy. If not specified, controller will match all
		// ingresses without ingress class annotation and ingresses of type alb
		IngressClass: "",
		LivenessProbe: LivenessProbe{
			FailureThreshold: 2, // nolint: gomnd
			HTTPGet: HTTPGet{
				Path:   "/healthz",
				Port:   61779, // nolint: gomnd
				Scheme: "HTTP",
			},
			InitialDelaySeconds: 30, // nolint: gomnd
			TimeoutSeconds:      10, // nolint: gomnd
		},
	}
}

// Values mimicks the chart values.yml
// nolint: maligned
type Values struct {
	ReplicaCount                              int                 `yaml:"replicaCount"`
	Image                                     Image               `yaml:"image"`
	ImagePullSecrets                          []interface{}       `yaml:"imagePullSecrets"`
	NameOverride                              string              `yaml:"nameOverride"`
	FullnameOverride                          string              `yaml:"fullnameOverride"`
	ClusterName                               interface{}         `yaml:"clusterName"`
	ServiceAccount                            ServiceAccount      `yaml:"serviceAccount"`
	Rbac                                      Rbac                `yaml:"rbac"`
	PodSecurityContext                        PodSecurityContext  `yaml:"podSecurityContext"`
	SecurityContext                           SecurityContext     `yaml:"securityContext"`
	TerminationGracePeriodSeconds             int                 `yaml:"terminationGracePeriodSeconds"`
	Resources                                 resources           `yaml:"resources"`
	PriorityClassName                         string              `yaml:"priorityClassName"`
	NodeSelector                              NodeSelector        `yaml:"nodeSelector"`
	Tolerations                               []interface{}       `yaml:"tolerations"`
	Affinity                                  Affinity            `yaml:"affinity"`
	PodAnnotations                            PodAnnotations      `yaml:"podAnnotations"`
	PodLabels                                 PodLabels           `yaml:"podLabels"`
	EnableCertManager                         bool                `yaml:"enableCertManager"`
	IngressClass                              string              `yaml:"ingressClass"`
	Region                                    interface{}         `yaml:"region"`
	VpcID                                     interface{}         `yaml:"vpcId"`
	AwsMaxRetries                             interface{}         `yaml:"awsMaxRetries,omitempty"`
	EnablePodReadinessGateInject              interface{}         `yaml:"enablePodReadinessGateInject,omitempty"`
	EnableShield                              interface{}         `yaml:"enableShield,omitempty"`
	EnableWaf                                 interface{}         `yaml:"enableWaf,omitempty"`
	EnableWafv2                               interface{}         `yaml:"enableWafv2,omitempty"`
	IngressMaxConcurrentReconciles            interface{}         `yaml:"ingressMaxConcurrentReconciles,omitempty"`
	LogLevel                                  interface{}         `yaml:"logLevel,omitempty"`
	MetricsBindAddr                           string              `yaml:"metricsBindAddr,omitempty"`
	WebhookBindPort                           interface{}         `yaml:"webhookBindPort,omitempty"`
	ServiceMaxConcurrentReconciles            interface{}         `yaml:"serviceMaxConcurrentReconciles,omitempty"`
	TargetgroupbindingMaxConcurrentReconciles interface{}         `yaml:"targetgroupbindingMaxConcurrentReconciles,omitempty"`
	SyncPeriod                                interface{}         `yaml:"syncPeriod,omitempty"`
	WatchNamespace                            interface{}         `yaml:"watchNamespace,omitempty"`
	LivenessProbe                             LivenessProbe       `yaml:"livenessProbe"`
	Env                                       interface{}         `yaml:"env,omitempty"`
	HostNetwork                               bool                `yaml:"hostNetwork,omitempty"`
	ExtraVolumeMounts                         interface{}         `yaml:"extraVolumeMounts,omitempty"`
	ExtraVolumes                              interface{}         `yaml:"extraVolumes,omitempty"`
	DefaultTags                               DefaultTags         `yaml:"defaultTags"`
	PodDisruptionBudget                       PodDisruptionBudget `yaml:"podDisruptionBudget"`
}

// Image ...
type Image struct {
	Repository string `yaml:"repository"`
	Tag        string `yaml:"tag"`
	PullPolicy string `yaml:"pullPolicy"`
}

// Annotations ...
type Annotations struct{}

// ServiceAccount ...
type ServiceAccount struct {
	Create      bool        `yaml:"create"`
	Annotations Annotations `yaml:"annotations"`
	Name        interface{} `yaml:"name"`
}

// Rbac ...
type Rbac struct {
	Create bool `yaml:"create"`
}

// PodSecurityContext ...
type PodSecurityContext struct {
	FsGroup int `yaml:"fsGroup"`
}

// SecurityContext ...
type SecurityContext struct {
	ReadOnlyRootFilesystem   bool `yaml:"readOnlyRootFilesystem"`
	RunAsNonRoot             bool `yaml:"runAsNonRoot"`
	AllowPrivilegeEscalation bool `yaml:"allowPrivilegeEscalation"`
}

type resources struct {
	Limits   resourceEntry `yaml:"limits"`
	Requests resourceEntry `yaml:"requests"`
}

type resourceEntry struct {
	CPU    string `yaml:"cpu"`
	Memory string `yaml:"memory"`
}

// NodeSelector ...
type NodeSelector struct{}

// Affinity ...
type Affinity struct{}

// PodAnnotations ...
type PodAnnotations struct{}

// PodLabels ...
type PodLabels struct{}

// HTTPGet ...
type HTTPGet struct {
	Path   string `yaml:"path"`
	Port   int    `yaml:"port"`
	Scheme string `yaml:"scheme"`
}

// LivenessProbe ...
type LivenessProbe struct {
	FailureThreshold    int     `yaml:"failureThreshold"`
	HTTPGet             HTTPGet `yaml:"httpGet"`
	InitialDelaySeconds int     `yaml:"initialDelaySeconds"`
	TimeoutSeconds      int     `yaml:"timeoutSeconds"`
}

// DefaultTags ...
type DefaultTags struct{}

// PodDisruptionBudget ...
type PodDisruptionBudget struct{}
