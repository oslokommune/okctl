// Package argocd provides a Helm chart for installing:
// - https://argoproj.github.io/argo-cd/
// - https://github.com/argoproj/argo-helm
package argocd

import (
	"fmt"
	"time"

	"github.com/oslokommune/okctl/pkg/helm"
)

// New returns an initialised Helm chart
func New(values *Values) *helm.Chart {
	return &helm.Chart{
		RepositoryName: "argo",
		RepositoryURL:  "https://argoproj.github.io/argo-helm",
		ReleaseName:    "argocd",
		Version:        "2.6.2",
		Chart:          "argo-cd",
		Namespace:      "argocd",
		Timeout:        5 * time.Minute, // nolint: gomnd
		Values:         values,
	}
}

// ValuesOpts contains input for creating the default values
type ValuesOpts struct {
	URL                  string
	HostName             string
	CertificateARN       string
	ClientID             string
	Organisation         string
	Team                 string
	RepoURL              string
	RepoName             string
	PrivateKeySecretName string
	PrivateKeySecretKey  string
}

// NewDefaultValues returns the default values for the chart
// nolint: gomnd
func NewDefaultValues(opts ValuesOpts) *Values {
	return &Values{
		InstallCRDs: false,
		Global: global{
			Image: image{
				Repository: "argoproj/argocd",
				Tag:        "v1.7.2",
				PullPolicy: "IfNotPresent",
			},
			SecurityContext: securityContext{
				RunAsNonRoot: true,
				RunAsGroup:   999,
				RunAsUser:    999,
				FsGroup:      999,
			},
		},
		Controller: controller{
			Name: "application-controller",
			Args: controllerArgs{
				StatusProcessors:    "20",
				OperationProcessors: "10",
				AppResyncPeriod:     "180",
			},
			LogLevel:      "info",
			ContainerPort: 8082,
			ReadinessProbe: probe{
				FailureThreshold:    3,
				InitialDelaySeconds: 10,
				PeriodSeconds:       10,
				SuccessThreshold:    1,
				TimeoutSeconds:      1,
			},
			LivenessProbe: probe{
				FailureThreshold:    3,
				InitialDelaySeconds: 10,
				PeriodSeconds:       10,
				SuccessThreshold:    1,
				TimeoutSeconds:      1,
			},
			Service: service{
				Port:     8082,
				PortName: "https-controller",
			},
			Resources: resources{
				Limits: resourceEntry{
					CPU:    "500m",
					Memory: "512Mi",
				},
				Requests: resourceEntry{
					CPU:    "250m",
					Memory: "256Mi",
				},
			},
			ServiceAccount: serviceAccount{
				Create: true,
				Name:   "argocd-application-controller",
			},
			ClusterAdminAccess: clusterAdminAccess{
				Enabled: true,
			},
		},
		Dex: dex{
			Enabled: true,
			Name:    "dex-server",
			Image: image{
				Repository: "quay.io/dexidp/dex",
				Tag:        "v2.22.0",
				PullPolicy: "IfNotPresent",
			},
			ServiceAccount: serviceAccount{
				Create: true,
				Name:   "argocd-dex-server",
			},
			VolumeMounts: []volumeMounts{
				{
					Name:      "static-files",
					MountPath: "/shared",
				},
			},
			Volumes: []volumes{
				{
					Name:     "static-files",
					EmptyDir: map[string]string{},
				},
			},
			ContainerPortHTTP: 5556,
			ServicePortHTTP:   5556,
			ContainerPortGrpc: 5557,
			ServicePortGrpc:   5557,
			Resources: resources{
				Limits: resourceEntry{
					CPU:    "50m",
					Memory: "64Mi",
				},
				Requests: resourceEntry{
					CPU:    "10m",
					Memory: "32Mi",
				},
			},
		},
		Redis: redis{
			Enabled:       true,
			Name:          "redis",
			ContainerPort: 6379,
			ServicePort:   6379,
			Image: image{
				Repository: "redis",
				Tag:        "5.0.8",
				PullPolicy: "IfNotPresent",
			},
			SecurityContext: securityContext{
				RunAsNonRoot: true,
				RunAsGroup:   1000,
				RunAsUser:    1000,
				FsGroup:      1000,
			},
			Resources: resources{
				Limits: resourceEntry{
					CPU:    "200m",
					Memory: "128Mi",
				},
				Requests: resourceEntry{
					CPU:    "100m",
					Memory: "64Mi",
				},
			},
		},
		Server: server{
			Name:     "server",
			Replicas: 1,
			Autoscaling: autoscaling{
				Enabled:                           false,
				MinReplicas:                       1,
				MaxReplicas:                       5,
				TargetCPUUtilizationPercentage:    50,
				TargetMemoryUtilizationPercentage: 50,
			},
			ExtraArgs:     []string{"--insecure"},
			LogLevel:      "info",
			ContainerPort: 8080,
			ReadinessProbe: probe{
				FailureThreshold:    3,
				InitialDelaySeconds: 10,
				PeriodSeconds:       10,
				SuccessThreshold:    1,
				TimeoutSeconds:      1,
			},
			LivenessProbe: probe{
				FailureThreshold:    3,
				InitialDelaySeconds: 10,
				PeriodSeconds:       10,
				SuccessThreshold:    1,
				TimeoutSeconds:      1,
			},
			Resources: resources{
				Limits: resourceEntry{
					CPU:    "100m",
					Memory: "128Mi",
				},
				Requests: resourceEntry{
					CPU:    "50m",
					Memory: "64Mi",
				},
			},
			Certificate: certificate{
				Enabled: false,
			},
			Service: serverService{
				Type:                     "ClusterIP",
				ServicePortHTTP:          80,
				ServicePortHTTPS:         443,
				ServicePortHTTPName:      "http",
				ServicePortHTTPSName:     "https",
				LoadBalancerIP:           "",
				LoadBalancerSourceRanges: nil,
			},
			Metrics: metrics{
				Enabled: false,
			},
			ServiceAccount: serviceAccount{
				Create: true,
				Name:   "argocd-server",
			},
			Ingress: ingress{
				Enabled: true,
				Annotations: map[string]string{
					"kubernetes.io/ingress.class":                    "alb",
					"alb.ingress.kubernetes.io/scheme":               "internet-facing",
					"alb.ingress.kubernetes.io/target-type":          "instance",
					"alb.ingress.kubernetes.io/healthcheck-path":     "/healthz",
					"alb.ingress.kubernetes.io/listen-ports":         `[{"HTTP":80},{"HTTPS":443}]`,
					"alb.ingress.kubernetes.io/actions.ssl-redirect": `{"Type":"redirect","RedirectConfig":{"Protocol":"HTTPS","Port":"443","StatusCode":"HTTP_301"}}`,
					"alb.ingress.kubernetes.io/certificate-arn":      opts.CertificateARN,
				},
				Hosts: []string{opts.HostName},
				Paths: []string{"/"},
				ExtraPaths: []path{
					{
						Path: "/*",
						Backend: backend{
							ServiceName: "ssl-redirect",
							ServicePort: "use-annotation",
						},
					},
				},
				HTTPS: false,
			},
			IngressGrpc: ingressGrpc{
				Enabled: false,
			},
			Route: route{
				Enabled: false,
			},
			Config: serverConfig{
				URL:                   opts.URL,
				UsersAnonymousEnabled: "false",
				DexConfig:             fmt.Sprintf(dexConfig, opts.ClientID, opts.Organisation, opts.Team),
				Repositories:          fmt.Sprintf(repositoriesConfig, opts.RepoURL, opts.RepoName, opts.PrivateKeySecretName, opts.PrivateKeySecretKey),
				AdminEnabled:          "false",
			},
			RBACConfig: rbacConfig{
				PolicyCSV: fmt.Sprintf(policyCSV, opts.Organisation, opts.Team),
				Scopes:    `[groups, email]`,
			},
			ClusterAdminAccess: clusterAdminAccess{
				Enabled: true,
			},
		},
		RepoServer: repoServer{
			Name:     "repo-server",
			Replicas: 1,
			Autoscaling: autoscaling{
				Enabled:                           false,
				MinReplicas:                       1,
				MaxReplicas:                       5,
				TargetCPUUtilizationPercentage:    50,
				TargetMemoryUtilizationPercentage: 50,
			},
			LogLevel:      "info",
			ContainerPort: 8081,
			ReadinessProbe: probe{
				FailureThreshold:    3,
				InitialDelaySeconds: 10,
				PeriodSeconds:       10,
				SuccessThreshold:    1,
				TimeoutSeconds:      1,
			},
			LivenessProbe: probe{
				FailureThreshold:    3,
				InitialDelaySeconds: 10,
				PeriodSeconds:       10,
				SuccessThreshold:    1,
				TimeoutSeconds:      1,
			},
			Resources: resources{
				Limits: resourceEntry{
					CPU:    "50m",
					Memory: "128Mi",
				},
				Requests: resourceEntry{
					CPU:    "10m",
					Memory: "64Mi",
				},
			},
			Service: service{
				Port:     8081,
				PortName: "https-repo-server",
			},
			Metrics: metrics{
				Enabled: false,
			},
			ServiceAccount: serviceAccount{
				Create: false,
			},
		},
		Configs: configs{
			SSHKnownHosts: sshKnownHosts,
			CreateSecret:  false,
		},
	}
}

