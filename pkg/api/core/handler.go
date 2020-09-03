package core

import (
	"net/http"
	"strings"

	logger2 "github.com/oslokommune/okctl/pkg/middleware/logger"

	"github.com/go-chi/chi"
	"github.com/go-kit/kit/endpoint"
	kit "github.com/go-kit/kit/transport/http"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/transport"
	"github.com/sirupsen/logrus"
)

// Endpoints defines all available endpoints
type Endpoints struct {
	CreateCluster                            endpoint.Endpoint
	DeleteCluster                            endpoint.Endpoint
	CreateVpc                                endpoint.Endpoint
	DeleteVpc                                endpoint.Endpoint
	CreateExternalSecretsPolicy              endpoint.Endpoint
	CreateExternalSecretsServiceAccount      endpoint.Endpoint
	CreateExternalSecretsHelmChart           endpoint.Endpoint
	CreateAlbIngressControllerServiceAccount endpoint.Endpoint
	CreateAlbIngressControllerPolicy         endpoint.Endpoint
	CreateAlbIngressControllerHelmChart      endpoint.Endpoint
	CreateExternalDNSPolicy                  endpoint.Endpoint
	CreateExternalDNSServiceAccount          endpoint.Endpoint
	CreateExternalDNSKubeDeployment          endpoint.Endpoint
	CreateDomain                             endpoint.Endpoint
	CreateCertificate                        endpoint.Endpoint
	CreateSecret                             endpoint.Endpoint
	CreateArgoCD                             endpoint.Endpoint
	CreateExternalSecrets                    endpoint.Endpoint
}

// MakeEndpoints returns the endpoints initialised with their
// corresponding service
func MakeEndpoints(s Services) Endpoints {
	return Endpoints{
		CreateCluster:                            makeCreateClusterEndpoint(s.Cluster),
		DeleteCluster:                            makeDeleteClusterEndpoint(s.Cluster),
		CreateVpc:                                makeCreateVpcEndpoint(s.Vpc),
		DeleteVpc:                                makeDeleteVpcEndpoint(s.Vpc),
		CreateExternalSecretsPolicy:              makeCreateExternalSecretsPolicyEndpoint(s.ManagedPolicy),
		CreateExternalSecretsServiceAccount:      makeCreateExternalSecretsServiceAccountEndpoint(s.ServiceAccount),
		CreateExternalSecretsHelmChart:           makeCreateExternalSecretsHelmChartEndpoint(s.Helm),
		CreateAlbIngressControllerServiceAccount: makeCreateAlbIngressControllerServiceAccountEndpoint(s.ServiceAccount),
		CreateAlbIngressControllerPolicy:         makeCreateAlbIngressControllerPolicyEndpoint(s.ManagedPolicy),
		CreateAlbIngressControllerHelmChart:      makeCreateAlbIngressControllerHelmChartEndpoint(s.Helm),
		CreateExternalDNSPolicy:                  makeCreateExternalDNSPolicyEndpoint(s.ManagedPolicy),
		CreateExternalDNSServiceAccount:          makeCreateExternalDNSServiceAccountEndpoint(s.ServiceAccount),
		CreateExternalDNSKubeDeployment:          makeCreateExternalDNSKubeDeploymentEndpoint(s.Kube),
		CreateDomain:                             makeCreateDomainEndpoint(s.Domain),
		CreateCertificate:                        makeCreateCertificateEndpoint(s.Certificate),
		CreateSecret:                             makeCreateSecret(s.Parameter),
		CreateArgoCD:                             makeCreateArgoCD(s.Helm),
		CreateExternalSecrets:                    makeCreateExternalSecretsEndpoint(s.Kube),
	}
}

// Handlers defines http handlers for processing requests
type Handlers struct {
	CreateCluster                            http.Handler
	DeleteCluster                            http.Handler
	CreateVpc                                http.Handler
	DeleteVpc                                http.Handler
	CreateExternalSecretsPolicy              http.Handler
	CreateExternalSecretsServiceAccount      http.Handler
	CreateExternalSecretsHelmChart           http.Handler
	CreateAlbIngressControllerServiceAccount http.Handler
	CreateAlbIngressControllerPolicy         http.Handler
	CreateAlbIngressControllerHelmChart      http.Handler
	CreateExternalDNSPolicy                  http.Handler
	CreateExternalDNSServiceAccount          http.Handler
	CreateExternalDNSKubeDeployment          http.Handler
	CreateDomain                             http.Handler
	CreateCertificate                        http.Handler
	CreateSecret                             http.Handler
	CreateArgoCD                             http.Handler
	CreateExternalSecrets                    http.Handler
}

// EncodeResponseType defines a type for responses
type EncodeResponseType string

const (
	// EncodeJSONResponse encodes as json when returning response
	EncodeJSONResponse EncodeResponseType = "json"
	// EncodeYAMLResponse encodes as yaml when returning response
	EncodeYAMLResponse EncodeResponseType = "yaml"
	// EncodeTextResponse encodes as text when returning response
	EncodeTextResponse EncodeResponseType = "text"
)

