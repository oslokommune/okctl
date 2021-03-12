package core

import (
	"net/http"

	logmd "github.com/oslokommune/okctl/pkg/middleware/logger"

	"github.com/go-chi/chi"
	"github.com/go-kit/kit/endpoint"
	kit "github.com/go-kit/kit/transport/http"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/transport"
	"github.com/sirupsen/logrus"
)

// Endpoints defines all available endpoints
type Endpoints struct {
	CreateCluster                                 endpoint.Endpoint
	DeleteCluster                                 endpoint.Endpoint
	CreateVpc                                     endpoint.Endpoint
	DeleteVpc                                     endpoint.Endpoint
	CreateExternalSecretsPolicy                   endpoint.Endpoint
	CreateExternalSecretsServiceAccount           endpoint.Endpoint
	CreateExternalSecretsHelmChart                endpoint.Endpoint
	CreateAlbIngressControllerServiceAccount      endpoint.Endpoint
	CreateAlbIngressControllerPolicy              endpoint.Endpoint
	CreateAlbIngressControllerHelmChart           endpoint.Endpoint
	CreateExternalDNSPolicy                       endpoint.Endpoint
	CreateExternalDNSServiceAccount               endpoint.Endpoint
	CreateExternalDNSKubeDeployment               endpoint.Endpoint
	CreateHostedZone                              endpoint.Endpoint
	CreateCertificate                             endpoint.Endpoint
	CreateSecret                                  endpoint.Endpoint
	DeleteSecret                                  endpoint.Endpoint
	CreateArgoCD                                  endpoint.Endpoint
	CreateExternalSecrets                         endpoint.Endpoint
	DeleteExternalSecretsPolicy                   endpoint.Endpoint
	DeleteAlbIngressControllerPolicy              endpoint.Endpoint
	DeleteExternalDNSPolicy                       endpoint.Endpoint
	DeleteHostedZone                              endpoint.Endpoint
	DeleteExternalSecretsServiceAccount           endpoint.Endpoint
	DeleteAlbIngressControllerServiceAccount      endpoint.Endpoint
	DeleteExternalDNSServiceAccount               endpoint.Endpoint
	CreateIdentityPool                            endpoint.Endpoint
	CreateIdentityPoolClient                      endpoint.Endpoint
	CreateIdentityPoolUser                        endpoint.Endpoint
	DeleteIdentityPool                            endpoint.Endpoint
	DeleteIdentityPoolClient                      endpoint.Endpoint
	CreateAWSLoadBalancerControllerServiceAccount endpoint.Endpoint
	DeleteAWSLoadBalancerControllerServiceAccount endpoint.Endpoint
	CreateAWSLoadBalancerControllerPolicy         endpoint.Endpoint
	DeleteAWSLoadBalancerControllerPolicy         endpoint.Endpoint
	CreateAWSLoadBalancerControllerHelmChart      endpoint.Endpoint
	DeleteCertificate                             endpoint.Endpoint
	DeleteNamespace                               endpoint.Endpoint
	DeleteCognitoCertificate                      endpoint.Endpoint
	CreateAutoscalerHelmChart                     endpoint.Endpoint
	CreateAutoscalerServiceAccount                endpoint.Endpoint
	DeleteAutoscalerServiceAccount                endpoint.Endpoint
	CreateAutoscalerPolicy                        endpoint.Endpoint
	DeleteAutoscalerPolicy                        endpoint.Endpoint
	CreateBlockstoragePolicy                      endpoint.Endpoint
	DeleteBlockstoragePolicy                      endpoint.Endpoint
	CreateBlockstorageServiceAccount              endpoint.Endpoint
	DeleteBlockstorageServiceAccount              endpoint.Endpoint
	CreateBlockstorageHelmChart                   endpoint.Endpoint
	CreateStorageClass                            endpoint.Endpoint
	CreateKubePrometheusStack                     endpoint.Endpoint
	CreateLokiHelmChart                           endpoint.Endpoint
	DeleteExternalSecrets                         endpoint.Endpoint
	CreatePromtailHelmChart                       endpoint.Endpoint
	CreateConfigMap                               endpoint.Endpoint
	DeleteConfigMap                               endpoint.Endpoint
	ScaleDeployment                               endpoint.Endpoint
	CreateHelmRelease                             endpoint.Endpoint
	DeleteHelmRelease                             endpoint.Endpoint
	CreatePolicy                                  endpoint.Endpoint
	DeletePolicy                                  endpoint.Endpoint
	CreateServiceAccount                          endpoint.Endpoint
	DeleteServiceAccount                          endpoint.Endpoint
	CreateNamespace                               endpoint.Endpoint
	CreatePostgresDatabase                        endpoint.Endpoint
	DeletePostgresDatabase                        endpoint.Endpoint
}

