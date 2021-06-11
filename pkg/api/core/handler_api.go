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
	CreateCluster                   endpoint.Endpoint
	DeleteCluster                   endpoint.Endpoint
	CreateVpc                       endpoint.Endpoint
	DeleteVpc                       endpoint.Endpoint
	CreateExternalDNSKubeDeployment endpoint.Endpoint
	CreateHostedZone                endpoint.Endpoint
	CreateCertificate               endpoint.Endpoint
	CreateSecret                    endpoint.Endpoint
	DeleteSecret                    endpoint.Endpoint
	CreateExternalSecrets           endpoint.Endpoint
	DeleteHostedZone                endpoint.Endpoint
	CreateIdentityPool              endpoint.Endpoint
	CreateIdentityPoolClient        endpoint.Endpoint
	CreateIdentityPoolUser          endpoint.Endpoint
	DeleteIdentityPool              endpoint.Endpoint
	DeleteIdentityPoolClient        endpoint.Endpoint
	DeleteIdentityPoolUser          endpoint.Endpoint
	DeleteCertificate               endpoint.Endpoint
	DeleteNamespace                 endpoint.Endpoint
	DeleteCognitoCertificate        endpoint.Endpoint
	CreateStorageClass              endpoint.Endpoint
	DeleteExternalSecrets           endpoint.Endpoint
	CreateConfigMap                 endpoint.Endpoint
	DeleteConfigMap                 endpoint.Endpoint
	ScaleDeployment                 endpoint.Endpoint
	CreateHelmRelease               endpoint.Endpoint
	DeleteHelmRelease               endpoint.Endpoint
	GetHelmRelease                  endpoint.Endpoint
	CreatePolicy                    endpoint.Endpoint
	DeletePolicy                    endpoint.Endpoint
	CreateServiceAccount            endpoint.Endpoint
	DeleteServiceAccount            endpoint.Endpoint
	CreateNamespace                 endpoint.Endpoint
	CreatePostgresDatabase          endpoint.Endpoint
	DeletePostgresDatabase          endpoint.Endpoint
	CreateS3Bucket                  endpoint.Endpoint
	DeleteS3Bucket                  endpoint.Endpoint
	CreateContainerRepository       endpoint.Endpoint
	DeleteContainerRepository       endpoint.Endpoint
}

// MakeEndpoints returns the endpoints initialised with their
// corresponding service
// nolint: funlen
func MakeEndpoints(s Services) Endpoints {
	return Endpoints{
		CreateCluster:                   makeCreateClusterEndpoint(s.Cluster),
		DeleteCluster:                   makeDeleteClusterEndpoint(s.Cluster),
		CreateVpc:                       makeCreateVpcEndpoint(s.Vpc),
		DeleteVpc:                       makeDeleteVpcEndpoint(s.Vpc),
		CreateExternalDNSKubeDeployment: makeCreateExternalDNSKubeDeploymentEndpoint(s.Kube),
		CreateHostedZone:                makeCreateHostedZoneEndpoint(s.Domain),
		CreateCertificate:               makeCreateCertificateEndpoint(s.Certificate),
		CreateSecret:                    makeCreateSecret(s.Parameter),
		DeleteSecret:                    makeDeleteSecret(s.Parameter),
		CreateExternalSecrets:           makeCreateExternalSecretsEndpoint(s.Kube),
		DeleteHostedZone:                makeDeleteHostedZoneEndpoint(s.Domain),
		CreateIdentityPool:              makeCreateIdentityPoolEndpoint(s.IdentityManager),
		CreateIdentityPoolClient:        makeCreateIdentityPoolClient(s.IdentityManager),
		CreateIdentityPoolUser:          makeCreateIdentityPoolUser(s.IdentityManager),
		DeleteIdentityPool:              makeDeleteIdentityPoolEndpoint(s.IdentityManager),
		DeleteIdentityPoolClient:        makeDeleteIdentityPoolClientEndpoint(s.IdentityManager),
		DeleteIdentityPoolUser:          makeDeleteIdentityPoolUserEndpoint(s.IdentityManager),
		DeleteCertificate:               makeDeleteCertificateEndpoint(s.Certificate),
		DeleteNamespace:                 makeDeleteNamespaceEndpoint(s.Kube),
		DeleteCognitoCertificate:        makeDeleteCognitoCertificateEndpoint(s.Certificate),
		CreateStorageClass:              makeCreateStorageClass(s.Kube),
		DeleteExternalSecrets:           makeDeleteExternalSecrets(s.Kube),
		CreateConfigMap:                 makeCreateConfigMapEndpoint(s.Kube),
		DeleteConfigMap:                 makeDeleteConfigMap(s.Kube),
		ScaleDeployment:                 makeScaleDeployment(s.Kube),
		CreateHelmRelease:               makeCreateHelmRelease(s.Helm),
		DeleteHelmRelease:               makeDeleteHelmRelease(s.Helm),
		GetHelmRelease:                  makeGetHelmRelease(s.Helm),
		CreatePolicy:                    makeCreatePolicyEndpoint(s.ManagedPolicy),
		DeletePolicy:                    makeDeletePolicyEndpoint(s.ManagedPolicy),
		CreateServiceAccount:            makeCreateServiceAccountEndpoint(s.ServiceAccount),
		DeleteServiceAccount:            makeDeleteServiceAccountEndpoint(s.ServiceAccount),
		CreateNamespace:                 makeCreateNamespace(s.Kube),
		CreatePostgresDatabase:          makeCreatePostgresDatabaseEndpoint(s.ComponentService),
		DeletePostgresDatabase:          makeDeletePostgresDatabaseEndpoint(s.ComponentService),
		CreateS3Bucket:                  makeCreateS3BucketEndpoint(s.ComponentService),
		DeleteS3Bucket:                  makeDeleteS3BucketEndpoint(s.ComponentService),
		CreateContainerRepository:       makeCreateContainerRepositoryEndpoint(s.ContainerRepositoryService),
		DeleteContainerRepository:       makeDeleteContainerRepositoryEndpoint(s.ContainerRepositoryService),
	}
}