// MakeHandlers returns all handlers initialised with encoders, decoders, etc
func MakeHandlers(responseType EncodeResponseType, endpoints Endpoints) *Handlers {
	var encoderFn kit.EncodeResponseFunc

	switch responseType {
	case EncodeJSONResponse:
		encoderFn = kit.EncodeJSONResponse
	case EncodeYAMLResponse:
		encoderFn = transport.EncodeYAMLResponse
	case EncodeTextResponse:
		encoderFn = transport.EncodeTextResponse
	}

	newServer := func(e endpoint.Endpoint, decodeFn kit.DecodeRequestFunc) http.Handler {
		return kit.NewServer(e, decodeFn, encoderFn)
	}

	return &Handlers{
		CreateCluster:                            newServer(endpoints.CreateCluster, decodeClusterCreateRequest),
		DeleteCluster:                            newServer(endpoints.DeleteCluster, decodeClusterDeleteRequest),
		CreateVpc:                                newServer(endpoints.CreateVpc, decodeVpcCreateRequest),
		DeleteVpc:                                newServer(endpoints.DeleteVpc, decodeVpcDeleteRequest),
		CreateExternalSecretsPolicy:              newServer(endpoints.CreateExternalSecretsPolicy, decodeCreateExternalSecretsPolicyRequest),
		CreateExternalSecretsServiceAccount:      newServer(endpoints.CreateExternalSecretsServiceAccount, decodeCreateExternalSecretsServiceAccount),
		CreateExternalSecretsHelmChart:           newServer(endpoints.CreateExternalSecretsHelmChart, decodeCreateExternalSecretsHelmChart),
		CreateAlbIngressControllerServiceAccount: newServer(endpoints.CreateAlbIngressControllerServiceAccount, decodeCreateAlbIngressControllerServiceAccount),
		CreateAlbIngressControllerPolicy:         newServer(endpoints.CreateAlbIngressControllerPolicy, decodeCreateAlbIngressControllerPolicyRequest),
		CreateAlbIngressControllerHelmChart:      newServer(endpoints.CreateAlbIngressControllerHelmChart, decodeCreateAlbIngressControllerHelmChart),
		CreateExternalDNSPolicy:                  newServer(endpoints.CreateExternalDNSPolicy, decodeCreateExternalDNSPolicyRequest),
		CreateExternalDNSServiceAccount:          newServer(endpoints.CreateExternalDNSServiceAccount, decodeCreateExternalDNSServiceAccount),
		CreateExternalDNSKubeDeployment:          newServer(endpoints.CreateExternalDNSKubeDeployment, decodeCreateExternalDNSKubeDeployment),
		CreateDomain:                             newServer(endpoints.CreateDomain, decodeCreateDomain),
		CreateCertificate:                        newServer(endpoints.CreateCertificate, decodeCreateCertificate),
		CreateSecret:                             newServer(endpoints.CreateSecret, decodeCreateSecret),
		CreateArgoCD:                             newServer(endpoints.CreateArgoCD, decodeCreateArgoCD),
		CreateExternalSecrets:                    newServer(endpoints.CreateExternalSecrets, decodeCreateExternalSecrets),
	}
}

// AttachRoutes creates a router and adds handlers to the corresponding routes
// nolint: funlen
func AttachRoutes(handlers *Handlers) http.Handler {
	r := chi.NewRouter()

	r.Route("/v1", func(r chi.Router) {
		r.Route("/clusters", func(r chi.Router) {
			r.Method(http.MethodPost, "/", handlers.CreateCluster)
			r.Method(http.MethodDelete, "/", handlers.DeleteCluster)
		})
		r.Route("/vpcs", func(r chi.Router) {
			r.Method(http.MethodPost, "/", handlers.CreateVpc)
			r.Method(http.MethodDelete, "/", handlers.DeleteVpc)
		})
		r.Route("/managedpolicies", func(r chi.Router) {
			r.Route("/externalsecrets", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateExternalSecretsPolicy)
			})
			r.Route("/albingresscontroller", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateAlbIngressControllerPolicy)
			})
			r.Route("/externaldns", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateExternalDNSPolicy)
			})
		})
		r.Route("/serviceaccounts", func(r chi.Router) {
			r.Route("/externalsecrets", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateExternalSecretsServiceAccount)
			})
			r.Route("/albingresscontroller", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateAlbIngressControllerServiceAccount)
			})
			r.Route("/externaldns", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateExternalDNSServiceAccount)
			})
		})
		r.Route("/helm", func(r chi.Router) {
			r.Route("/externalsecrets", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateExternalSecretsHelmChart)
			})
			r.Route("/albingresscontroller", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateAlbIngressControllerHelmChart)
			})
			r.Route("/argocd", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateArgoCD)
			})
		})
		r.Route("/kube", func(r chi.Router) {
			r.Route("/externaldns", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateExternalDNSKubeDeployment)
			})
			r.Route("/externalsecrets", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateExternalSecrets)
			})
		})
		r.Route("/domains", func(r chi.Router) {
			r.Method(http.MethodPost, "/", handlers.CreateDomain)
		})
		r.Route("/certificates", func(r chi.Router) {
			r.Method(http.MethodPost, "/", handlers.CreateCertificate)
		})
		r.Route("/parameters", func(r chi.Router) {
			r.Route("/secret", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateSecret)
			})
		})
	})

	return r
}

