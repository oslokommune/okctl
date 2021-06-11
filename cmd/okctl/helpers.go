package main

import (
	"github.com/oslokommune/okctl/pkg/okctl"
)

var (
	awsCredentialsType    string //nolint:gochecknoglobals
	githubCredentialsType string //nolint:gochecknoglobals
)

func enableServiceUserAuthentication(o *okctl.Okctl) {
	o.AWSCredentialsType = awsCredentialsType
	o.GithubCredentialsType = githubCredentialsType
}
