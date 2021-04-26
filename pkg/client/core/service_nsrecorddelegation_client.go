package core

import (
	"fmt"

	"github.com/go-git/go-billy/v5/memfs"

	"github.com/oslokommune/okctl/pkg/git"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/github"
)

type nsRecordDelegationService struct {
	githubAPI github.Githuber
}

func (s *nsRecordDelegationService) RevokeDomainDelegation(opts client.RevokeDomainDelegationOpts) error {
	err := opts.Validate()
	if err != nil {
		return err
	}

	delegator := git.NewNameserverDelegator(
		true,
		git.DefaultWorkingDir,
		git.RepositoryStagerClone(git.DefaultRepositoryURL()),
		memfs.New(),
	)

	result, err := delegator.RevokeDelegation(opts.PrimaryHostedZoneFQDN, false)
	if err != nil {
		return fmt.Errorf("revoking dns zone delegation: %w", err)
	}

	if result.ModifiedRepository {
		err = s.githubAPI.CreatePullRequest(&github.PullRequest{
			Organisation:      github.DefaultOrg,
			Repository:        github.DefaultAWSInfrastructureRepository,
			SourceBranch:      result.Branch,
			DestinationBranch: github.DefaultAWSInfrastructurePrimaryBranch,
			Title:             fmt.Sprintf("❌ Hosted Zone revocation request for %s", opts.PrimaryHostedZoneFQDN),
			Body:              "Autogenerated by okctl",
			Labels:            opts.Labels,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *nsRecordDelegationService) InitiateDomainDelegation(opts client.InitiateDomainDelegationOpts) error {
	err := opts.Validate()
	if err != nil {
		return err
	}

	delegator := git.NewNameserverDelegator(
		true,
		git.DefaultWorkingDir,
		git.RepositoryStagerClone(git.DefaultRepositoryURL()),
		memfs.New(),
	)

	result, err := delegator.CreateDelegation(opts.PrimaryHostedZoneFQDN, opts.Nameservers, false)
	if err != nil {
		return fmt.Errorf("initiating dns zone delegation: %w", err)
	}

	if result.ModifiedRepository {
		err = s.githubAPI.CreatePullRequest(&github.PullRequest{
			Organisation:      github.DefaultOrg,
			Repository:        github.DefaultAWSInfrastructureRepository,
			SourceBranch:      result.Branch,
			DestinationBranch: github.DefaultAWSInfrastructurePrimaryBranch,
			Title:             fmt.Sprintf("✅ Hosted Zone delegation request for %s", opts.PrimaryHostedZoneFQDN),
			Body:              "Autogenerated by okctl",
			Labels:            opts.Labels,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// NewNameserverHandlerService initializes a new NSRecordDelegationService
func NewNameserverHandlerService(githubAPI github.Githuber) client.NSRecordDelegationService {
	return &nsRecordDelegationService{
		githubAPI: githubAPI,
	}
}