// MakeEndpoints returns the endpoints initialised with their
// corresponding service
// nolint: funlen
func MakeEndpoints(s Services) Endpoints {
	return Endpoints{
		CreateCluster:                                 makeCreateClusterEndpoint(s.Cluster),
		DeleteCluster:                                 makeDeleteClusterEndpoint(s.Cluster),
		CreateVpc:                                     makeCreateVpcEndpoint(s.Vpc),
		DeleteVpc:                                     makeDeleteVpcEndpoint(s.Vpc),
		CreateExternalSecretsPolicy:                   makeCreateExternalSecretsPolicyEndpoint(s.ManagedPolicy),
		CreateExternalSecretsServiceAccount:           makeCreateExternalSecretsServiceAccountEndpoint(s.ServiceAccount),
		CreateExternalSecretsHelmChart:                makeCreateExternalSecretsHelmChartEndpoint(s.Helm),
		CreateAlbIngressControllerServiceAccount:      makeCreateAlbIngressControllerServiceAccountEndpoint(s.ServiceAccount),
		CreateAlbIngressControllerPolicy:              makeCreateAlbIngressControllerPolicyEndpoint(s.ManagedPolicy),
		CreateAlbIngressControllerHelmChart:           makeCreateAlbIngressControllerHelmChartEndpoint(s.Helm),
		CreateExternalDNSPolicy:                       makeCreateExternalDNSPolicyEndpoint(s.ManagedPolicy),
		CreateExternalDNSServiceAccount:               makeCreateExternalDNSServiceAccountEndpoint(s.ServiceAccount),
		CreateExternalDNSKubeDeployment:               makeCreateExternalDNSKubeDeploymentEndpoint(s.Kube),
		CreateHostedZone:                              makeCreateHostedZoneEndpoint(s.Domain),
		CreateCertificate:                             makeCreateCertificateEndpoint(s.Certificate),
		CreateSecret:                                  makeCreateSecret(s.Parameter),
		DeleteSecret:                                  makeDeleteSecret(s.Parameter),
		CreateArgoCD:                                  makeCreateArgoCD(s.Helm),
		CreateExternalSecrets:                         makeCreateExternalSecretsEndpoint(s.Kube),
		DeleteExternalSecretsPolicy:                   makeDeleteExternalSecretsPolicyEndpoint(s.ManagedPolicy),
		DeleteAlbIngressControllerPolicy:              makeDeleteAlbIngressControllerPolicyEndpoint(s.ManagedPolicy),
		DeleteExternalDNSPolicy:                       makeDeleteExternalDNSPolicyEndpoint(s.ManagedPolicy),
		DeleteHostedZone:                              makeDeleteHostedZoneEndpoint(s.Domain),
		DeleteExternalSecretsServiceAccount:           makeDeleteExternalSecretsServiceAccountEndpoint(s.ServiceAccount),
		DeleteAlbIngressControllerServiceAccount:      makeDeleteAlbIngressControllerServiceAccountEndpoint(s.ServiceAccount),
		DeleteExternalDNSServiceAccount:               makeDeleteExternalDNSServiceAccountEndpoint(s.ServiceAccount),
		CreateIdentityPool:                            makeCreateIdentityPoolEndpoint(s.IdentityManager),
		CreateIdentityPoolClient:                      makeCreateIdentityPoolClient(s.IdentityManager),
		CreateIdentityPoolUser:                        makeCreateIdentityPoolUser(s.IdentityManager),
		DeleteIdentityPool:                            makeDeleteIdentityPoolEndpoint(s.IdentityManager),
		DeleteIdentityPoolClient:                      makeDeleteIdentityPoolClientEndpoint(s.IdentityManager),
		CreateAWSLoadBalancerControllerServiceAccount: makeCreateAWSLoadBalancerControllerServiceAccountEndpoint(s.ServiceAccount),
		DeleteAWSLoadBalancerControllerServiceAccount: makeDeleteAWSLoadBalancerControllerServiceAccountEndpoint(s.ServiceAccount),
		CreateAWSLoadBalancerControllerPolicy:         makeCreateAWSLoadBalancerControllerPolicyEndpoint(s.ManagedPolicy),
		DeleteAWSLoadBalancerControllerPolicy:         makeDeleteAWSLoadBalancerControllerPolicyEndpoint(s.ManagedPolicy),
		CreateAWSLoadBalancerControllerHelmChart:      makeCreateAWSLoadBalancerControllerHelmChartEndpoint(s.Helm),
		DeleteCertificate:                             makeDeleteCertificateEndpoint(s.Certificate),
		DeleteNamespace:                               makeDeleteNamespaceEndpoint(s.Kube),
		DeleteCognitoCertificate:                      makeDeleteCognitoCertificateEndpoint(s.Certificate),
		CreateAutoscalerHelmChart:                     makeCreateAutoscalerHelmChartEndpoint(s.Helm),
		CreateAutoscalerServiceAccount:                makeCreateAutoscalerServiceAccountEndpoint(s.ServiceAccount),
		DeleteAutoscalerServiceAccount:                makeDeleteAutoscalerServiceAccountEndpoint(s.ServiceAccount),
		CreateAutoscalerPolicy:                        makeCreateAutoscalerPolicyEndpoint(s.ManagedPolicy),
		DeleteAutoscalerPolicy:                        makeDeleteAutoscalerPolicyEndpoint(s.ManagedPolicy),
		CreateBlockstoragePolicy:                      makeCreateBlockstoragePolicyEndpoint(s.ManagedPolicy),
		DeleteBlockstoragePolicy:                      makeDeleteBlockstoragePolicyEndpoint(s.ManagedPolicy),
		CreateBlockstorageServiceAccount:              makeCreateBlockstorageServiceAccountEndpoint(s.ServiceAccount),
		DeleteBlockstorageServiceAccount:              makeDeleteBlockstorageServiceAccountEndpoint(s.ServiceAccount),
		CreateBlockstorageHelmChart:                   makeCreateBlockstorageHelmChartEndpoint(s.Helm),
		CreateStorageClass:                            makeCreateStorageClass(s.Kube),
		CreateKubePrometheusStack:                     makeCreateKubePrometheusStack(s.Helm),
		CreateLokiHelmChart:                           makeCreateLokiHelmChartEndpoint(s.Helm),
		DeleteExternalSecrets:                         makeDeleteExternalSecrets(s.Kube),
		CreatePromtailHelmChart:                       makeCreatePromtailHelmChartEndpoint(s.Helm),
		CreateConfigMap:                               makeCreateConfigMapEndpoint(s.Kube),
		DeleteConfigMap:                               makeDeleteConfigMap(s.Kube),
		ScaleDeployment:                               makeScaleDeployment(s.Kube),
		CreateHelmRelease:                             makeCreateHelmRelease(s.Helm),
		DeleteHelmRelease:                             makeDeleteHelmRelease(s.Helm),
		CreatePolicy:                                  makeCreatePolicyEndpoint(s.ManagedPolicy),
		DeletePolicy:                                  makeDeletePolicyEndpoint(s.ManagedPolicy),
		CreateServiceAccount:                          makeCreateServiceAccountEndpoint(s.ServiceAccount),
		DeleteServiceAccount:                          makeDeleteServiceAccountEndpoint(s.ServiceAccount),
		CreateNamespace:                               makeCreateNamespace(s.Kube),
		CreatePostgresDatabase:                        makeCreatePostgresDatabaseEndpoint(s.ComponentService),
		DeletePostgresDatabase:                        makeDeletePostgresDatabaseEndpoint(s.ComponentService),
	}
}

