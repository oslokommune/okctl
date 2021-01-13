package core

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/core/api/git"
	"github.com/oslokommune/okctl/pkg/github"
	"github.com/oslokommune/okctl/pkg/spinner"
)

type nameserverRecordDelegationService struct {
	spin      spinner.Spinner
	githubAPI github.Githuber
}

// CreateNameserverRecordDelegationRequest knows how to request a NS record delegation
func (n *nameserverRecordDelegationService) CreateNameserverRecordDelegationRequest(
	opts *client.CreateNameserverDelegationRequestOpts,
) (*client.NameserverRecord, error) {
	delegator := git.NewNameserverDelegator(n.githubAPI)

	request, err := delegator.CreateNameserverDelegationRequest(opts.PrimaryHostedZoneFQDN, opts.Nameservers)
	if err != nil {
		return nil, fmt.Errorf("error creating nameserver record: %w", err)
	}

	if request.IsSubmitted() {
		return request.Record, nil
	}

	err = request.Submit()
	if err != nil {
		return nil, fmt.Errorf("error submitting nameserver delegation request: %w", err)
	}

	return request.Record, nil
}

// NewNameserverHandlerService initializes a new NameserverRecordDelegationService
func NewNameserverHandlerService(githubAPI github.Githuber, spin spinner.Spinner) client.NameserverRecordDelegationService {
	return &nameserverRecordDelegationService{
		spin:      spin,
		githubAPI: githubAPI,
	}
}