// Handlers defines http handlers for processing requests
type Handlers struct {
	CreateCluster                   http.Handler
	DeleteCluster                   http.Handler
	CreateVpc                       http.Handler
	DeleteVpc                       http.Handler
	CreateExternalDNSKubeDeployment http.Handler
	CreateHostedZone                http.Handler
	CreateCertificate               http.Handler
	CreateSecret                    http.Handler
	DeleteSecret                    http.Handler
	CreateExternalSecrets           http.Handler
	DeleteHostedZone                http.Handler
	CreateIdentityPool              http.Handler
	CreateIdentityPoolClient        http.Handler
	CreateIdentityPoolUser          http.Handler
	DeleteIdentityPool              http.Handler
	DeleteIdentityPoolClient        http.Handler
	DeleteIdentityPoolUser          http.Handler
	DeleteCertificate               http.Handler
	DeleteNamespace                 http.Handler
	DeleteCognitoCertificate        http.Handler
	CreateStorageClass              http.Handler
	DeleteExternalSecrets           http.Handler
	CreateConfigMap                 http.Handler
	DeleteConfigMap                 http.Handler
	ScaleDeployment                 http.Handler
	CreateHelmRelease               http.Handler
	DeleteHelmRelease               http.Handler
	GetHelmRelease                  http.Handler
	CreatePolicy                    http.Handler
	DeletePolicy                    http.Handler
	CreateServiceAccount            http.Handler
	DeleteServiceAccount            http.Handler
	CreateNamespace                 http.Handler
	CreatePostgresDatabase          http.Handler
	DeletePostgresDatabase          http.Handler
	CreateS3Bucket                  http.Handler
	DeleteS3Bucket                  http.Handler
	CreateContainerRepository       http.Handler
	DeleteContainerRepository       http.Handler
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
		CreateCluster:                   newServer(endpoints.CreateCluster, decodeClusterCreateRequest),
		DeleteCluster:                   newServer(endpoints.DeleteCluster, decodeClusterDeleteRequest),
		CreateVpc:                       newServer(endpoints.CreateVpc, decodeVpcCreateRequest),
		DeleteVpc:                       newServer(endpoints.DeleteVpc, decodeVpcDeleteRequest),
		CreateExternalDNSKubeDeployment: newServer(endpoints.CreateExternalDNSKubeDeployment, decodeCreateExternalDNSKubeDeployment),
		CreateHostedZone:                newServer(endpoints.CreateHostedZone, decodeCreateHostedZone),
		CreateCertificate:               newServer(endpoints.CreateCertificate, decodeCreateCertificate),
		CreateSecret:                    newServer(endpoints.CreateSecret, decodeCreateSecret),
		DeleteSecret:                    newServer(endpoints.DeleteSecret, decodeDeleteSecret),
		CreateExternalSecrets:           newServer(endpoints.CreateExternalSecrets, decodeCreateExternalSecrets),
		DeleteHostedZone:                newServer(endpoints.DeleteHostedZone, decodeDeleteHostedZone),
		CreateIdentityPool:              newServer(endpoints.CreateIdentityPool, decodeCreateIdentityPool),
		CreateIdentityPoolClient:        newServer(endpoints.CreateIdentityPoolClient, decodeCreateIdentityPoolClient),
		CreateIdentityPoolUser:          newServer(endpoints.CreateIdentityPoolUser, decodeCreateIdentityPoolUser),
		DeleteIdentityPool:              newServer(endpoints.DeleteIdentityPool, decodeDeleteIdentityPool),
		DeleteIdentityPoolClient:        newServer(endpoints.DeleteIdentityPoolClient, decodeDeleteIdentityPoolClient),
		DeleteIdentityPoolUser:          newServer(endpoints.DeleteIdentityPoolUser, decodeDeleteIdentityPoolUser),
		DeleteCertificate:               newServer(endpoints.DeleteCertificate, decodeDeleteCertificate),
		DeleteNamespace:                 newServer(endpoints.DeleteNamespace, decodeDeleteNamespace),
		DeleteCognitoCertificate:        newServer(endpoints.DeleteCognitoCertificate, decodeDeleteCognitoCertificate),
		CreateStorageClass:              newServer(endpoints.CreateStorageClass, decodeCreateStorageClass),
		DeleteExternalSecrets:           newServer(endpoints.DeleteExternalSecrets, decodeDeleteExternalSecrets),
		CreateConfigMap:                 newServer(endpoints.CreateConfigMap, decodeCreateConfigMap),
		DeleteConfigMap:                 newServer(endpoints.DeleteConfigMap, decodeDeleteConfigMap),
		ScaleDeployment:                 newServer(endpoints.ScaleDeployment, decodeScaleDeployment),
		CreateHelmRelease:               newServer(endpoints.CreateHelmRelease, decodeCreateHelmRelease),
		DeleteHelmRelease:               newServer(endpoints.DeleteHelmRelease, decodeDeleteHelmRelease),
		GetHelmRelease:                  newServer(endpoints.GetHelmRelease, decodeGetHelmRelease),
		CreatePolicy:                    newServer(endpoints.CreatePolicy, decodeStructRequest(&api.CreatePolicyOpts{})),
		DeletePolicy:                    newServer(endpoints.DeletePolicy, decodeStructRequest(&api.DeletePolicyOpts{})),
		CreateServiceAccount:            newServer(endpoints.CreateServiceAccount, decodeStructRequest(&api.CreateServiceAccountOpts{})),
		DeleteServiceAccount:            newServer(endpoints.DeleteServiceAccount, decodeStructRequest(&api.DeleteServiceAccountOpts{})),
		CreateNamespace:                 newServer(endpoints.CreateNamespace, decodeStructRequest(&api.CreateNamespaceOpts{})),
		CreatePostgresDatabase:          newServer(endpoints.CreatePostgresDatabase, decodeStructRequest(&api.CreatePostgresDatabaseOpts{})),
		DeletePostgresDatabase:          newServer(endpoints.DeletePostgresDatabase, decodeStructRequest(&api.DeletePostgresDatabaseOpts{})),
		CreateS3Bucket:                  newServer(endpoints.CreateS3Bucket, decodeStructRequest(&api.CreateS3BucketOpts{})),
		DeleteS3Bucket:                  newServer(endpoints.DeleteS3Bucket, decodeStructRequest(&api.DeleteS3BucketOpts{})),
		CreateContainerRepository:       newServer(endpoints.CreateContainerRepository, decodeStructRequest(&api.CreateContainerRepositoryOpts{})),
		DeleteContainerRepository:       newServer(endpoints.DeleteContainerRepository, decodeStructRequest(&api.DeleteContainerRepositoryOpts{})),
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
		})
		r.Route("/serviceaccounts", func(r chi.Router) {
			r.Method(http.MethodPost, "/", handlers.CreateServiceAccount)
			r.Method(http.MethodDelete, "/", handlers.DeleteServiceAccount)
		})
		r.Route("/helm", func(r chi.Router) {
			r.Route("/releases", func(r chi.Router) {
				r.Method(http.MethodGet, "/", handlers.GetHelmRelease)
				r.Method(http.MethodPost, "/", handlers.CreateHelmRelease)
				r.Method(http.MethodDelete, "/", handlers.DeleteHelmRelease)
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
					r.Method(http.MethodDelete, "/", handlers.DeleteIdentityPoolUser)
				})
			})
		})
		r.Route("/components", func(r chi.Router) {
			r.Route("/postgres", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreatePostgresDatabase)
				r.Method(http.MethodDelete, "/", handlers.DeletePostgresDatabase)
			})
			r.Route("/s3bucket", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateS3Bucket)
				r.Method(http.MethodDelete, "/", handlers.DeleteS3Bucket)
			})
		})
		r.Route("/containerrepositories", func(r chi.Router) {
			r.Method(http.MethodPost, "/", handlers.CreateContainerRepository)
			r.Method(http.MethodDelete, "/", handlers.DeleteContainerRepository)
		})
	})

	return r
}