// Handlers defines http handlers for processing requests
type Handlers struct {
	CreateCluster                                 http.Handler
	DeleteCluster                                 http.Handler
	CreateVpc                                     http.Handler
	DeleteVpc                                     http.Handler
	CreateExternalSecretsPolicy                   http.Handler
	CreateExternalSecretsServiceAccount           http.Handler
	CreateExternalSecretsHelmChart                http.Handler
	CreateAlbIngressControllerServiceAccount      http.Handler
	CreateAlbIngressControllerPolicy              http.Handler
	CreateAlbIngressControllerHelmChart           http.Handler
	CreateExternalDNSPolicy                       http.Handler
	CreateExternalDNSServiceAccount               http.Handler
	CreateExternalDNSKubeDeployment               http.Handler
	CreateHostedZone                              http.Handler
	CreateCertificate                             http.Handler
	CreateSecret                                  http.Handler
	DeleteSecret                                  http.Handler
	CreateArgoCD                                  http.Handler
	CreateExternalSecrets                         http.Handler
	DeleteExternalSecretsPolicy                   http.Handler
	DeleteAlbIngressControllerPolicy              http.Handler
	DeleteExternalDNSPolicy                       http.Handler
	DeleteHostedZone                              http.Handler
	DeleteExternalSecretsServiceAccount           http.Handler
	DeleteAlbIngressControllerServiceAccount      http.Handler
	DeleteExternalDNSServiceAccount               http.Handler
	CreateIdentityPool                            http.Handler
	CreateIdentityPoolClient                      http.Handler
	CreateIdentityPoolUser                        http.Handler
	DeleteIdentityPool                            http.Handler
	DeleteIdentityPoolClient                      http.Handler
	CreateAWSLoadBalancerControllerServiceAccount http.Handler
	DeleteAWSLoadBalancerControllerServiceAccount http.Handler
	CreateAWSLoadBalancerControllerPolicy         http.Handler
	DeleteAWSLoadBalancerControllerPolicy         http.Handler
	CreateAWSLoadBalancerControllerHelmChart      http.Handler
	DeleteCertificate                             http.Handler
	DeleteNamespace                               http.Handler
	DeleteCognitoCertificate                      http.Handler
	CreateAutoscalerHelmChart                     http.Handler
	CreateAutoscalerServiceAccount                http.Handler
	DeleteAutoscalerServiceAccount                http.Handler
	CreateAutoscalerPolicy                        http.Handler
	DeleteAutoscalerPolicy                        http.Handler
	CreateBlockstoragePolicy                      http.Handler
	DeleteBlockstoragePolicy                      http.Handler
	CreateBlockstorageServiceAccount              http.Handler
	DeleteBlockstorageServiceAccount              http.Handler
	CreateBlockstorageHelmChart                   http.Handler
	CreateStorageClass                            http.Handler
	CreateKubePrometheusStack                     http.Handler
	CreateLokiHelmChart                           http.Handler
	DeleteExternalSecrets                         http.Handler
	CreatePromtailHelmChart                       http.Handler
	CreateConfigMap                               http.Handler
	DeleteConfigMap                               http.Handler
	ScaleDeployment                               http.Handler
	CreateHelmRelease                             http.Handler
	DeleteHelmRelease                             http.Handler
	CreatePolicy                                  http.Handler
	DeletePolicy                                  http.Handler
	CreateServiceAccount                          http.Handler
	DeleteServiceAccount                          http.Handler
	CreateNamespace                               http.Handler
	CreatePostgresDatabase                        http.Handler
	DeletePostgresDatabase                        http.Handler
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
// nolint: funlen
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
		CreateCluster:                                 newServer(endpoints.CreateCluster, decodeClusterCreateRequest),
		DeleteCluster:                                 newServer(endpoints.DeleteCluster, decodeClusterDeleteRequest),
		CreateVpc:                                     newServer(endpoints.CreateVpc, decodeVpcCreateRequest),
		DeleteVpc:                                     newServer(endpoints.DeleteVpc, decodeVpcDeleteRequest),
		CreateExternalSecretsPolicy:                   newServer(endpoints.CreateExternalSecretsPolicy, decodeCreateExternalSecretsPolicyRequest),
		CreateExternalSecretsServiceAccount:           newServer(endpoints.CreateExternalSecretsServiceAccount, decodeCreateExternalSecretsServiceAccount),
		CreateExternalSecretsHelmChart:                newServer(endpoints.CreateExternalSecretsHelmChart, decodeCreateExternalSecretsHelmChart),
		CreateAlbIngressControllerServiceAccount:      newServer(endpoints.CreateAlbIngressControllerServiceAccount, decodeCreateAlbIngressControllerServiceAccount),
		CreateAlbIngressControllerPolicy:              newServer(endpoints.CreateAlbIngressControllerPolicy, decodeCreateAlbIngressControllerPolicyRequest),
		CreateAlbIngressControllerHelmChart:           newServer(endpoints.CreateAlbIngressControllerHelmChart, decodeCreateAlbIngressControllerHelmChart),
		CreateExternalDNSPolicy:                       newServer(endpoints.CreateExternalDNSPolicy, decodeCreateExternalDNSPolicyRequest),
		CreateExternalDNSServiceAccount:               newServer(endpoints.CreateExternalDNSServiceAccount, decodeCreateExternalDNSServiceAccount),
		CreateExternalDNSKubeDeployment:               newServer(endpoints.CreateExternalDNSKubeDeployment, decodeCreateExternalDNSKubeDeployment),
		CreateHostedZone:                              newServer(endpoints.CreateHostedZone, decodeCreateHostedZone),
		CreateCertificate:                             newServer(endpoints.CreateCertificate, decodeCreateCertificate),
		CreateSecret:                                  newServer(endpoints.CreateSecret, decodeCreateSecret),
		DeleteSecret:                                  newServer(endpoints.DeleteSecret, decodeDeleteSecret),
		CreateArgoCD:                                  newServer(endpoints.CreateArgoCD, decodeCreateArgoCD),
		CreateExternalSecrets:                         newServer(endpoints.CreateExternalSecrets, decodeCreateExternalSecrets),
		DeleteExternalSecretsPolicy:                   newServer(endpoints.DeleteExternalSecretsPolicy, decodeIDRequest),
		DeleteAlbIngressControllerPolicy:              newServer(endpoints.DeleteAlbIngressControllerPolicy, decodeIDRequest),
		DeleteExternalDNSPolicy:                       newServer(endpoints.DeleteExternalDNSPolicy, decodeIDRequest),
		DeleteHostedZone:                              newServer(endpoints.DeleteHostedZone, decodeDeleteHostedZone),
		DeleteExternalSecretsServiceAccount:           newServer(endpoints.DeleteExternalSecretsServiceAccount, decodeIDRequest),
		DeleteAlbIngressControllerServiceAccount:      newServer(endpoints.DeleteAlbIngressControllerServiceAccount, decodeIDRequest),
		DeleteExternalDNSServiceAccount:               newServer(endpoints.DeleteExternalDNSServiceAccount, decodeIDRequest),
		CreateIdentityPool:                            newServer(endpoints.CreateIdentityPool, decodeCreateIdentityPool),
		CreateIdentityPoolClient:                      newServer(endpoints.CreateIdentityPoolClient, decodeCreateIdentityPoolClient),
		CreateIdentityPoolUser:                        newServer(endpoints.CreateIdentityPoolUser, decodeCreateIdentityPoolUser),
		DeleteIdentityPool:                            newServer(endpoints.DeleteIdentityPool, decodeDeleteIdentityPool),
		DeleteIdentityPoolClient:                      newServer(endpoints.DeleteIdentityPoolClient, decodeDeleteIdentityPoolClient),
		CreateAWSLoadBalancerControllerServiceAccount: newServer(endpoints.CreateAWSLoadBalancerControllerServiceAccount, decodeCreateAWSLoadBalancerControllerServiceAccount),
		DeleteAWSLoadBalancerControllerServiceAccount: newServer(endpoints.DeleteAWSLoadBalancerControllerServiceAccount, decodeIDRequest),
		CreateAWSLoadBalancerControllerPolicy:         newServer(endpoints.CreateAWSLoadBalancerControllerPolicy, decodeCreateAWSLoadBalancerControllerPolicyRequest),
		DeleteAWSLoadBalancerControllerPolicy:         newServer(endpoints.DeleteAWSLoadBalancerControllerPolicy, decodeIDRequest),
		CreateAWSLoadBalancerControllerHelmChart:      newServer(endpoints.CreateAWSLoadBalancerControllerHelmChart, decodeCreateAWSLoadBalancerControllerHelmChart),
		DeleteCertificate:                             newServer(endpoints.DeleteCertificate, decodeDeleteCertificate),
		DeleteNamespace:                               newServer(endpoints.DeleteNamespace, decodeDeleteNamespace),
		DeleteCognitoCertificate:                      newServer(endpoints.DeleteCognitoCertificate, decodeDeleteCognitoCertificate),
		CreateAutoscalerHelmChart:                     newServer(endpoints.CreateAutoscalerHelmChart, decodeCreateAutoscalerHelmChart),
		CreateAutoscalerServiceAccount:                newServer(endpoints.CreateAutoscalerServiceAccount, decodeCreateAutoscalerServiceAccount),
		DeleteAutoscalerServiceAccount:                newServer(endpoints.DeleteAutoscalerServiceAccount, decodeIDRequest),
		CreateAutoscalerPolicy:                        newServer(endpoints.CreateAutoscalerPolicy, decodeCreateAutoscalerPolicy),
		DeleteAutoscalerPolicy:                        newServer(endpoints.DeleteAutoscalerPolicy, decodeIDRequest),
		CreateBlockstoragePolicy:                      newServer(endpoints.CreateBlockstoragePolicy, decodeCreateBlockstoragePolicy),
		DeleteBlockstoragePolicy:                      newServer(endpoints.DeleteBlockstoragePolicy, decodeIDRequest),
		CreateBlockstorageServiceAccount:              newServer(endpoints.CreateBlockstorageServiceAccount, decodeCreateBlockstorageServiceAccount),
		DeleteBlockstorageServiceAccount:              newServer(endpoints.DeleteBlockstorageServiceAccount, decodeIDRequest),
		CreateBlockstorageHelmChart:                   newServer(endpoints.CreateBlockstorageHelmChart, decodeCreateBlockstorageHelmChart),
		CreateStorageClass:                            newServer(endpoints.CreateStorageClass, decodeCreateStorageClass),
		CreateKubePrometheusStack:                     newServer(endpoints.CreateKubePrometheusStack, decodeCreateKubePrometheusStackHelmChart),
		CreateLokiHelmChart:                           newServer(endpoints.CreateLokiHelmChart, decodeCreateLokiHelmChart),
		DeleteExternalSecrets:                         newServer(endpoints.DeleteExternalSecrets, decodeDeleteExternalSecrets),
		CreatePromtailHelmChart:                       newServer(endpoints.CreatePromtailHelmChart, decodeCreatePromtailHelmChart),
		CreateConfigMap:                               newServer(endpoints.CreateConfigMap, decodeCreateConfigMap),
		DeleteConfigMap:                               newServer(endpoints.DeleteConfigMap, decodeDeleteConfigMap),
		ScaleDeployment:                               newServer(endpoints.ScaleDeployment, decodeScaleDeployment),
		CreateHelmRelease:                             newServer(endpoints.CreateHelmRelease, decodeCreateHelmRelease),
		DeleteHelmRelease:                             newServer(endpoints.DeleteHelmRelease, decodeDeleteHelmRelease),
		CreatePolicy:                                  newServer(endpoints.CreatePolicy, decodeStructRequest(&api.CreatePolicyOpts{})),
		DeletePolicy:                                  newServer(endpoints.DeletePolicy, decodeStructRequest(&api.DeletePolicyOpts{})),
		CreateServiceAccount:                          newServer(endpoints.CreateServiceAccount, decodeStructRequest(&api.CreateServiceAccountOpts{})),
		DeleteServiceAccount:                          newServer(endpoints.DeleteServiceAccount, decodeStructRequest(&api.DeleteServiceAccountOpts{})),
		CreateNamespace:                               newServer(endpoints.CreateNamespace, decodeStructRequest(&api.CreateNamespaceOpts{})),
		CreatePostgresDatabase:                        newServer(endpoints.CreatePostgresDatabase, decodeStructRequest(&api.CreatePostgresDatabaseOpts{})),
		DeletePostgresDatabase:                        newServer(endpoints.DeletePostgresDatabase, decodeStructRequest(&api.DeletePostgresDatabaseOpts{})),
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
			r.Method(http.MethodPost, "/", handlers.CreatePolicy)
			r.Method(http.MethodDelete, "/", handlers.DeletePolicy)
			r.Route("/externalsecrets", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateExternalSecretsPolicy)
				r.Method(http.MethodDelete, "/", handlers.DeleteExternalSecretsPolicy)
			})
			r.Route("/albingresscontroller", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateAlbIngressControllerPolicy)
				r.Method(http.MethodDelete, "/", handlers.DeleteAlbIngressControllerPolicy)
			})
			r.Route("/awsloadbalancercontroller", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateAWSLoadBalancerControllerPolicy)
				r.Method(http.MethodDelete, "/", handlers.DeleteAWSLoadBalancerControllerPolicy)
			})
			r.Route("/externaldns", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateExternalDNSPolicy)
				r.Method(http.MethodDelete, "/", handlers.DeleteExternalDNSPolicy)
			})
			r.Route("/autoscaler", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateAutoscalerPolicy)
				r.Method(http.MethodDelete, "/", handlers.DeleteAutoscalerPolicy)
			})
			r.Route("/blockstorage", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateBlockstoragePolicy)
				r.Method(http.MethodDelete, "/", handlers.DeleteBlockstoragePolicy)
			})
		})
		r.Route("/serviceaccounts", func(r chi.Router) {
			r.Method(http.MethodPost, "/", handlers.CreateServiceAccount)
			r.Method(http.MethodDelete, "/", handlers.DeleteServiceAccount)
			r.Route("/externalsecrets", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateExternalSecretsServiceAccount)
				r.Method(http.MethodDelete, "/", handlers.DeleteExternalSecretsServiceAccount)
			})
			r.Route("/albingresscontroller", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateAlbIngressControllerServiceAccount)
				r.Method(http.MethodDelete, "/", handlers.DeleteAlbIngressControllerServiceAccount)
			})
			r.Route("/awsloadbalancercontroller", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateAWSLoadBalancerControllerServiceAccount)
				r.Method(http.MethodDelete, "/", handlers.DeleteAWSLoadBalancerControllerServiceAccount)
			})
			r.Route("/externaldns", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateExternalDNSServiceAccount)
				r.Method(http.MethodDelete, "/", handlers.DeleteExternalDNSServiceAccount)
			})
			r.Route("/autoscaler", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateAutoscalerServiceAccount)
				r.Method(http.MethodDelete, "/", handlers.DeleteAutoscalerServiceAccount)
			})
			r.Route("/blockstorage", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateBlockstorageServiceAccount)
				r.Method(http.MethodDelete, "/", handlers.DeleteBlockstorageServiceAccount)
			})
		})
		r.Route("/helm", func(r chi.Router) {
			r.Route("/externalsecrets", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateExternalSecretsHelmChart)
			})
			r.Route("/albingresscontroller", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateAlbIngressControllerHelmChart)
			})
			r.Route("/awsloadbalancercontroller", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateAWSLoadBalancerControllerHelmChart)
			})
			r.Route("/argocd", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateArgoCD)
			})
			r.Route("/autoscaler", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateAutoscalerHelmChart)
			})
			r.Route("/blockstorage", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateBlockstorageHelmChart)
			})
			r.Route("/kubepromstack", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateKubePrometheusStack)
			})
			r.Route("/loki", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateLokiHelmChart)
			})
			r.Route("/releases", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateHelmRelease)
				r.Method(http.MethodDelete, "/", handlers.DeleteHelmRelease)
			})
			r.Route("/promtail", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreatePromtailHelmChart)
			})
		})
		r.Route("/kube", func(r chi.Router) {
			r.Route("/externaldns", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateExternalDNSKubeDeployment)
			})
			r.Route("/externalsecrets", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateExternalSecrets)
				r.Method(http.MethodDelete, "/", handlers.DeleteExternalSecrets)
			})
			r.Route("/namespaces", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateNamespace)
				r.Method(http.MethodDelete, "/", handlers.DeleteNamespace)
			})
			r.Route("/storageclasses", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateStorageClass)
			})
			r.Route("/configmaps", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateConfigMap)
				r.Method(http.MethodDelete, "/", handlers.DeleteConfigMap)
			})
			r.Route("/scale", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.ScaleDeployment)
			})
		})
		r.Route("/domains", func(r chi.Router) {
			r.Route("/hostedzones", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateHostedZone)
				r.Method(http.MethodDelete, "/", handlers.DeleteHostedZone)
			})
		})
		r.Route("/certificates", func(r chi.Router) {
			r.Method(http.MethodPost, "/", handlers.CreateCertificate)
			r.Method(http.MethodDelete, "/", handlers.DeleteCertificate)
			r.Route("/cognito", func(r chi.Router) {
				r.Method(http.MethodDelete, "/", handlers.DeleteCognitoCertificate)
			})
		})
		r.Route("/parameters", func(r chi.Router) {
			r.Route("/secret", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateSecret)
				r.Method(http.MethodDelete, "/", handlers.DeleteSecret)
			})
		})
		r.Route("/identitymanagers", func(r chi.Router) {
			r.Route("/pools", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateIdentityPool)
				r.Method(http.MethodDelete, "/", handlers.DeleteIdentityPool)
				r.Route("/clients", func(r chi.Router) {
					r.Method(http.MethodPost, "/", handlers.CreateIdentityPoolClient)
					r.Method(http.MethodDelete, "/", handlers.DeleteIdentityPoolClient)
				})
				r.Route("/users", func(r chi.Router) {
					r.Method(http.MethodPost, "/", handlers.CreateIdentityPoolUser)
				})
			})
		})
		r.Route("/components", func(r chi.Router) {
			r.Route("/postgres", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreatePostgresDatabase)
				r.Method(http.MethodDelete, "/", handlers.DeletePostgresDatabase)
			})
		})
	})

	return r
}

