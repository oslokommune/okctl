package core

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-kit/kit/endpoint"
	kit "github.com/go-kit/kit/transport/http"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/middleware"
	"github.com/oslokommune/okctl/pkg/transport"
	"github.com/sirupsen/logrus"
)

// Endpoints defines all available endpoints
type Endpoints struct {
	CreateCluster               endpoint.Endpoint
	DeleteCluster               endpoint.Endpoint
	CreateClusterConfig         endpoint.Endpoint
	CreateVpc                   endpoint.Endpoint
	DeleteVpc                   endpoint.Endpoint
	CreateExternalSecretsPolicy endpoint.Endpoint
}

// MakeEndpoints returns the endpoints initialised with their
// corresponding service
func MakeEndpoints(s Services) Endpoints {
	return Endpoints{
		CreateCluster:               makeCreateClusterEndpoint(s.Cluster),
		DeleteCluster:               makeDeleteClusterEndpoint(s.Cluster),
		CreateClusterConfig:         makeCreateClusterConfigEndpoint(s.ClusterConfig),
		CreateVpc:                   makeCreateVpcEndpoint(s.Vpc),
		DeleteVpc:                   makeDeleteVpcEndpoint(s.Vpc),
		CreateExternalSecretsPolicy: makeCreateExternalSecretsPolicyEndpoint(s.ManagedPolicy),
	}
}

// Handlers defines http handlers for processing requests
type Handlers struct {
	CreateCluster               http.Handler
	DeleteCluster               http.Handler
	CreateClusterConfig         http.Handler
	CreateVpc                   http.Handler
	DeleteVpc                   http.Handler
	CreateExternalSecretsPolicy http.Handler
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
		CreateCluster:               newServer(endpoints.CreateCluster, decodeClusterCreateRequest),
		DeleteCluster:               newServer(endpoints.DeleteCluster, decodeClusterDeleteRequest),
		CreateClusterConfig:         newServer(endpoints.CreateClusterConfig, decodeClusterConfigCreateRequest),
		CreateVpc:                   newServer(endpoints.CreateVpc, decodeVpcCreateRequest),
		DeleteVpc:                   newServer(endpoints.DeleteVpc, decodeVpcDeleteRequest),
		CreateExternalSecretsPolicy: newServer(endpoints.CreateExternalSecretsPolicy, decodeCreateExternalSecretsPolicyRequest),
	}
}

// AttachRoutes creates a router and adds handlers to the corresponding routes
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
		r.Route("/clusterconfigs", func(r chi.Router) {
			r.Method(http.MethodPost, "/", handlers.CreateClusterConfig)
		})
		r.Route("/managedpolicies", func(r chi.Router) {
			r.Route("/externalsecrets", func(r chi.Router) {
				r.Method(http.MethodPost, "/", handlers.CreateExternalSecretsPolicy)
			})
		})
	})

	return r
}

// Services defines all available services
type Services struct {
	Cluster       api.ClusterService
	ClusterConfig api.ClusterConfigService
	Vpc           api.VpcService
	ManagedPolicy api.ManagedPolicyService
}

// EndpointOption makes it easy to enable and disable the endpoint
// middlewares
type EndpointOption func(Endpoints) Endpoints

const (
	clusterTag       = "cluster"
	clusterConfigTag = "clusterConfig"
	vpcTag           = "vpc"
	managedPolicies  = "managedPolicies"
	externalSecrets  = "externalSecrets"
)

// InstrumentEndpoints adds instrumentation to the endpoints
func InstrumentEndpoints(logger *logrus.Logger) EndpointOption {
	return func(endpoints Endpoints) Endpoints {
		return Endpoints{
			CreateCluster:               middleware.Logging(logger, clusterTag, "create")(endpoints.CreateCluster),
			DeleteCluster:               middleware.Logging(logger, clusterTag, "delete")(endpoints.DeleteCluster),
			CreateClusterConfig:         middleware.Logging(logger, clusterConfigTag, "create")(endpoints.CreateClusterConfig),
			CreateVpc:                   middleware.Logging(logger, vpcTag, "create")(endpoints.CreateVpc),
			DeleteVpc:                   middleware.Logging(logger, vpcTag, "delete")(endpoints.DeleteVpc),
			CreateExternalSecretsPolicy: middleware.Logging(logger, strings.Join([]string{managedPolicies, externalSecrets}, "/"), "create")(endpoints.CreateExternalSecretsPolicy), // nolint: lll
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
