package client

import (
	"github.com/oslokommune/okctl/pkg/api"
)

// SecurityGroupAPI defines functionality required by the Security Group API
type SecurityGroupAPI interface {
	api.SecurityGroupCRUDer
	api.SecurityGroupRuleCRUDer
}
