package storm

import (
	"errors"
	"time"

	"github.com/oslokommune/okctl/pkg/api"
)

// ErrNotFound is a not found error
var ErrNotFound = errors.New("not found")

// Metadata contains some useful metadata
// about a struct stored in storm
type Metadata struct {
	Identifier int `storm:"id,increment"`
	CreatedAt  time.Time
	UpdatedAt  time.Time `storm:"index"`
	Deleted    bool
}

// NewMetadata returns initialised metadata state
func NewMetadata() Metadata {
	return Metadata{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// ID contains the same content as an api.ID
// just modified for use with storm
type ID struct {
	Region       string
	AWSAccountID string
	ClusterName  string
}

// NewID returns an ID constructed from an
// api.ID
func NewID(id api.ID) ID {
	return ID{
		Region:       id.Region,
		AWSAccountID: id.AWSAccountID,
		ClusterName:  id.ClusterName,
	}
}

// Convert to an api.ID
func (i ID) Convert() api.ID {
	return api.ID{
		Region:       i.Region,
		AWSAccountID: i.AWSAccountID,
		ClusterName:  i.ClusterName,
	}
}