// Services defines all available services
type Services struct {
	Cluster          api.ClusterService
	Vpc              api.VpcService
	ManagedPolicy    api.ManagedPolicyService
	ServiceAccount   api.ServiceAccountService
	Helm             api.HelmService
	Kube             api.KubeService
	Domain           api.DomainService
	Certificate      api.CertificateService
	Parameter        api.ParameterService
	IdentityManager  api.IdentityManagerService
	ComponentService api.ComponentService
}

// EndpointOption makes it easy to enable and disable the endpoint
// middlewares
type EndpointOption func(Endpoints) Endpoints

const (
	clusterTag                   = "clusterService"
	vpcTag                       = "vpc"
	managedPoliciesTag           = "managedPolicies"
	externalSecretsTag           = "externalSecrets"
	serviceAccountsTag           = "serviceAccounts"
	helmTag                      = "helm"
	albIngressControllerTag      = "albingresscontroller"
	awsLoadBalancerControllerTag = "awsloadbalancercontroller"
	externalDNSTag               = "externaldns"
	kubeTag                      = "kube"
	domainTag                    = "domain"
	hostedZoneTag                = "hostedZone"
	certificateTag               = "certificate"
	parameterTag                 = "parameter"
	secretTag                    = "secret"
	argocdTag                    = "argocd"
	identityManagerTag           = "identitymanager"
	identityPoolTag              = "identitypool"
	identityPoolClientTag        = "identitypoolclient"
	identityPoolUserTag          = "identitypooluser"
	namespaceTag                 = "namespace"
	cognitoTag                   = "cognito"
	autoscalerTag                = "autoscaler"
	blockstorageTag              = "blockstorage"
	storageclassTag              = "storageclass"
	kubePrometheusStackTag       = "kubeprometheusstack"
	lokiTag                      = "loki"
	releasesTag                  = "releases"
	promtailTag                  = "promtail"
	configMapTag                 = "configmap"
	scaleTag                     = "scale"
	postgresTag                  = "postgres"
	componentsTag                = "components"
)