const sshKnownHosts = `bitbucket.org ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAubiN81eDcafrgMeLzaFPsw2kNvEcqTKl/VqLat/MaB33pZy0y3rJZtnqwR2qOOvbwKZYKiEO1O6VqNEBxKvJJelCq0dTXWT5pbO2gDXC6h6QDXCaHo6pOHGPUy+YBaGQRGuSusMEASYiWunYN0vCAI8QaXnWMXNMdFP3jHAJH0eDsoiGnLPBlBp4TNm6rYI74nMzgz3B9IikW4WVK+dc8KZJZWYjAuORU3jc1c/NPskD2ASinf8v3xnfXeukU0sJ5N6m5E8VLjObPEO+mN2t/FZTMZLiFqPWc/ALSqnMnnhwrNi2rbfg/rd/IpL8Le3pSBne8+seeFVBoGqzHM9yXw==
github.com ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAq2A7hRGmdnm9tUDbO9IDSwBK6TbQa+PXYPCPy6rbTrTtw7PHkccKrpp0yVhp5HdEIcKr6pLlVDBfOLX9QUsyCOV0wzfjIJNlGEYsdlLJizHhbn2mUjvSAHQqZETYP81eFzLQNnPHt4EVVUh7VfDESU84KezmD5QlWpXLmvU31/yMf+Se8xhHTvKSCZIFImWwoG6mbUoWf9nzpIoaSjB+weqqUUmpaaasXVal72J+UX2B+2RPW3RcT0eOzQgqlJL3RKrTJvdsjE3JEAvGq3lGHSZXy28G3skua2SmVi/w4yCE6gbODqnTWlg7+wC604ydGXA8VJiS5ap43JXiUFFAaQ==
gitlab.com ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBFSMqzJeV9rUzU4kWitGjeR4PWSa29SPqJ1fVkhtj3Hw9xjLVXVYrU9QlYWrOLXBpQ6KWjbjTDTdDkoohFzgbEY=
gitlab.com ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIAfuCHKVTjquxvt6CM6tdG4SLp1Btn/nOeHHE5UOzRdf
gitlab.com ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCsj2bNKTBSpIYDEGk9KxsGh3mySTRgMtXL583qmBpzeQ+jqCMRgBqB98u3z++J1sKlXHWfM9dyhSevkMwSbhoR8XIq/U0tCNyokEi/ueaBMCvbcTHhO7FcwzY92WK4Yt0aGROY5qX2UKSeOvuP4D6TPqKF1onrSzH9bx9XUf2lEdWT/ia1NEKjunUqu1xOB/StKDHMoX4/OKyIzuS0q/T1zOATthvasJFoPrAjkohTyaDUz2LN5JoH839hViyEG82yB+MjcFV5MU3N1l1QL3cVUCh93xSaua1N85qivl+siMkPGbO5xR/En4iEY6K2XPASUEMaieWVNTRCtJ4S8H+9
ssh.dev.azure.com ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC7Hr1oTWqNqOlzGJOfGJ4NakVyIzf1rXYd4d7wo6jBlkLvCA4odBlL0mDUyZ0/QUfTTqeu+tm22gOsv+VrVTMk6vwRU75gY/y9ut5Mb3bR5BV58dKXyq9A9UeB5Cakehn5Zgm6x1mKoVyf+FFn26iYqXJRgzIZZcZ5V6hrE0Qg39kZm4az48o0AUbf6Sp4SLdvnuMa2sVNwHBboS7EJkm57XQPVU3/QpyNLHbWDdzwtrlS+ez30S3AdYhLKEOxAG8weOnyrtLJAUen9mTkol8oII1edf7mWWbWVf0nBmly21+nZcmCTISQBtdcyPaEno7fFQMDD26/s0lfKob4Kw8H
vs-ssh.visualstudio.com ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC7Hr1oTWqNqOlzGJOfGJ4NakVyIzf1rXYd4d7wo6jBlkLvCA4odBlL0mDUyZ0/QUfTTqeu+tm22gOsv+VrVTMk6vwRU75gY/y9ut5Mb3bR5BV58dKXyq9A9UeB5Cakehn5Zgm6x1mKoVyf+FFn26iYqXJRgzIZZcZ5V6hrE0Qg39kZm4az48o0AUbf6Sp4SLdvnuMa2sVNwHBboS7EJkm57XQPVU3/QpyNLHbWDdzwtrlS+ez30S3AdYhLKEOxAG8weOnyrtLJAUen9mTkol8oII1edf7mWWbWVf0nBmly21+nZcmCTISQBtdcyPaEno7fFQMDD26/s0lfKob4Kw8H
`