// Services defines all available services
type Services struct {
	Cluster        api.ClusterService
	Vpc            api.VpcService
	ManagedPolicy  api.ManagedPolicyService
	ServiceAccount api.ServiceAccountService
	Helm           api.HelmService
	Kube           api.KubeService
	Domain         api.DomainService
	Certificate    api.CertificateService
	Parameter      api.ParameterService
}

// EndpointOption makes it easy to enable and disable the endpoint
// middlewares
type EndpointOption func(Endpoints) Endpoints

const (
	clusterTag              = "cluster"
	vpcTag                  = "vpc"
	managedPoliciesTag      = "managedPolicies"
	externalSecretsTag      = "externalSecrets"
	serviceAccountsTag      = "serviceAccounts"
	helmTag                 = "helm"
	albIngressControllerTag = "albingresscontroller"
	externalDNSTag          = "externaldns"
	kubeTag                 = "kube"
	domainTag               = "domain"
	certificateTag          = "certificate"
	parameterTag            = "parameter"
	secretTag               = "secret"
	argocdTag               = "argocd"
)

// InstrumentEndpoints adds instrumentation to the endpoints
// nolint: lll
func InstrumentEndpoints(logger *logrus.Logger) EndpointOption {
	return func(endpoints Endpoints) Endpoints {
		return Endpoints{
			CreateCluster:                            logger2.Logging(logger, clusterTag, "create")(endpoints.CreateCluster),
			DeleteCluster:                            logger2.Logging(logger, clusterTag, "delete")(endpoints.DeleteCluster),
			CreateVpc:                                logger2.Logging(logger, vpcTag, "create")(endpoints.CreateVpc),
			DeleteVpc:                                logger2.Logging(logger, vpcTag, "delete")(endpoints.DeleteVpc),
			CreateExternalSecretsPolicy:              logger2.Logging(logger, strings.Join([]string{managedPoliciesTag, externalSecretsTag}, "/"), "create")(endpoints.CreateExternalSecretsPolicy),
			CreateExternalSecretsServiceAccount:      logger2.Logging(logger, strings.Join([]string{serviceAccountsTag, externalSecretsTag}, "/"), "create")(endpoints.CreateExternalSecretsServiceAccount),
			CreateExternalSecretsHelmChart:           logger2.Logging(logger, strings.Join([]string{helmTag, externalSecretsTag}, "/"), "create")(endpoints.CreateExternalSecretsHelmChart),
			CreateAlbIngressControllerServiceAccount: logger2.Logging(logger, strings.Join([]string{serviceAccountsTag, albIngressControllerTag}, "/"), "create")(endpoints.CreateAlbIngressControllerServiceAccount),
			CreateAlbIngressControllerPolicy:         logger2.Logging(logger, strings.Join([]string{managedPoliciesTag, albIngressControllerTag}, "/"), "create")(endpoints.CreateAlbIngressControllerPolicy),
			CreateAlbIngressControllerHelmChart:      logger2.Logging(logger, strings.Join([]string{helmTag, albIngressControllerTag}, "/"), "create")(endpoints.CreateAlbIngressControllerHelmChart),
			CreateExternalDNSPolicy:                  logger2.Logging(logger, strings.Join([]string{managedPoliciesTag, externalDNSTag}, "/"), "create")(endpoints.CreateExternalDNSPolicy),
			CreateExternalDNSServiceAccount:          logger2.Logging(logger, strings.Join([]string{serviceAccountsTag, externalDNSTag}, "/"), "create")(endpoints.CreateExternalDNSServiceAccount),
			CreateExternalDNSKubeDeployment:          logger2.Logging(logger, strings.Join([]string{kubeTag, externalDNSTag}, "/"), "create")(endpoints.CreateExternalDNSKubeDeployment),
			CreateDomain:                             logger2.Logging(logger, domainTag, "create")(endpoints.CreateDomain),
			CreateCertificate:                        logger2.Logging(logger, certificateTag, "create")(endpoints.CreateCertificate),
			CreateSecret:                             logger2.Logging(logger, strings.Join([]string{parameterTag, secretTag}, "/"), "create")(endpoints.CreateSecret),
			CreateArgoCD:                             logger2.Logging(logger, strings.Join([]string{helmTag, argocdTag}, "/"), "create")(endpoints.CreateArgoCD),
			CreateExternalSecrets:                    logger2.Logging(logger, strings.Join([]string{kubeTag, externalSecretsTag}, "/"), "create")(endpoints.CreateExternalSecrets),
		}
	}
}

// GenerateEndpoints is a convenience function for decorating endpoints
func GenerateEndpoints(s Services, opts ...EndpointOption) Endpoints {
	endpoints := MakeEndpoints(s)
	for _, opt := range opts {
		endpoints = opt(endpoints)
	}

	return endpoints
}
