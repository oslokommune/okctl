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

type Endpoints struct {
	CreateCluster endpoint.Endpoint
	DeleteCluster endpoint.Endpoint
}

func MakeEndpoints(s Services) Endpoints {
	return Endpoints{
		CreateCluster: makeCreateClusterEndpoint(s.Cluster),
		DeleteCluster: makeDeleteClusterEndpoint(s.Cluster),
	}
}

type Handlers struct {
	CreateCluster http.Handler
	DeleteCluster http.Handler
}

type EncodeResponseType string

const (
	EncodeJSONResponse EncodeResponseType = "json"
	EncodeYAMLResponse EncodeResponseType = "yaml"
	EncodeTextResponse EncodeResponseType = "text"
)

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

type Services struct {
	Cluster api.ClusterService
}

type EndpointOption func(Endpoints) Endpoints

func InstrumentEndpoints(logger *logrus.Logger) EndpointOption {
	return func(endpoints Endpoints) Endpoints {
		return Endpoints{
			CreateCluster: middleware.Logging(logger)(endpoints.CreateCluster),
			DeleteCluster: middleware.Logging(logger)(endpoints.DeleteCluster),
		}
	}
}

func GenerateEndpoints(s Services, opts ...EndpointOption) Endpoints {
	endpoints := MakeEndpoints(s)
	for _, opt := range opts {
		endpoints = opt(endpoints)
	}

	return endpoints
}