const policyCSV = `g, %s:%s, role:admin
`

const repositoriesConfig = `- url: %s
  type: git
  name: %s
  sshPrivateKeySecret:
    name: %s
    key: %s
`

const dexConfig = `connectors:
- type: gihub
  id: github
  name: Github
  config:
    clientID: %s
    clientSecret: $dex.github.clientSecret
    orgs:
    - name: %s
      teams:
      - %s
`

// Values contains the parameters we map up
type Values struct {
	InstallCRDs bool       `yaml:"installCRDs"`
	Global      global     `yaml:"global"`
	Controller  controller `yaml:"controller"`
	Dex         dex        `yaml:"dex"`
	Redis       redis      `yaml:"redis"`
	Server      server     `yaml:"server"`
	RepoServer  repoServer `yaml:"repoServer"`
	Configs     configs    `yaml:"configs"`
}

type configs struct {
	SSHKnownHosts string `yaml:"knownHosts.data.ssh_known_hosts"`
	CreateSecret  bool   `yaml:"secret.createSecret"`
}

type repoServer struct {
	Name           string         `yaml:"name"`
	Replicas       int            `yaml:"replicas"`
	Autoscaling    autoscaling    `yaml:"autoscaling"`
	LogLevel       string         `yaml:"logLevel"`
	ContainerPort  int            `yaml:"containerPort"`
	ReadinessProbe probe          `yaml:"readinessProbe"`
	LivenessProbe  probe          `yaml:"livenessProbe"`
	Resources      resources      `yaml:"resources"`
	Service        service        `yaml:"service"`
	Metrics        metrics        `yaml:"metrics"`
	ServiceAccount serviceAccount `yaml:"serviceAccount"`
}