// Services defines all available services
type Services struct {
	Cluster                    api.ClusterService
	Vpc                        api.VpcService
	ManagedPolicy              api.ManagedPolicyService
	ServiceAccount             api.ServiceAccountService
	Helm                       api.HelmService
	Kube                       api.KubeService
	Domain                     api.DomainService
	Certificate                api.CertificateService
	Parameter                  api.ParameterService
	IdentityManager            api.IdentityManagerService
	ComponentService           api.ComponentService
	ContainerRepositoryService api.ContainerRepositoryService
}

// EndpointOption makes it easy to enable and disable the endpoint
// middlewares
type EndpointOption func(Endpoints) Endpoints

const (
	clusterTag             = "clusterService"
	vpcTag                 = "vpc"
	managedPoliciesTag     = "managedPolicies"
	externalSecretsTag     = "externalSecrets"
	serviceAccountsTag     = "serviceAccounts"
	helmTag                = "helm"
	externalDNSTag         = "externaldns"
	kubeTag                = "kube"
	domainTag              = "domain"
	hostedZoneTag          = "hostedZone"
	certificateTag         = "certificate"
	parameterTag           = "parameter"
	secretTag              = "secret"
	identityManagerTag     = "identitymanager"
	identityPoolTag        = "identitypool"
	identityPoolClientTag  = "identitypoolclient"
	identityPoolUserTag    = "identitypooluser"
	namespaceTag           = "namespace"
	cognitoTag             = "cognito"
	storageclassTag        = "storageclass"
	releasesTag            = "releases"
	configMapTag           = "configmap"
	scaleTag               = "scale"
	postgresTag            = "postgres"
	componentsTag          = "components"
	s3bucketTag            = "s3bucket"
	containerRepositoryTag = "containerrepository"
)

