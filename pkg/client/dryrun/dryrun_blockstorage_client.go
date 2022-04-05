package dryrun

import (
	"context"
	"fmt"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"io"
)

type blockstorageService struct {
	out io.Writer
}

func (b blockstorageService) CreateBlockstorage(_ context.Context, _ client.CreateBlockstorageOpts) (*client.Blockstorage, error) {
	fmt.Fprintf(b.out, formatCreate("Blockstorage controller"))

	return &client.Blockstorage{}, nil
}

func (b blockstorageService) DeleteBlockstorage(_ context.Context, _ api.ID) error {
	fmt.Fprintf(b.out, formatDelete("Blockstorage controller"))

	return nil
}