// nolint: maligned
type server struct {
	Name               string             `yaml:"name"`
	Replicas           int                `yaml:"replicas"`
	Autoscaling        autoscaling        `yaml:"autoscaling"`
	ExtraArgs          []string           `yaml:"extraArgs"`
	LogLevel           string             `yaml:"logLevel"`
	ContainerPort      int                `yaml:"containerPort"`
	ReadinessProbe     probe              `yaml:"readinessProbe"`
	LivenessProbe      probe              `yaml:"livenessProbe"`
	Resources          resources          `yaml:"resources"`
	Certificate        certificate        `yaml:"certificate"`
	Service            serverService      `yaml:"service"`
	Metrics            metrics            `yaml:"metrics"`
	ServiceAccount     serviceAccount     `yaml:"serviceAccount"`
	Ingress            ingress            `yaml:"ingress"`
	IngressGrpc        ingressGrpc        `yaml:"ingressGrpc"`
	Route              route              `yaml:"route"`
	Config             serverConfig       `yaml:"config"`
	RBACConfig         rbacConfig         `yaml:"rbacConfig"`
	ClusterAdminAccess clusterAdminAccess `yaml:"clusterAdminAccess"`
}

// https://github.com/argoproj/argo-cd/blob/master/docs/operator-manual/argocd-cm.yaml
type serverConfig struct {
	URL                   string `yaml:"url"`
	UsersAnonymousEnabled string `yaml:"users.anonymous.enabled"`
	DexConfig             string `yaml:"dex.config"`
	Repositories          string `yaml:"repositories"`
	AdminEnabled          string `yaml:"admin.enabled"`
}

