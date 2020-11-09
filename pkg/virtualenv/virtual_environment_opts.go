package virtualenv

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// VirtualEnvironmentOpts contains the required inputs
type VirtualEnvironmentOpts struct {
	Region                 string
	AWSAccountID           string
	Environment            string
	Repository             string
	ClusterName            string
	UserDataDir            string
	Debug                  bool
	KubectlBinaryDir       string
	AwsIamAuthenticatorDir string
	OsEnvVars			   []string
	Ps1Dir                 string
}

// Validate the inputs
func (o *VirtualEnvironmentOpts) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.Environment, validation.Required),
		validation.Field(&o.AWSAccountID, validation.Required),
		validation.Field(&o.Region, validation.Required),
		validation.Field(&o.ClusterName, validation.Required),
		validation.Field(&o.UserDataDir, validation.Required),
		validation.Field(&o.KubectlBinaryDir, validation.Required),
		validation.Field(&o.AwsIamAuthenticatorDir, validation.Required),
	)
}
