package dryrun

import (
	"context"
	"fmt"
	"github.com/oslokommune/okctl/pkg/client"
	"io"
)

type postgresService struct {
	out io.Writer
}

func (p postgresService) CreatePostgresDatabase(_ context.Context, _ client.CreatePostgresDatabaseOpts) (*client.PostgresDatabase, error) {
	fmt.Fprintf(p.out, formatCreate("PostgreSQL"))

	return &client.PostgresDatabase{}, nil
}

func (p postgresService) DeletePostgresDatabase(_ context.Context, _ client.DeletePostgresDatabaseOpts) error {
	fmt.Fprintf(p.out, formatDelete("PostgreSQL"))

	return nil
}

func (p postgresService) GetPostgresDatabase(_ context.Context, _ client.GetPostgresDatabaseOpts) (*client.PostgresDatabase, error) {
	panic("implement me")
}