type rbacConfig struct {
	PolicyCSV string `yaml:"policy.csv"`
	Scopes    string `yaml:"scopes"`
}

type serverService struct {
	Annotations              map[string]string `yaml:"annotations"`
	Labels                   map[string]string `yaml:"labels"`
	Type                     string            `yaml:"type"`
	ServicePortHTTP          int               `yaml:"servicePortHttp"`
	ServicePortHTTPS         int               `yaml:"servicePortHttps"`
	ServicePortHTTPName      string            `yaml:"servicePortHttpName"`
	ServicePortHTTPSName     string            `yaml:"servicePortHttpsName"`
	LoadBalancerIP           string            `yaml:"loadBalancerIP"`
	LoadBalancerSourceRanges []interface{}     `yaml:"loadBalancerSourceRanges"`
}

type redis struct {
	Enabled         bool            `yaml:"enabled"`
	Name            string          `yaml:"name"`
	ContainerPort   int             `yaml:"containerPort"`
	ServicePort     int             `yaml:"servicePort"`
	Image           image           `yaml:"image"`
	SecurityContext securityContext `yaml:"securityContext"`
	Resources       resources       `yaml:"resources"`
}

type dex struct {
	Enabled           bool           `yaml:"enabled"`
	Name              string         `yaml:"name"`
	Image             image          `yaml:"image"`
	ServiceAccount    serviceAccount `yaml:"serviceAccount"`
	VolumeMounts      []volumeMounts `yaml:"volumeMounts"`
	Volumes           []volumes      `yaml:"volumes"`
	ContainerPortHTTP int            `yaml:"containerPortHttp"`
	ServicePortHTTP   int            `yaml:"servicePortHttp"`
	ContainerPortGrpc int            `yaml:"containerPortGrpc"`
	ServicePortGrpc   int            `yaml:"servicePortGrpc"`
	Resources         resources      `yaml:"resources"`
}