// InstrumentEndpoints adds instrumentation to the endpoints
// nolint: lll funlen
func InstrumentEndpoints(logger *logrus.Logger) EndpointOption {
	return func(endpoints Endpoints) Endpoints {
		return Endpoints{
			CreateCluster:                   logmd.Logging(logger, "create", clusterTag)(endpoints.CreateCluster),
			DeleteCluster:                   logmd.Logging(logger, "delete", clusterTag)(endpoints.DeleteCluster),
			CreateVpc:                       logmd.Logging(logger, "create", vpcTag)(endpoints.CreateVpc),
			DeleteVpc:                       logmd.Logging(logger, "delete", vpcTag)(endpoints.DeleteVpc),
			CreateExternalDNSKubeDeployment: logmd.Logging(logger, "create", kubeTag, externalDNSTag)(endpoints.CreateExternalDNSKubeDeployment),
			CreateHostedZone:                logmd.Logging(logger, "create", domainTag, hostedZoneTag)(endpoints.CreateHostedZone),
			CreateCertificate:               logmd.Logging(logger, "create", certificateTag)(endpoints.CreateCertificate),
			CreateSecret:                    logmd.Logging(logger, "create", parameterTag, secretTag)(endpoints.CreateSecret),
			DeleteSecret:                    logmd.Logging(logger, "delete", parameterTag, secretTag)(endpoints.DeleteSecret),
			CreateExternalSecrets:           logmd.Logging(logger, "create", kubeTag, externalSecretsTag)(endpoints.CreateExternalSecrets),
			DeleteHostedZone:                logmd.Logging(logger, "delete", domainTag, hostedZoneTag)(endpoints.DeleteHostedZone),
			CreateIdentityPool:              logmd.Logging(logger, "create", identityManagerTag, identityPoolTag)(endpoints.CreateIdentityPool),
			CreateIdentityPoolClient:        logmd.Logging(logger, "create", identityManagerTag, identityPoolTag, identityPoolClientTag)(endpoints.CreateIdentityPoolClient),
			CreateIdentityPoolUser:          logmd.Logging(logger, "create", identityManagerTag, identityPoolTag, identityPoolUserTag)(endpoints.CreateIdentityPoolUser),
			DeleteIdentityPool:              logmd.Logging(logger, "delete", identityManagerTag, identityPoolTag)(endpoints.DeleteIdentityPool),
			DeleteIdentityPoolClient:        logmd.Logging(logger, "delete", identityManagerTag, identityPoolClientTag)(endpoints.DeleteIdentityPoolClient),
			DeleteIdentityPoolUser:          logmd.Logging(logger, "delete", identityManagerTag, identityPoolTag, identityPoolUserTag)(endpoints.DeleteIdentityPoolUser),
			DeleteCertificate:               logmd.Logging(logger, "delete", certificateTag)(endpoints.DeleteCertificate),
			DeleteNamespace:                 logmd.Logging(logger, "delete", kubeTag, namespaceTag)(endpoints.DeleteNamespace),
			DeleteCognitoCertificate:        logmd.Logging(logger, "delete", certificateTag, cognitoTag)(endpoints.DeleteCognitoCertificate),
			CreateStorageClass:              logmd.Logging(logger, "create", kubeTag, storageclassTag)(endpoints.CreateStorageClass),
			DeleteExternalSecrets:           logmd.Logging(logger, "delete", kubeTag, externalSecretsTag)(endpoints.DeleteExternalSecrets),
			CreateConfigMap:                 logmd.Logging(logger, "create", kubeTag, configMapTag)(endpoints.CreateConfigMap),
			DeleteConfigMap:                 logmd.Logging(logger, "delete", kubeTag, configMapTag)(endpoints.DeleteConfigMap),
			ScaleDeployment:                 logmd.Logging(logger, "create", kubeTag, scaleTag)(endpoints.ScaleDeployment),
			CreateHelmRelease:               logmd.Logging(logger, "create", helmTag, releasesTag)(endpoints.CreateHelmRelease),
			DeleteHelmRelease:               logmd.Logging(logger, "delete", helmTag, releasesTag)(endpoints.DeleteHelmRelease),
			GetHelmRelease:                  logmd.Logging(logger, "get", helmTag, releasesTag)(endpoints.GetHelmRelease),
			CreatePolicy:                    logmd.Logging(logger, "create", managedPoliciesTag)(endpoints.CreatePolicy),
			DeletePolicy:                    logmd.Logging(logger, "delete", managedPoliciesTag)(endpoints.DeletePolicy),
			CreateServiceAccount:            logmd.Logging(logger, "create", serviceAccountsTag)(endpoints.CreateServiceAccount),
			DeleteServiceAccount:            logmd.Logging(logger, "delete", serviceAccountsTag)(endpoints.DeleteServiceAccount),
			CreateNamespace:                 logmd.Logging(logger, "create", kubeTag, namespaceTag)(endpoints.CreateNamespace),
			CreatePostgresDatabase:          logmd.Logging(logger, "create", componentsTag, postgresTag)(endpoints.CreatePostgresDatabase),
			DeletePostgresDatabase:          logmd.Logging(logger, "delete", componentsTag, postgresTag)(endpoints.DeletePostgresDatabase),
			CreateS3Bucket:                  logmd.Logging(logger, "create", componentsTag, s3bucketTag)(endpoints.CreateS3Bucket),
			DeleteS3Bucket:                  logmd.Logging(logger, "delete", componentsTag, s3bucketTag)(endpoints.DeleteS3Bucket),
			CreateContainerRepository:       logmd.Logging(logger, "create", componentsTag, containerRepositoryTag)(endpoints.CreateContainerRepository),
			DeleteContainerRepository:       logmd.Logging(logger, "delete", componentsTag, containerRepositoryTag)(endpoints.DeleteContainerRepository),
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
