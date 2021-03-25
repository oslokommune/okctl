package storm

import (
	"time"

	"github.com/oslokommune/okctl/pkg/api"
)

// Metadata contains some useful metadata
// about a struct stored in storm
type Metadata struct {
	CreatedAt time.Time
	UpdatedAt time.Time `storm:"index"`
	Deleted   bool
}

// ID contains the same content as an api.ID
// just modified for use with storm
type ID struct {
	Region       string
	AWSAccountID string
	Environment  string
	Repository   string
	ClusterName  string
}

// NewID returns an ID constructed from an
// api.ID
func NewID(id *api.ID) *ID {
	return &ID{
		Region:       id.Region,
		AWSAccountID: id.AWSAccountID,
		Environment:  id.Environment,
		Repository:   id.Repository,
		ClusterName:  id.ClusterName,
	}
}

// Convert to an api.ID
func (i *ID) Convert() *api.ID {
	return &api.ID{
		Region:       i.Region,
		AWSAccountID: i.AWSAccountID,
		Environment:  i.Environment,
		Repository:   i.Repository,
		ClusterName:  i.ClusterName,
	}
}
