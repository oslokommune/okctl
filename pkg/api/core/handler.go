package core

import (
	"net/http"

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
	CreateCluster endpoint.Endpoint
	DeleteCluster endpoint.Endpoint
}

// MakeEndpoints returns the endpoints initialised with their
// corresponding service
func MakeEndpoints(s Services) Endpoints {
	return Endpoints{
		CreateCluster: makeCreateClusterEndpoint(s.Cluster),
		DeleteCluster: makeDeleteClusterEndpoint(s.Cluster),
	}
}

// Handlers defines http handlers for processing requests
type Handlers struct {
	CreateCluster http.Handler
	DeleteCluster http.Handler
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
		CreateCluster: newServer(endpoints.CreateCluster, decodeClusterCreateRequest),
		DeleteCluster: newServer(endpoints.DeleteCluster, decodeClusterDeleteRequest),
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
	})

	return r
}

// Services defines all available services
type Services struct {
	Cluster api.ClusterService
}

// EndpointOption makes it easy to enable and disable the endpoint
// middlewares
type EndpointOption func(Endpoints) Endpoints

// InstrumentEndpoints adds instrumentation to the endpoints
func InstrumentEndpoints(logger *logrus.Logger) EndpointOption {
	return func(endpoints Endpoints) Endpoints {
		return Endpoints{
			CreateCluster: middleware.Logging(logger)(endpoints.CreateCluster),
			DeleteCluster: middleware.Logging(logger)(endpoints.DeleteCluster),
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