// InstrumentEndpoints adds instrumentation to the endpoints
// nolint: lll funlen
func InstrumentEndpoints(logger *logrus.Logger) EndpointOption {
	return func(endpoints Endpoints) Endpoints {
		return Endpoints{
			CreateCluster:                                 logmd.Logging(logger, "create", clusterTag)(endpoints.CreateCluster),
			DeleteCluster:                                 logmd.Logging(logger, "delete", clusterTag)(endpoints.DeleteCluster),
			CreateVpc:                                     logmd.Logging(logger, "create", vpcTag)(endpoints.CreateVpc),
			DeleteVpc:                                     logmd.Logging(logger, "delete", vpcTag)(endpoints.DeleteVpc),
			CreateExternalSecretsPolicy:                   logmd.Logging(logger, "create", managedPoliciesTag, externalSecretsTag)(endpoints.CreateExternalSecretsPolicy),
			CreateExternalSecretsServiceAccount:           logmd.Logging(logger, "create", serviceAccountsTag, externalSecretsTag)(endpoints.CreateExternalSecretsServiceAccount),
			CreateExternalSecretsHelmChart:                logmd.Logging(logger, "create", helmTag, externalSecretsTag)(endpoints.CreateExternalSecretsHelmChart),
			CreateAlbIngressControllerServiceAccount:      logmd.Logging(logger, "create", serviceAccountsTag, albIngressControllerTag)(endpoints.CreateAlbIngressControllerServiceAccount),
			CreateAlbIngressControllerPolicy:              logmd.Logging(logger, "create", managedPoliciesTag, albIngressControllerTag)(endpoints.CreateAlbIngressControllerPolicy),
			CreateAlbIngressControllerHelmChart:           logmd.Logging(logger, "create", helmTag, albIngressControllerTag)(endpoints.CreateAlbIngressControllerHelmChart),
			CreateExternalDNSPolicy:                       logmd.Logging(logger, "create", managedPoliciesTag, externalDNSTag)(endpoints.CreateExternalDNSPolicy),
			CreateExternalDNSServiceAccount:               logmd.Logging(logger, "create", serviceAccountsTag, externalDNSTag)(endpoints.CreateExternalDNSServiceAccount),
			CreateExternalDNSKubeDeployment:               logmd.Logging(logger, "create", kubeTag, externalDNSTag)(endpoints.CreateExternalDNSKubeDeployment),
			CreateHostedZone:                              logmd.Logging(logger, "create", domainTag, hostedZoneTag)(endpoints.CreateHostedZone),
			CreateCertificate:                             logmd.Logging(logger, "create", certificateTag)(endpoints.CreateCertificate),
			CreateSecret:                                  logmd.Logging(logger, "create", parameterTag, secretTag)(endpoints.CreateSecret),
			DeleteSecret:                                  logmd.Logging(logger, "delete", parameterTag, secretTag)(endpoints.DeleteSecret),
			CreateArgoCD:                                  logmd.Logging(logger, "create", helmTag, argocdTag)(endpoints.CreateArgoCD),
			CreateExternalSecrets:                         logmd.Logging(logger, "create", kubeTag, externalSecretsTag)(endpoints.CreateExternalSecrets),
			DeleteExternalSecretsPolicy:                   logmd.Logging(logger, "delete", managedPoliciesTag, externalSecretsTag)(endpoints.DeleteExternalSecretsPolicy),
			DeleteAlbIngressControllerPolicy:              logmd.Logging(logger, "delete", managedPoliciesTag, albIngressControllerTag)(endpoints.DeleteAlbIngressControllerPolicy),
			DeleteExternalDNSPolicy:                       logmd.Logging(logger, "delete", managedPoliciesTag, externalDNSTag)(endpoints.DeleteExternalDNSPolicy),
			DeleteHostedZone:                              logmd.Logging(logger, "delete", domainTag, hostedZoneTag)(endpoints.DeleteHostedZone),
			DeleteExternalSecretsServiceAccount:           logmd.Logging(logger, "delete", serviceAccountsTag, externalSecretsTag)(endpoints.DeleteExternalSecretsServiceAccount),
			DeleteAlbIngressControllerServiceAccount:      logmd.Logging(logger, "delete", serviceAccountsTag, albIngressControllerTag)(endpoints.DeleteAlbIngressControllerServiceAccount),
			DeleteExternalDNSServiceAccount:               logmd.Logging(logger, "delete", serviceAccountsTag, externalDNSTag)(endpoints.DeleteExternalDNSServiceAccount),
			CreateIdentityPool:                            logmd.Logging(logger, "create", identityManagerTag, identityPoolTag)(endpoints.CreateIdentityPool),
			CreateIdentityPoolClient:                      logmd.Logging(logger, "create", identityManagerTag, identityPoolTag, identityPoolClientTag)(endpoints.CreateIdentityPoolClient),
			CreateIdentityPoolUser:                        logmd.Logging(logger, "create", identityManagerTag, identityPoolTag, identityPoolUserTag)(endpoints.CreateIdentityPoolUser),
			DeleteIdentityPool:                            logmd.Logging(logger, "delete", identityManagerTag, identityPoolTag)(endpoints.DeleteIdentityPool),
			DeleteIdentityPoolClient:                      logmd.Logging(logger, "delete", identityManagerTag, identityPoolClientTag)(endpoints.DeleteIdentityPoolClient),
			CreateAWSLoadBalancerControllerServiceAccount: logmd.Logging(logger, "create", serviceAccountsTag, awsLoadBalancerControllerTag)(endpoints.CreateAWSLoadBalancerControllerServiceAccount),
			DeleteAWSLoadBalancerControllerServiceAccount: logmd.Logging(logger, "delete", serviceAccountsTag, awsLoadBalancerControllerTag)(endpoints.DeleteAWSLoadBalancerControllerServiceAccount),
			CreateAWSLoadBalancerControllerPolicy:         logmd.Logging(logger, "create", managedPoliciesTag, awsLoadBalancerControllerTag)(endpoints.CreateAWSLoadBalancerControllerPolicy),
			DeleteAWSLoadBalancerControllerPolicy:         logmd.Logging(logger, "delete", managedPoliciesTag, awsLoadBalancerControllerTag)(endpoints.DeleteAWSLoadBalancerControllerPolicy),
			CreateAWSLoadBalancerControllerHelmChart:      logmd.Logging(logger, "create", helmTag, awsLoadBalancerControllerTag)(endpoints.CreateAWSLoadBalancerControllerHelmChart),
			DeleteCertificate:                             logmd.Logging(logger, "delete", certificateTag)(endpoints.DeleteCertificate),
			DeleteNamespace:                               logmd.Logging(logger, "delete", kubeTag, namespaceTag)(endpoints.DeleteNamespace),
			DeleteCognitoCertificate:                      logmd.Logging(logger, "delete", certificateTag, cognitoTag)(endpoints.DeleteCognitoCertificate),
			CreateAutoscalerHelmChart:                     logmd.Logging(logger, "create", helmTag, autoscalerTag)(endpoints.CreateAutoscalerHelmChart),
			CreateAutoscalerServiceAccount:                logmd.Logging(logger, "create", serviceAccountsTag, autoscalerTag)(endpoints.CreateAutoscalerServiceAccount),
			DeleteAutoscalerServiceAccount:                logmd.Logging(logger, "delete", serviceAccountsTag, autoscalerTag)(endpoints.DeleteAutoscalerServiceAccount),
			CreateAutoscalerPolicy:                        logmd.Logging(logger, "create", managedPoliciesTag, autoscalerTag)(endpoints.CreateAutoscalerPolicy),
			DeleteAutoscalerPolicy:                        logmd.Logging(logger, "delete", managedPoliciesTag, autoscalerTag)(endpoints.DeleteAutoscalerPolicy),
			CreateBlockstoragePolicy:                      logmd.Logging(logger, "create", managedPoliciesTag, blockstorageTag)(endpoints.CreateBlockstoragePolicy),
			DeleteBlockstoragePolicy:                      logmd.Logging(logger, "delete", managedPoliciesTag, blockstorageTag)(endpoints.DeleteBlockstoragePolicy),
			CreateBlockstorageServiceAccount:              logmd.Logging(logger, "create", serviceAccountsTag, blockstorageTag)(endpoints.CreateBlockstorageServiceAccount),
			DeleteBlockstorageServiceAccount:              logmd.Logging(logger, "delete", serviceAccountsTag, blockstorageTag)(endpoints.DeleteBlockstorageServiceAccount),
			CreateBlockstorageHelmChart:                   logmd.Logging(logger, "create", helmTag, blockstorageTag)(endpoints.CreateBlockstorageHelmChart),
			CreateStorageClass:                            logmd.Logging(logger, "create", kubeTag, storageclassTag)(endpoints.CreateStorageClass),
			CreateKubePrometheusStack:                     logmd.Logging(logger, "create", helmTag, kubePrometheusStackTag)(endpoints.CreateKubePrometheusStack),
			CreateLokiHelmChart:                           logmd.Logging(logger, "create", helmTag, lokiTag)(endpoints.CreateLokiHelmChart),
			DeleteExternalSecrets:                         logmd.Logging(logger, "delete", kubeTag, externalSecretsTag)(endpoints.DeleteExternalSecrets),
			CreatePromtailHelmChart:                       logmd.Logging(logger, "create", helmTag, promtailTag)(endpoints.CreatePromtailHelmChart),
			CreateConfigMap:                               logmd.Logging(logger, "create", kubeTag, configMapTag)(endpoints.CreateConfigMap),
			DeleteConfigMap:                               logmd.Logging(logger, "delete", kubeTag, configMapTag)(endpoints.DeleteConfigMap),
			ScaleDeployment:                               logmd.Logging(logger, "create", kubeTag, scaleTag)(endpoints.ScaleDeployment),
			CreateHelmRelease:                             logmd.Logging(logger, "create", helmTag, releasesTag)(endpoints.CreateHelmRelease),
			DeleteHelmRelease:                             logmd.Logging(logger, "delete", helmTag, releasesTag)(endpoints.DeleteHelmRelease),
			CreatePolicy:                                  logmd.Logging(logger, "create", managedPoliciesTag)(endpoints.CreatePolicy),
			DeletePolicy:                                  logmd.Logging(logger, "delete", managedPoliciesTag)(endpoints.DeletePolicy),
			CreateServiceAccount:                          logmd.Logging(logger, "create", serviceAccountsTag)(endpoints.CreateServiceAccount),
			DeleteServiceAccount:                          logmd.Logging(logger, "delete", serviceAccountsTag)(endpoints.DeleteServiceAccount),
			CreateNamespace:                               logmd.Logging(logger, "create", kubeTag, namespaceTag)(endpoints.CreateNamespace),
			CreatePostgresDatabase:                        logmd.Logging(logger, "create", componentsTag, postgresTag)(endpoints.CreatePostgresDatabase),
			DeletePostgresDatabase:                        logmd.Logging(logger, "delete", componentsTag, postgresTag)(endpoints.DeletePostgresDatabase),
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
