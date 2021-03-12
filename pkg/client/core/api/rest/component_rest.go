package rest

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

const (
	// TargetComponentPostgres matches the REST API route
	TargetComponentPostgres = "components/postgres/"
)

type componentAPI struct {
	client *HTTPClient
}

func (c *componentAPI) CreatePostgresDatabase(opts api.CreatePostgresDatabaseOpts) (*api.PostgresDatabase, error) {
	into := &api.PostgresDatabase{}
	return into, c.client.DoPost(TargetComponentPostgres, &opts, into)
}

func (c *componentAPI) DeletePostgresDatabase(opts api.DeletePostgresDatabaseOpts) error {
	return c.client.DoDelete(TargetComponentPostgres, &opts)
}

// NewComponentAPI returns an initialised REST API invoker
func NewComponentAPI(client *HTTPClient) client.ComponentAPI {
	return &componentAPI{
		client: client,
	}
}
