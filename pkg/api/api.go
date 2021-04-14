package api

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	clusterMinLength = 3
	clusterMaxLength = 64
)

// ID contains the state that uniquely identifies a cluster
type ID struct {
	Region       string
	AWSAccountID string
	ClusterName  string
}

// Validate the identifier
func (i ID) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.Region, validation.Required),
		validation.Field(&i.ClusterName, validation.Required, validation.Length(clusterMinLength, clusterMaxLength)),
		validation.Field(&i.AWSAccountID, validation.Required, validation.Match(regexp.MustCompile("^[0-9]{12}$")).Error("must consist of 12 digits")),
	)
}