type controller struct {
	Name               string             `yaml:"name"`
	Args               controllerArgs     `yaml:"args"`
	LogLevel           string             `yaml:"logLevel"`
	ContainerPort      int                `yaml:"containerPort"`
	ReadinessProbe     probe              `yaml:"readinessProbe"`
	LivenessProbe      probe              `yaml:"livenessProbe"`
	Service            service            `yaml:"service"`
	Resources          resources          `yaml:"resources"`
	ServiceAccount     serviceAccount     `yaml:"serviceAccount"`
	ClusterAdminAccess clusterAdminAccess `yaml:"clusterAdminAccess"`
}

type controllerArgs struct {
	StatusProcessors    string `yaml:"statusProcessors"`
	OperationProcessors string `yaml:"operationProcessors"`
	AppResyncPeriod     string `yaml:"appResyncPeriod"`
}

type ingress struct {
	Enabled     bool              `yaml:"enabled"`
	Annotations map[string]string `yaml:"annotations"`
	Labels      map[string]string `yaml:"labels"`
	Hosts       []string          `yaml:"hosts"`
	Paths       []string          `yaml:"paths"`
	ExtraPaths  []path            `yaml:"extraPaths"`
	HTTPS       bool              `yaml:"https"`
}

type path struct {
	Path    string  `yaml:"path"`
	Backend backend `yaml:"backend"`
}

type ingressGrpc struct {
	Enabled bool `yaml:"enabled"`
}

type backend struct {
	ServiceName string `yaml:"serviceName"`
	ServicePort string `yaml:"servicePort"`
}

type autoscaling struct {
	Enabled                           bool `yaml:"enabled"`
	MinReplicas                       int  `yaml:"minReplicas"`
	MaxReplicas                       int  `yaml:"maxReplicas"`
	TargetCPUUtilizationPercentage    int  `yaml:"targetCPUUtilizationPercentage"`
	TargetMemoryUtilizationPercentage int  `yaml:"targetMemoryUtilizationPercentage"`
}

type route struct {
	Enabled bool `yaml:"enabled"`
}

type certificate struct {
	Enabled bool `yaml:"enabled"`
}

type volumeMounts struct {
	Name      string `yaml:"name"`
	MountPath string `yaml:"mountPath"`
}

type volumes struct {
	Name     string            `yaml:"name"`
	EmptyDir map[string]string `yaml:"emptyDir"`
}

type serviceAccount struct {
	Create      bool              `yaml:"create"`
	Name        string            `yaml:"name"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

type global struct {
	Image           image           `yaml:"image"`
	SecurityContext securityContext `yaml:"securityContext"`
}

type image struct {
	Repository string `yaml:"repository"`
	Tag        string `yaml:"tag"`
	PullPolicy string `yaml:"imagePullPolicy"`
}

type securityContext struct {
	RunAsNonRoot bool `yaml:"runAsNonRoot"`
	RunAsGroup   int  `yaml:"runAsGroup"`
	RunAsUser    int  `yaml:"runAsUser"`
	FsGroup      int  `yaml:"fsGroup"`
}

type resources struct {
	Limits   resourceEntry `yaml:"limits"`
	Requests resourceEntry `yaml:"requests"`
}

type resourceEntry struct {
	CPU    string `yaml:"cpu"`
	Memory string `yaml:"memory"`
}

type probe struct {
	FailureThreshold    int `yaml:"failureThreshold"`
	InitialDelaySeconds int `yaml:"initialDelaySeconds"`
	PeriodSeconds       int `yaml:"periodSeconds"`
	SuccessThreshold    int `yaml:"successThreshold"`
	TimeoutSeconds      int `yaml:"timeoutSeconds"`
}

type service struct {
	Annotations map[string]string `yaml:"annotations,omitempty"`
	Labels      map[string]string `yaml:"labels,omitempty"`
	Port        int               `yaml:"port"`
	PortName    string            `yaml:"portName"`
}

type metrics struct {
	Enabled bool `yaml:"enabled"`
}

type clusterAdminAccess struct {
	Enabled bool `yaml:"enabled"`
}
